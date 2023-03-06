package app

import (
	"fmt"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	ctypes "github.com/cybercongress/go-cyber/types"

	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	ibcante "github.com/cosmos/ibc-go/v4/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v4/modules/core/keeper"
	bandwidthkeeper "github.com/cybercongress/go-cyber/x/bandwidth/keeper"
	bandwidthtypes "github.com/cybercongress/go-cyber/x/bandwidth/types"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
)

// HandlerOptions are the options required for constructing a default SDK AnteHandler.
type HandlerBaseOptions struct {
	AccountKeeper   ante.AccountKeeper
	BankKeeper      bankkeeper.Keeper
	FeegrantKeeper  ante.FeegrantKeeper
	SignModeHandler authsigning.SignModeHandler
	SigGasConsumer  func(meter sdk.GasMeter, sig signing.SignatureV2, params types.Params) error
}

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	HandlerBaseOptions

	BandwidthMeter    *bandwidthkeeper.BandwidthMeter
	IBCKeeper         *ibckeeper.Keeper
	WasmConfig        *wasmTypes.WasmConfig
	TXCounterStoreKey sdk.StoreKey
}

func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for AnteHandler")
	}

	if options.BankKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for AnteHandler")
	}

	if options.SignModeHandler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	if options.WasmConfig == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "wasm config is required for ante builder")
	}

	if options.TXCounterStoreKey == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "tx counter key is required for ante builder")
	}

	if options.BandwidthMeter == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "bandwidth meter is required for AnteHandler")
	}

	var sigGasConsumer = options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit), // after setup context to enforce limits early
		wasmkeeper.NewCountTXDecorator(options.TXCounterStoreKey),
		ante.NewRejectExtensionOptionsDecorator(),
		NewMempoolFeeDecorator(), // overwrite ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		// overwrite ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper),
		NewDeductFeeBandRouterDecorator(options.AccountKeeper, options.BankKeeper, options.BandwidthMeter, options.FeegrantKeeper),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewAnteDecorator(options.IBCKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}

type DeductFeeBandDecorator struct {
	ak             ante.AccountKeeper
	bankKeeper     bankkeeper.Keeper
	bandMeter      *bandwidthkeeper.BandwidthMeter
	feegrantKeeper ante.FeegrantKeeper
}

func NewDeductFeeBandRouterDecorator(
	ak ante.AccountKeeper,
	bk bankkeeper.Keeper,
	bm *bandwidthkeeper.BandwidthMeter,
	fk ante.FeegrantKeeper,
) DeductFeeBandDecorator {
	return DeductFeeBandDecorator{
		ak:             ak,
		bankKeeper:     bk,
		bandMeter:      bm,
		feegrantKeeper: fk,
	}
}

