package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/x/bandwidth/types"
)

// Genesis accounts should contains fully restored bandwidth on block 0
func InitGenesis(ctx sdk.Context, bwHandler types.BandwidthMeter, bwKeeper types.Keeper, addresses []sdk.AccAddress) {

	for _, address := range addresses {
		accMaxBw := bwHandler.GetAccMaxBandwidth(ctx, address)
		bwKeeper.SetAccBandwidth(ctx, types.NewGenesisAccBandwidth(address, accMaxBw))
	}
}
