package rank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState(keeper.GetParams(ctx))
}
