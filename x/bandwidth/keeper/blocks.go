package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/bandwidth/types"
)

func (bm BandwidthMeter) SetBlockBandwidth(ctx sdk.Context, blockNumber uint64, value uint64) {
	store := ctx.KVStore(bm.storeKey)
	store.Set(types.BlockStoreKey(blockNumber), sdk.Uint64ToBigEndian(value))
}

func (bm BandwidthMeter) GetValuesForPeriod(ctx sdk.Context, period uint64) map[uint64]uint64 {
	store := ctx.KVStore(bm.storeKey)

	windowStart := ctx.BlockHeight() - int64(period) + 1
	if windowStart <= 0 {
		windowStart = 1
	}

	result := make(map[uint64]uint64)
	for blockNumber := uint64(windowStart); blockNumber <= uint64(ctx.BlockHeight()); blockNumber++ {
		result[blockNumber] = sdk.BigEndianToUint64(store.Get(types.BlockStoreKey(blockNumber)))
	}
	return result
}
