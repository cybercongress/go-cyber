package bandwidth

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/cybercongress/cyberd/x/bandwidth/types"
)

var _ Keeper = BaseAccBandwidthKeeper{}

type BaseAccBandwidthKeeper struct {
	key *sdk.KVStoreKey
}

func NewAccBandwidthKeeper(key *sdk.KVStoreKey) BaseAccBandwidthKeeper {
	return BaseAccBandwidthKeeper{key: key}
}

func (bk BaseAccBandwidthKeeper) SetAccBandwidth(ctx sdk.Context, bandwidth AcсBandwidth) {
	bwBytes, _ := json.Marshal(bandwidth)
	ctx.KVStore(bk.key).Set(bandwidth.Address, bwBytes)
}

func (bk BaseAccBandwidthKeeper) GetAccBandwidth(ctx sdk.Context, addr sdk.AccAddress) (bw AcсBandwidth) {
	bwBytes := ctx.KVStore(bk.key).Get(addr)
	if bwBytes == nil {
		return AcсBandwidth{
			Address:          addr,
			RemainedValue:    0,
			LastUpdatedBlock: ctx.BlockHeight(),
			MaxValue:         0,
		}
	}
	err := json.Unmarshal(bwBytes, &bw)
	if err != nil {
		// should not happen
		panic("bandwidth index is broken")
	}
	return
}
