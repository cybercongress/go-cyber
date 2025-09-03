package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	"github.com/cybercongress/go-cyber/v6/x/bandwidth/types"
)

func InitGenesis(ctx sdk.Context, bm *BandwidthMeter, _ authkeeper.AccountKeeper, data types.GenesisState) {
	if err := bm.SetParams(ctx, data.Params); err != nil {
		panic(err)
	}

	bm.currentCreditPrice = bm.GetBandwidthPrice(ctx, data.Params.BasePrice)
}

func ExportGenesis(ctx sdk.Context, bm *BandwidthMeter) *types.GenesisState {
	return types.NewGenesisState(bm.GetParams(ctx))
}
