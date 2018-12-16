package bandwidth

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/x/bandwidth/types"
)

type AccStakeProvider interface {
	GetAccStakePercentage(ctx sdk.Context, address sdk.AccAddress) float64
}

type AccountBandwidthKeeper interface {
	SetAccBandwidth(ctx sdk.Context, bandwidth types.AcсBandwidth)
	// Returns updated to current acc bandwidth
	GetCurrentAccBandwidth(ctx sdk.Context, address sdk.AccAddress) types.AcсBandwidth
	GetAccMaxBandwidth(ctx sdk.Context, address sdk.AccAddress) int64
}

type BaseAccBandwidthKeeper struct {
	key           *sdk.KVStoreKey
	stakeProvider AccStakeProvider
}

func NewAccountBandwidthKeeper(key *sdk.KVStoreKey, stakeProvider AccStakeProvider) BaseAccBandwidthKeeper {
	return BaseAccBandwidthKeeper{
		key:           key,
		stakeProvider: stakeProvider,
	}
}

func (bk BaseAccBandwidthKeeper) GetAccMaxBandwidth(ctx sdk.Context, addr sdk.AccAddress) int64 {
	accStakePercentage := bk.stakeProvider.GetAccStakePercentage(ctx, addr)
	return int64(accStakePercentage * float64(MaxNetworkBandwidth) / 2)
}

func (bk BaseAccBandwidthKeeper) SetAccBandwidth(ctx sdk.Context, bandwidth types.AcсBandwidth) {
	bwBytes, _ := json.Marshal(bandwidth)
	ctx.KVStore(bk.key).Set(bandwidth.Address, bwBytes)
}

func (bk BaseAccBandwidthKeeper) GetCurrentAccBandwidth(ctx sdk.Context, address sdk.AccAddress) types.AcсBandwidth {

	bwBytes := ctx.KVStore(bk.key).Get(address)
	var accBw types.AcсBandwidth
	if bwBytes == nil {
		accBw = types.AcсBandwidth{
			Address:          address,
			RemainedValue:    0,
			LastUpdatedBlock: ctx.BlockHeight(),
			MaxValue:         0,
		}
	} else {
		_ = json.Unmarshal(bwBytes, &accBw)
	}
	accMaxBw := bk.GetAccMaxBandwidth(ctx, address)
	accBw.UpdateMax(accMaxBw, ctx.BlockHeight(), RecoveryPeriod)
	return accBw
}
