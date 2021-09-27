package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/bandwidth/types"
)

func (bm BandwidthMeter) SetAccountBandwidth(ctx sdk.Context, ab types.AccountBandwidth) {
	ctx.KVStore(bm.storeKey).Set(types.AccountStoreKey(ab.Address), bm.cdc.MustMarshal(&ab))
}

func (bm BandwidthMeter) GetAccountBandwidth(ctx sdk.Context, address sdk.AccAddress) (ab types.AccountBandwidth) {
	bwBytes := ctx.KVStore(bm.storeKey).Get(types.AccountStoreKey(address.String()))
	if bwBytes == nil {
		return types.AccountBandwidth{
			Address:          address.String(),
			RemainedValue:    0,
			LastUpdatedBlock: uint64(ctx.BlockHeight()),
			MaxValue:         0,
		}
	}
	bm.cdc.MustUnmarshal(bwBytes, &ab)
	return
}
