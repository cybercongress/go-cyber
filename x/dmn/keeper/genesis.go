package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v5/x/dmn/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	err := k.SetParams(ctx, data.Params)
	if err != nil {
		panic(err)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	params := k.GetParams(ctx)

	return types.NewGenesisState(params)
}
