package liquidity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v6/x/liquidity/keeper"
	"github.com/cybercongress/go-cyber/v6/x/liquidity/types"
)

// InitGenesis new liquidity genesis
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	keeper.InitGenesis(ctx, data)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	return keeper.ExportGenesis(ctx)
}
