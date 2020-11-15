package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/bandwidth/types"
)

func (bk BandwidthMeter) SetBlockBandwidth(ctx sdk.Context, blockNumber uint64, value uint64) {
	store := ctx.KVStore(bk.storeKey)
	store.Set(types.BlockStoreKey(blockNumber), sdk.Uint64ToBigEndian(value))
}

func (bk BandwidthMeter) GetValuesForPeriod(ctx sdk.Context, period uint64) map[uint64]uint64 {
	store := ctx.KVStore(bk.storeKey)

	windowStart := uint64(ctx.BlockHeight()) - period + 1
	if windowStart <= 0 {
		windowStart = 1
	}

	result := make(map[uint64]uint64)
	for blockNumber := windowStart; blockNumber <= uint64(ctx.BlockHeight()); blockNumber++ {
		result[uint64(blockNumber)] = BigEndianToUint64(store.Get(types.BlockStoreKey(blockNumber)))
	}

	return result
}