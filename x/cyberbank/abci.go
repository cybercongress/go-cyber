package cyberbank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/cyberbank/keeper"
)

func EndBlocker(ctx sdk.Context, k *keeper.IndexedKeeper) {
	k.EndBlocker(ctx)
}