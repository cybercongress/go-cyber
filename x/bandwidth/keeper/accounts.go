package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/bandwidth/types"
)

func (bk BandwidthMeter) SetAccountBandwidth(ctx sdk.Context, ab types.AcсountBandwidth) {
	bwBytes, _ := json.Marshal(ab)
	ctx.KVStore(bk.storeKey).Set(types.AccountStoreKey(ab.Address), bwBytes)
}

func (bk BandwidthMeter) GetAccountBandwidth(ctx sdk.Context, address sdk.AccAddress) (ak types.AcсountBandwidth) {
	bwBytes := ctx.KVStore(bk.storeKey).Get(types.AccountStoreKey(address))
	if bwBytes == nil {
		return types.AcсountBandwidth{
			Address:          address,
			RemainedValue:    0,
			LastUpdatedBlock: uint64(ctx.BlockHeight()),
			MaxValue:         0,
		}
	}
	err := json.Unmarshal(bwBytes, &ak)
	if err != nil {
		panic("bandwidth index is broken")
	}
	return
}
