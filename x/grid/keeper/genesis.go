package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v4/x/grid/types"
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

// TODO fix empty case
//"grid" : {
//	"routes" : [{
//		"destination" : "",
//		"value" : [],
//		"source" : "",
//		"name" : ""
//	}],
//	"params" : {
//		"max_routes" : 16
//	}
//}

func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	routes := k.GetAllRoutes(ctx)

	return types.NewGenesisState(params, routes)
}
