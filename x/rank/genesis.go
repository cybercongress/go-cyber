package rank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/cyberd/x/rank/exported"
)

func InitGenesis(ctx sdk.Context, keeper exported.StateKeeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
}

func ExportGenesis(ctx sdk.Context, keeper exported.StateKeeper) GenesisState {
	return NewGenesisState(keeper.GetParams(ctx))
}
