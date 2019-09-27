package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper interface {
	SetParams(ctx sdk.Context, params Params)
	GetParams(ctx sdk.Context) (params Params)
}
