package keeper

import (
	"github.com/cybercongress/go-cyber/x/rank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, keeper StateKeeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)
}

func ExportGenesis(ctx sdk.Context, keeper StateKeeper) *types.GenesisState {
	return types.NewGenesisState(keeper.GetParams(ctx))
}
