package bandwidth

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/x/bandwidth/types"
)

type AccountBandwidthKeeper interface {
	SetAccountBandwidth(ctx sdk.Context, bandwidth types.AccountBandwidth)
	GetAccountBandwidth(address sdk.AccAddress, ctx sdk.Context) (types.AccountBandwidth, error)
}

type BaseAccountBandwidthKeeper struct {
	key *sdk.KVStoreKey
}

func (bk BaseAccountBandwidthKeeper) SetAccountBandwidth(ctx sdk.Context, bandwidth types.AccountBandwidth) {
	bwBytes, _ := json.Marshal(bandwidth)
	ctx.KVStore(bk.key).Set(bandwidth.Address, bwBytes)
}

func (bk BaseAccountBandwidthKeeper) GetAccountBandwidth(address sdk.AccAddress, ctx sdk.Context) (bw types.AccountBandwidth, err error) {
	bwBytes := ctx.KVStore(bk.key).Get(address)
	if bwBytes == nil {
		return types.AccountBandwidth{
			Address:          address,
			RemainedValue:    0,
			LastUpdatedBlock: ctx.BlockHeight(),
			MaxValue:         0,
		}, nil
	}
	err = json.Unmarshal(bwBytes, &bw)
	return
}

func NewAccountBandwidthKeeper(key *sdk.KVStoreKey) BaseAccountBandwidthKeeper {
	return BaseAccountBandwidthKeeper{key: key}
}
