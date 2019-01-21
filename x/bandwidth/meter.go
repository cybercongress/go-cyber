package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/cyberd/x/bandwidth/types"
)

var _ types.BandwidthMeter = BaseBandwidthMeter{}

type BaseBandwidthMeter struct {
	// data providers
	accKeeper     auth.AccountKeeper
	stakeProvider types.AccStakeProvider
	bwKeeper      types.Keeper

	// bw configuration
	msgCost types.MsgBandwidthCost
}

func NewBaseMeter(
	ak auth.AccountKeeper, sp types.AccStakeProvider, bwKeeper types.Keeper, msgCost types.MsgBandwidthCost,
) BaseBandwidthMeter {

	return BaseBandwidthMeter{
		accKeeper:     ak,
		stakeProvider: sp,
		bwKeeper:      bwKeeper,
		msgCost:       msgCost,
	}
}

func (h BaseBandwidthMeter) GetTxCost(ctx sdk.Context, price float64, tx sdk.Tx) int64 {
	bandwidthForTx := TxCost
	for _, msg := range tx.GetMsgs() {
		bandwidthForTx = bandwidthForTx + h.msgCost(msg)
	}
	return int64(float64(bandwidthForTx) * price)
}

func (h BaseBandwidthMeter) GetAccMaxBandwidth(ctx sdk.Context, addr sdk.AccAddress) int64 {
	accStakePercentage := h.stakeProvider.GetAccStakePercentage(ctx, addr)
	return int64(accStakePercentage * float64(DesirableNetworkBandwidthForRecoveryPeriod))
}

func (h BaseBandwidthMeter) GetCurrentAccBandwidth(ctx sdk.Context, address sdk.AccAddress) types.AcсBandwidth {
	accBw := h.bwKeeper.GetAccBandwidth(ctx, address)
	accMaxBw := h.GetAccMaxBandwidth(ctx, address)
	accBw.UpdateMax(accMaxBw, ctx.BlockHeight(), RecoveryPeriod)
	return accBw
}

func (h BaseBandwidthMeter) UpdateAccMaxBandwidth(ctx sdk.Context, address sdk.AccAddress) {
	bw := h.GetCurrentAccBandwidth(ctx, address)
	h.bwKeeper.SetAccBandwidth(ctx, bw)
}

//
// Performs bw consumption for given acc
// To get right number, should be called after tx delivery with bw state obtained prior delivery
//
// Pseudo code:
// bw := getCurrentBw(addr)
// bwCost := deliverTx(tx)
// consumeBw(bw, bwCost)
func (h BaseBandwidthMeter) ConsumeAccBandwidth(ctx sdk.Context, bw types.AcсBandwidth, amt int64) {
	bw.Consume(amt)
	h.bwKeeper.SetAccBandwidth(ctx, bw)
	bw = h.GetCurrentAccBandwidth(ctx, bw.Address)
	h.bwKeeper.SetAccBandwidth(ctx, bw)
}
