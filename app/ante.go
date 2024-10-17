package app

import (
	"fmt"
	"math"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	graphtypes "github.com/cybercongress/go-cyber/v5/x/graph/types"

	// antefee "github.com/cosmos/cosmos-sdk/x/auth/ante/fee"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	ctypes "github.com/cybercongress/go-cyber/v5/types"

	ibcante "github.com/cosmos/ibc-go/v7/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	bandwidthkeeper "github.com/cybercongress/go-cyber/v5/x/bandwidth/keeper"
	bandwidthtypes "github.com/cybercongress/go-cyber/v5/x/bandwidth/types"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	ante.HandlerOptions

	BandwidthMeter    *bandwidthkeeper.BandwidthMeter
	IBCKeeper         *ibckeeper.Keeper
	WasmConfig        *wasmtypes.WasmConfig
	WasmKeeper        *wasmkeeper.Keeper
	TXCounterStoreKey storetypes.StoreKey

	TxEncoder sdk.TxEncoder
}

func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "account keeper is required for AnteHandler")
	}

	if options.BankKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "bank keeper is required for AnteHandler")
	}

	if options.SignModeHandler == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	if options.WasmConfig == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "wasm config is required for ante builder")
	}

	if options.TXCounterStoreKey == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "tx counter key is required for ante builder")
	}

	if options.BandwidthMeter == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "bandwidth meter is required for AnteHandler")
	}

	sigGasConsumer := options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit), // after setup context to enforce limits early
		wasmkeeper.NewCountTXDecorator(options.TXCounterStoreKey),
		wasmkeeper.NewGasRegisterDecorator(options.WasmKeeper.GetGasRegister()),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		// TODO rewrite fee decorator
		NewDeductFeeBandRouterDecorator(options.AccountKeeper, options.BankKeeper, options.BandwidthMeter, options.FeegrantKeeper, options.TxFeeChecker),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}

type DeductFeeBandDecorator struct {
	accountKeeper  ante.AccountKeeper
	bankKeeper     authtypes.BankKeeper
	bandwidthMeter *bandwidthkeeper.BandwidthMeter
	feegrantKeeper ante.FeegrantKeeper
	txFeeChecker   ante.TxFeeChecker
}

func NewDeductFeeBandRouterDecorator(
	ak ante.AccountKeeper,
	bk authtypes.BankKeeper,
	bm *bandwidthkeeper.BandwidthMeter,
	fk ante.FeegrantKeeper,
	tfc ante.TxFeeChecker,
) DeductFeeBandDecorator {
	if tfc == nil {
		tfc = checkTxFeeWithValidatorMinGasPrices
	}

	return DeductFeeBandDecorator{
		accountKeeper:  ak,
		bankKeeper:     bk,
		bandwidthMeter: bm,
		feegrantKeeper: fk,
		txFeeChecker:   tfc,
	}
}

// TODO check line by line
func (dfd DeductFeeBandDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	var (
		priority int64
		err      error
	)

	feeFlag := false
	feeSplitFlag := false
	bandwidthFlag := false
	program2pay := sdk.AccAddress{}

	for _, msg := range tx.GetMsgs() {
		if (sdk.MsgTypeURL(msg) == sdk.MsgTypeURL(&graphtypes.MsgCyberlink{})) {
			bandwidthFlag = true
		}
		if (sdk.MsgTypeURL(msg) != sdk.MsgTypeURL(&graphtypes.MsgCyberlink{})) {
			feeFlag = true
		}
		if sdk.MsgTypeURL(msg) == sdk.MsgTypeURL(&wasmtypes.MsgExecuteContract{}) {
			executeTx, ok := msg.(*wasm.MsgExecuteContract)
			if !ok {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Msg must be a MsgExecuteContract")
			}
			feeSplitFlag = true
			program2pay, err = sdk.AccAddressFromBech32(executeTx.Contract)
			if err != nil {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Address must be valid")
			}
		}
	}

	// not allow to merge bandwidth/fee billing in one transaction
	if bandwidthFlag && (feeFlag || feeSplitFlag) {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrNotSupported, "Support only batch of cyberlinks")
	}

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if !simulate && ctx.BlockHeight() > 0 && feeTx.GetGas() == 0 {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidGasLimit, "must provide positive gas")
	}

	fee := feeTx.GetFee()
	if !simulate {
		if feeFlag || feeSplitFlag {
			fee, priority, err = dfd.txFeeChecker(ctx, tx)
			if err != nil {
				return ctx, err
			}
		}
	}

	if feeFlag || feeSplitFlag {
		err = dfd.checkDeductFee(ctx, tx, fee, program2pay)
		if err != nil {
			return ctx, err
		}
		newCtx := ctx.WithPriority(priority)

		return next(newCtx, tx, simulate)
	}

	bandwidthPayer := feeTx.FeePayer()
	bandwidthGranter := feeTx.FeeGranter()
	deductBandwidthFrom := bandwidthPayer

	if bandwidthGranter != nil {
		deductBandwidthFrom = bandwidthGranter
	}

	deductBandwidthFromAcc := dfd.accountKeeper.GetAccount(ctx, deductBandwidthFrom)
	if deductBandwidthFromAcc == nil {
		return ctx, sdkerrors.ErrUnknownAddress.Wrapf("fee payer address: %s does not exist", deductBandwidthFrom)
	}

	txCost := dfd.bandwidthMeter.GetPricedTotalCyberlinksCost(ctx, tx)
	accountBandwidth := dfd.bandwidthMeter.GetCurrentAccountBandwidth(ctx, deductBandwidthFromAcc.GetAddress())

	currentBlockSpentBandwidth := dfd.bandwidthMeter.GetCurrentBlockSpentBandwidth(ctx)
	maxBlockBandwidth := dfd.bandwidthMeter.GetMaxBlockBandwidth(ctx)

	if !accountBandwidth.HasEnoughRemained(txCost) {
		return ctx, bandwidthtypes.ErrNotEnoughBandwidth
	}

	if (txCost + currentBlockSpentBandwidth) > maxBlockBandwidth {
		return ctx, bandwidthtypes.ErrExceededMaxBlockBandwidth
	}

	isDeliverTx := !ctx.IsCheckTx() && !ctx.IsReCheckTx() && !simulate
	if isDeliverTx {
		err = dfd.bandwidthMeter.ConsumeAccountBandwidth(ctx, accountBandwidth, txCost)
		if err != nil {
			return ctx, err
		}
		dfd.bandwidthMeter.AddToBlockBandwidth(ctx, txCost)
	}

	newCtx := ctx.WithPriority(int64(math.MaxInt64))

	return next(newCtx, tx, simulate)
}

