package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/joinresistance/space-pussy/x/grid/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	// TODO check EnergyGrid balance aggregated value on genesis
	k.SetParams(ctx, data.Params)
	err := k.SetRoutes(ctx, data.Routes)
	if err != nil {
		panic(err)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	routes := k.GetAllRoutes(ctx)

	return types.NewGenesisState(params, routes)
}
