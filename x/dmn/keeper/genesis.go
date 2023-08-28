package keeper

import (
	"github.com/cybercongress/go-cyber/x/dmn/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
}

func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	params := k.GetParams(ctx)

	return types.NewGenesisState(params)
}
