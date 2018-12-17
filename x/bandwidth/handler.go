package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/cyberd/x/bandwidth/types"
)

var _ types.Handler = BaseHandler{}

type BaseHandler struct {
	// data providers
	accKeeper     auth.AccountKeeper
	stakeProvider types.AccStakeProvider
	bwKeeper      types.Keeper

	// bw configuration
	msgCost types.MsgBandwidthCost
}

func NewHandler(
	ak auth.AccountKeeper, sp types.AccStakeProvider, bwKeeper types.Keeper, msgCost types.MsgBandwidthCost,
) BaseHandler {

	return BaseHandler{
		accKeeper:     ak,
		stakeProvider: sp,
		bwKeeper:      bwKeeper,
		msgCost:       msgCost,
	}
}

func (h BaseHandler) GetTxCost(ctx sdk.Context, price float64, tx sdk.Tx) int64 {
	bandwidthForTx := TxCost
	for _, msg := range tx.GetMsgs() {
		bandwidthForTx = bandwidthForTx + h.msgCost(msg)
	}
	return bandwidthForTx
}

func (h BaseHandler) GetAccMaxBandwidth(ctx sdk.Context, addr sdk.AccAddress) int64 {
	accStakePercentage := h.stakeProvider.GetAccStakePercentage(ctx, addr)
	return int64(accStakePercentage * float64(MaxNetworkBandwidth) / 2)
}

func (h BaseHandler) GetCurrentAccBandwidth(ctx sdk.Context, address sdk.AccAddress) types.AcсBandwidth {
	accBw := h.bwKeeper.GetAccBandwidth(ctx, address)
	accMaxBw := h.GetAccMaxBandwidth(ctx, address)
	accBw.UpdateMax(accMaxBw, ctx.BlockHeight(), RecoveryPeriod)
	return accBw
}

// Double save for case:
// When acc send coins, we should consume bw before cutting max bw.
func (h BaseHandler) ConsumeAccBandwidth(ctx sdk.Context, bw types.AcсBandwidth, amt int64) {
	bw.Consume(amt)
	h.bwKeeper.SetAccBandwidth(ctx, bw)
	bw = h.GetCurrentAccBandwidth(ctx, bw.Address)
	h.bwKeeper.SetAccBandwidth(ctx, bw)
}
