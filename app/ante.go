package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cybercongress/go-cyber/x/bandwidth"
)

func NewAnteHandler(
	ak keeper.AccountKeeper,
	abk *bandwidth.BandwidthMeter,
	supplyKeeper types.SupplyKeeper,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		ante.NewSetPubKeyDecorator(ak),
		ante.NewValidateSigCountDecorator(ak),
		ante.NewDeductFeeDecorator(ak, supplyKeeper),
		NewDeductBandwidthDecorator(ak, abk),
		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		ante.NewSigVerificationDecorator(ak),
		ante.NewIncrementSequenceDecorator(ak),
	)
}

var (
	_ FeeTx = (*types.StdTx)(nil)
)

type FeeTx interface {
	sdk.Tx
	GetGas() uint64
	GetFee() sdk.Coins
	FeePayer() sdk.AccAddress
}


type DeductBandwidthDecorator struct {
	ak         auth.AccountKeeper
	bm		   *bandwidth.BandwidthMeter
}

func NewDeductBandwidthDecorator(ak auth.AccountKeeper, bm *bandwidth.BandwidthMeter) DeductBandwidthDecorator {
	return DeductBandwidthDecorator{
		ak:         ak,
		bm: bm,
	}
}

func (dbd DeductBandwidthDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feePayer := feeTx.FeePayer()
	feePayerAcc := dbd.ak.GetAccount(ctx, feePayer)

	if feePayerAcc == nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	txCost := dbd.bm.GetPricedTxCost(ctx, tx)
	accountBandwidth := dbd.bm.GetCurrentAccountBandwidth(ctx, feePayerAcc.GetAddress())

	currentBlockSpentBandwidth := dbd.bm.GetCurrentBlockSpentBandwidth(ctx)
	maxBlockBandwidth := dbd.bm.GetMaxBlockBandwidth(ctx)

	if !accountBandwidth.HasEnoughRemained(txCost) {
		return ctx, bandwidth.ErrNotEnoughBandwidth
	} else if (uint64(txCost) + currentBlockSpentBandwidth) > maxBlockBandwidth  {
		return ctx, bandwidth.ErrExceededMaxBlockBandwidth
	} else {
		dbd.bm.ConsumeAccountBandwidth(ctx, accountBandwidth, txCost)
		dbd.bm.AddToBlockBandwidth(txCost)
	}

	return next(ctx, tx, simulate)
}