func (drd DeductFeeBandDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {

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

	if addr := drd.ak.GetModuleAddress(types.FeeCollectorName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.FeeCollectorName))
	}

	fee := feeTx.GetFee()
	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()
	deductFeesFrom := feePayer

	// if feegranter set deduct fee from feegranter account.
	// this works with only when feegrant enabled.
	if feeGranter != nil {
		if drd.feegrantKeeper == nil {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "fee grants are not enabled")
		} else if !feeGranter.Equals(feePayer) {
			err := drd.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, fee, tx.GetMsgs())

			if err != nil {
				return ctx, sdkerrors.Wrapf(err, "%s not allowed to pay fees from %s", feeGranter, feePayer)
			}
		}

		deductFeesFrom = feeGranter
	}
	deductFeesFromAcc := drd.ak.GetAccount(ctx, deductFeesFrom)
	if deductFeesFromAcc == nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	// if feeFlag also true than fee of the whole tx will go to contract
	if feeSplitFlag {
		if !feeTx.GetFee().IsZero() {
			err = DeductFees(drd.bankKeeper, ctx, deductFeesFromAcc, feeTx.GetFee(), program2pay)
			if err != nil {
				return ctx, err
			}
		}

		events := sdk.Events{sdk.NewEvent(sdk.EventTypeTx,
			sdk.NewAttribute(sdk.AttributeKeyFee, feeTx.GetFee().String()),
		)}
		ctx.EventManager().EmitEvents(events)

		return next(ctx, tx, simulate)
	}
	// default sdk case
	if feeFlag {
		if !feeTx.GetFee().IsZero() {
			err = DeductFees(drd.bankKeeper, ctx, deductFeesFromAcc, feeTx.GetFee(), nil)
			if err != nil {
				return ctx, err
			}
		}

		events := sdk.Events{sdk.NewEvent(sdk.EventTypeTx,
			sdk.NewAttribute(sdk.AttributeKeyFee, feeTx.GetFee().String()),
		)}
		ctx.EventManager().EmitEvents(events)

		return next(ctx, tx, simulate)
	}

	txCost := drd.bandMeter.GetPricedTotalCyberlinksCost(ctx, tx)
	accountBandwidth := drd.bandMeter.GetCurrentAccountBandwidth(ctx, deductFeesFromAcc.GetAddress())

	currentBlockSpentBandwidth := drd.bandMeter.GetCurrentBlockSpentBandwidth(ctx)
	maxBlockBandwidth := drd.bandMeter.GetMaxBlockBandwidth(ctx)

	if !simulate {
		if !accountBandwidth.HasEnoughRemained(txCost) {
			return ctx, bandwidthtypes.ErrNotEnoughBandwidth
		} else if (txCost + currentBlockSpentBandwidth) > maxBlockBandwidth {
			return ctx, bandwidthtypes.ErrExceededMaxBlockBandwidth
		} else {
			if !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
				err = drd.bandMeter.ConsumeAccountBandwidth(ctx, accountBandwidth, txCost)
				if err != nil {
					return ctx, err
				}
				// TODO think to add to transient store
				drd.bandMeter.AddToBlockBandwidth(txCost)
			}
		}
	}
	return next(ctx, tx, simulate)
}

// DeductFees deducts fees from the given account.
//
// NOTE: We could use the BankKeeper (in addition to the AccountKeeper, because
// the BankKeeper doesn't give us accounts), but it seems easier to do this.
func DeductFees(bankKeeper bankkeeper.Keeper, ctx sdk.Context, acc authtypes.AccountI, fees sdk.Coins, program2pay sdk.AccAddress) error {
	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	if program2pay == nil {
		err := bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, fees)
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

		err = bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, toValidatorsAmount)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
		}
	}

	return nil
}

/*-------------------------------------------------------*/

// MempoolFeeDecorator will check if the transaction's fee is at least as large
// as the local validator's minimum gasFee (defined in validator config).
// If fee is too low, decorator returns error and tx is rejected from mempool.
// Note this only applies when ctx.CheckTx = true
// If fee is high enough or not CheckTx, then call next AnteHandler
// CONTRACT: Tx must implement FeeTx to use MempoolFeeDecorator
type MempoolFeeDecorator struct{}

func NewMempoolFeeDecorator() MempoolFeeDecorator {
	return MempoolFeeDecorator{}
}

func (mfd MempoolFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}
	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

	cyberlinksNoFeeFlag := true
	for _, msg := range tx.GetMsgs() {
		if (sdk.MsgTypeURL(msg) != sdk.MsgTypeURL(&graphtypes.MsgCyberlink{})) {
			cyberlinksNoFeeFlag = false
			break
		}
	}

	// Ensure that the provided fees meet a minimum threshold for the validator,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if ctx.IsCheckTx() && !simulate && !cyberlinksNoFeeFlag {
		minGasPrices := ctx.MinGasPrices()
		if !minGasPrices.IsZero() {
			requiredFees := make(sdk.Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := sdk.NewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}

			if !feeCoins.IsAnyGTE(requiredFees) {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, requiredFees)
			}
		}
	}
	return next(ctx, tx, simulate)
}