// Wraps fee.checkDeductFee with addition of program2pay address in case of contract execution
func (dfd DeductFeeBandDecorator) checkDeductFee(ctx sdk.Context, sdkTx sdk.Tx, fee sdk.Coins, program2pay sdk.AccAddress) error {
	feeTx, ok := sdkTx.(sdk.FeeTx)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if addr := dfd.accountKeeper.GetModuleAddress(authtypes.FeeCollectorName); addr == nil {
		return fmt.Errorf("fee collector module account (%s) has not been set", authtypes.FeeCollectorName)
	}

	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()
	deductFeesFrom := feePayer

	// if feegranter set deduct fee from feegranter account.
	// this works with only when feegrant enabled.
	if feeGranter != nil {
		if dfd.feegrantKeeper == nil {
			return sdkerrors.ErrInvalidRequest.Wrap("fee grants are not enabled")
		} else if !feeGranter.Equals(feePayer) {
			err := dfd.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, fee, sdkTx.GetMsgs())
			if err != nil {
				return sdkerrors.Wrapf(err, "%s does not allow to pay fees for %s", feeGranter, feePayer)
			}
		}

		deductFeesFrom = feeGranter
	}

	deductFeesFromAcc := dfd.accountKeeper.GetAccount(ctx, deductFeesFrom)
	if deductFeesFromAcc == nil {
		return sdkerrors.ErrUnknownAddress.Wrapf("fee payer address: %s does not exist", deductFeesFrom)
	}

	// deduct the fees
	if !fee.IsZero() {
		if program2pay == nil {
			err := DeductFees(dfd.bankKeeper, ctx, deductFeesFromAcc, fee, nil)
			if err != nil {
				return err
			}
		} else {
			err := DeductFees(dfd.bankKeeper, ctx, deductFeesFromAcc, fee, program2pay)
			if err != nil {
				return err
			}
		}
	}

	events := sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeTx,
			sdk.NewAttribute(sdk.AttributeKeyFee, fee.String()),
			sdk.NewAttribute(sdk.AttributeKeyFeePayer, deductFeesFrom.String()),
		),
	}
	ctx.EventManager().EmitEvents(events)

	return nil
}

// DeductFees deducts fees from the given account.
// Wraps fee.DeductFees with addition of program2pay address in case of contract execution
func DeductFees(bankKeeper authtypes.BankKeeper, ctx sdk.Context, acc authtypes.AccountI, fees sdk.Coins, program2pay sdk.AccAddress) error {
	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	if program2pay.Empty() {
		err := bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), authtypes.FeeCollectorName, fees)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
		}
	} else {
		feeInCYB := sdk.NewDec(fees.AmountOf(ctypes.CYB).Int64())
		toProgram := feeInCYB.Mul(sdk.NewDecWithPrec(80, 2))
		toValidators := feeInCYB.Sub(toProgram)

		toValidatorsAmount := sdk.NewCoins(sdk.NewCoin(ctypes.CYB, toValidators.RoundInt()))
		toProgramAmount := sdk.NewCoins(sdk.NewCoin(ctypes.CYB, toProgram.RoundInt()))

		err := bankKeeper.SendCoins(ctx, acc.GetAddress(), program2pay, toProgramAmount)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
		}

		err = bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), authtypes.FeeCollectorName, toValidatorsAmount)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
		}
	}

	return nil
}

// checkTxFeeWithValidatorMinGasPrices implements the default fee logic, where the minimum price per
// unit of gas is fixed and set by each validator, can the tx priority is computed from the gas price.
func checkTxFeeWithValidatorMinGasPrices(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return nil, 0, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

	// Ensure that the provided fees meet a minimum threshold for the validator,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if ctx.IsCheckTx() {
		minGasPrices := ctx.MinGasPrices()
		if !minGasPrices.IsZero() {
			requiredFees := make(sdk.Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := sdkmath.LegacyNewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}

			if !feeCoins.IsAnyGTE(requiredFees) {
				return nil, 0, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, requiredFees)
			}
		}
	}

	priority := getTxPriority(feeCoins, int64(gas))
	return feeCoins, priority, nil
}

// getTxPriority returns a naive tx priority based on the amount of the smallest denomination of the gas price
// provided in a transaction.
// NOTE: This implementation should be used with a great consideration as it opens potential attack vectors
// where txs with multiple coins could not be prioritize as expected.
func getTxPriority(fee sdk.Coins, gas int64) int64 {
	var priority int64
	for _, c := range fee {
		p := int64(math.MaxInt64)
		gasPrice := c.Amount.QuoRaw(gas)
		if gasPrice.IsInt64() {
			p = gasPrice.Int64()
		}
		if priority == 0 || p < priority {
			priority = p
		}
	}

	return priority
}
