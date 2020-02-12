package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/x/bandwidth/internal/types"
)

// Genesis accounts should contains fully restored bandwidth on block 0
func InitGenesis(ctx sdk.Context, handler types.BandwidthMeter,
	keeper AccountBandwidthKeeper, addresses []sdk.AccAddress, data GenesisState) {

	keeper.SetParams(ctx, data.Params)
	for _, address := range addresses {
		accMaxBw := handler.GetAccMaxBandwidth(ctx, address)
		keeper.SetAccountBandwidth(ctx, types.NewGenesisAccountBandwidth(address, accMaxBw))
	}
}

func ExportGenesis(ctx sdk.Context, keeper AccountBandwidthKeeper) GenesisState {
	return NewGenesisState(keeper.GetParams(ctx))
}
