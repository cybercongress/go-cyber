package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/joinresistance/space-pussy/x/rank/types"
)

func InitGenesis(ctx sdk.Context, keeper StateKeeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)
}

func ExportGenesis(ctx sdk.Context, keeper StateKeeper) *types.GenesisState {
	return types.NewGenesisState(keeper.GetParams(ctx))
}
