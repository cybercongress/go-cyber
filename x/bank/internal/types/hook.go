package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// `from` and `to` can be nill
type CoinsTransferHook = func(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress)
