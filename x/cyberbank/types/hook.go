package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type CoinsTransferHook = func(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress)
