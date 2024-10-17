package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v5/x/grid/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	// TODO check EnergyGrid balance aggregated value on genesis
	if err := k.SetParams(ctx, data.Params); err != nil {
		panic(err)
	}
	err := k.SetRoutes(ctx, data.Routes)
	if err != nil {
		panic(err)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	//routes := k.GetAllRoutes(ctx)
	// Create an empty slice of types.Route
	var routes []types.Route

	return types.NewGenesisState(params, routes)
}
