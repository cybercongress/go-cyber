package app

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	//bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	ctypes "github.com/cybercongress/go-cyber/types"

	//ctypes "github.com/cybercongress/go-cyber/types"
	bandwidthkeeper "github.com/cybercongress/go-cyber/x/bandwidth/keeper"
	bandwidthtypes "github.com/cybercongress/go-cyber/x/bandwidth/types"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(
	ak keeper.AccountKeeper, bankKeeper bankkeeper.Keeper,
	abk *bandwidthkeeper.BandwidthMeter,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewRejectExtensionOptionsDecorator(),
		NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.TxTimeoutHeightDecorator{},
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		ante.NewRejectFeeGranterDecorator(),
		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		NewDeductFeeBandRouterDecorator(ak, bankKeeper, abk),
		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		ante.NewSigVerificationDecorator(ak, signModeHandler),
		ante.NewIncrementSequenceDecorator(ak),
	)
}

// DeductFeeDecorator deducts fees from the first signer of the tx
// If the first signer does not have the funds to pay for the fees, return with InsufficientFunds error
// Call next AnteHandler if fees successfully deducted
// CONTRACT: Tx must implement FeeTx interface to use DeductFeeDecorator
type DeductFeeBandRouterDecorator struct {
	ak  keeper.AccountKeeper
	bk 	bankkeeper.Keeper
	bm	*bandwidthkeeper.BandwidthMeter
}

func NewDeductFeeBandRouterDecorator(ak keeper.AccountKeeper, bk bankkeeper.Keeper, bm *bandwidthkeeper.BandwidthMeter) DeductFeeBandRouterDecorator {
	return DeductFeeBandRouterDecorator{
		ak: ak,
		bk: bk,
		bm: bm,
	}
}

func (drd DeductFeeBandRouterDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {

	feeFlag := false
	feeSplitFlag := false
	bandwidthFlag := false
	ai2pay := sdk.AccAddress{}

	// temporary boundary to resolve node stuck on cosmovisor scan buffer limit
	if len(tx.GetMsgs()) > 25 {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrNotSupported, "Support only less than 25 msgs per tx")
	}

	// TODO optimize flat set
	for _, msg := range tx.GetMsgs() {
		if (msg.Route() != graphtypes.RouterKey && msg.Route() != wasm.RouterKey) {
			feeFlag = true
			//break
		}
		if (msg.Route() == wasm.RouterKey && msg.Type() != "execute") {
			feeFlag = true
			//break
		}
		if (msg.Route() == wasm.RouterKey && msg.Type() == "execute") {
			executeTx, ok := msg.(*wasm.MsgExecuteContract)
			if !ok {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Msg must be a MsgExecuteContract")
				//feeFlag = true
				//break
			}
			feeSplitFlag = true
			ai2pay, _ = sdk.AccAddressFromBech32(executeTx.Contract)
		}
		if (msg.Route() == graphtypes.RouterKey) {
			bandwidthFlag = true
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

	feePayer := feeTx.FeePayer()
	feePayerAcc := drd.ak.GetAccount(ctx, feePayer)

	if feePayerAcc == nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	// if feeFlag also true than fee of the whole tx will go to contract
	if feeSplitFlag {
		if !feeTx.GetFee().IsZero() {
			err = DeductFees(drd.bk, ctx, feePayerAcc, feeTx.GetFee(), ai2pay)
			if err != nil {
				return ctx, err
			}
		}

		return next(ctx, tx, simulate)
	}
	// default sdk case
	if feeFlag {
		if !feeTx.GetFee().IsZero() {
			err = DeductFees(drd.bk, ctx, feePayerAcc, feeTx.GetFee(), nil)
			if err != nil {
				return ctx, err
			}
		}

		return next(ctx, tx, simulate)
	}

	txCost := drd.bm.GetPricedTotalCyberlinksCost(ctx, tx)
	accountBandwidth := drd.bm.GetCurrentAccountBandwidth(ctx, feePayerAcc.GetAddress())

	currentBlockSpentBandwidth := drd.bm.GetCurrentBlockSpentBandwidth(ctx)
	maxBlockBandwidth := drd.bm.GetMaxBlockBandwidth(ctx)

	if !accountBandwidth.HasEnoughRemained(txCost) {
		return ctx, bandwidthtypes.ErrNotEnoughBandwidth
	} else if (txCost + currentBlockSpentBandwidth) > maxBlockBandwidth  {
		return ctx, bandwidthtypes.ErrExceededMaxBlockBandwidth
	} else {
		if !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
			_ = drd.bm.ConsumeAccountBandwidth(ctx, accountBandwidth, txCost)
			drd.bm.AddToBlockBandwidth(txCost)
		}
	}

	return next(ctx, tx, simulate)
}

// DeductFees deducts fees from the given account.
//
// NOTE: We could use the BankKeeper (in addition to the AccountKeeper, because
// the BankKeeper doesn't give us accounts), but it seems easier to do this.
func DeductFees(bankKeeper bankkeeper.Keeper, ctx sdk.Context, acc authtypes.AccountI, fees sdk.Coins, ai sdk.AccAddress) error {
	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	if ai == nil {
		err := bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, fees)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
		}
	} else {
		feeInCYB := sdk.NewDec(fees.AmountOf(ctypes.CYB).Int64())
		if feeInCYB.IsZero() == true {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, string(""))
		}
		toContract := feeInCYB.Mul(sdk.NewDecWithPrec(80,2))
		toValidators := feeInCYB.Sub(toContract)

		toValidatorsAmount := sdk.NewCoins(sdk.NewCoin(ctypes.CYB, toValidators.RoundInt()))
		toContractAmount := sdk.NewCoins(sdk.NewCoin(ctypes.CYB, toContract.RoundInt()))

		err := bankKeeper.SendCoins(ctx, acc.GetAddress(), ai, toContractAmount)
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
		if (msg.Route() != graphtypes.RouterKey) {
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
