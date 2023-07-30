package keeper

import (
	"github.com/cybercongress/go-cyber/x/bandwidth/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (bm BandwidthMeter) SetAccountBandwidth(ctx sdk.Context, ab types.NeuronBandwidth) {
	ctx.KVStore(bm.storeKey).Set(types.AccountStoreKey(ab.Neuron), bm.cdc.MustMarshal(&ab))
}

func (bm BandwidthMeter) GetAccountBandwidth(ctx sdk.Context, address sdk.AccAddress) (ab types.NeuronBandwidth) {
	bwBytes := ctx.KVStore(bm.storeKey).Get(types.AccountStoreKey(address.String()))
	if bwBytes == nil {
		return types.NeuronBandwidth{
			Neuron:           address.String(),
			RemainedValue:    0,
			LastUpdatedBlock: uint64(ctx.BlockHeight()),
			MaxValue:         0,
		}
	}
	bm.cdc.MustUnmarshal(bwBytes, &ab)
	return
}
