package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v4/x/rank/types"
)

func InitGenesis(ctx sdk.Context, sk StateKeeper, data types.GenesisState) {
	if err := sk.SetParams(ctx, data.Params); err != nil {
		panic(err)
	}
}

func ExportGenesis(ctx sdk.Context, keeper StateKeeper) *types.GenesisState {
	return types.NewGenesisState(keeper.GetParams(ctx))
}
