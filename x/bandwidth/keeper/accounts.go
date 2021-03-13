package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/bandwidth/types"
)

func (bk BandwidthMeter) SetAccountBandwidth(ctx sdk.Context, ab types.AccountBandwidth) {
	ctx.KVStore(bk.StoreKey).Set(types.AccountStoreKey(ab.Address), bk.Cdc.MustMarshalBinaryBare(&ab))
}

func (bk BandwidthMeter) GetAccountBandwidth(ctx sdk.Context, address sdk.AccAddress) (ab types.AccountBandwidth) {
	bwBytes := ctx.KVStore(bk.StoreKey).Get(types.AccountStoreKey(address.String()))
	if bwBytes == nil {
		return types.AccountBandwidth{
			Address:          address.String(),
			RemainedValue:    0,
			LastUpdatedBlock: uint64(ctx.BlockHeight()),
			MaxValue:         0,
		}
	}
	bk.Cdc.MustUnmarshalBinaryBare(bwBytes, &ab)
	return
}
