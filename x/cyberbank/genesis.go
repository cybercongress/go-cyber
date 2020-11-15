package cyberbank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/cyberbank/keeper"
)

func InitGenesis(ctx sdk.Context, k *keeper.IndexedKeeper) {
	k.InitGenesis(ctx)
}

