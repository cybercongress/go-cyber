package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/deep-foundation/deep-chain/x/dmn/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
}

func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	params := k.GetParams(ctx)

	return types.NewGenesisState(params)
}
