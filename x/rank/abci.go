package rank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v5/x/rank/keeper"
)

func EndBlocker(ctx sdk.Context, k *keeper.StateKeeper) {
	k.EndBlocker(ctx)
}
