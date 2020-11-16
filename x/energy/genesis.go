package energy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/x/energy/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
	k.SetRoutes(ctx, data.Routes)
}

func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {
	params := k.GetParams(ctx)
	routes := k.GetAllRoutes(ctx)

	return types.NewGenesisState(params, routes)
}
