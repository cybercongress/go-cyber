package bandwidth

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/types"
)

type BandwidthStateKeeper interface {
	SetAccountBandwidth(ctx sdk.Context, bandwidth types.AccountBandwidthState)
	GetAccountBandwidth(address sdk.AccAddress, ctx sdk.Context) (types.AccountBandwidthState, error)
}

type BaseBandwidthStateKeeper struct {
	key *sdk.KVStoreKey
}

func (bk *BaseBandwidthStateKeeper) SetAccountBandwidth(ctx sdk.Context, bandwidth types.AccountBandwidthState) {
	bwBytes, _ := json.Marshal(bandwidth)
	ctx.KVStore(bk.key).Set(bandwidth.Address, bwBytes)
}

func (bk *BaseBandwidthStateKeeper) GetAccountBandwidth(address sdk.AccAddress, ctx sdk.Context) (bw types.AccountBandwidthState, err error) {
	bwBytes := ctx.KVStore(bk.key).Get(address)
	err = json.Unmarshal(bwBytes, &bw)
	return
}
