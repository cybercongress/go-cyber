package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type CoinsTransferHook = func(ctx sdk.Context, from, to sdk.AccAddress)
