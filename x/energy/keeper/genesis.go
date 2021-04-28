package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/x/energy/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
	err := k.SetRoutes(ctx, data.Routes); if err != nil {
		panic(err)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	routes := k.GetAllRoutes(ctx)

	return types.NewGenesisState(params, routes)
}
