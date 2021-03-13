package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, k *IndexedKeeper) {
	k.InitGenesis(ctx)
}

