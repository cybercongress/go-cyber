package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper interface {
	SetAccBandwidth(ctx sdk.Context, bandwidth AcсBandwidth)
	GetAccBandwidth(ctx sdk.Context, address sdk.AccAddress) (bw AcсBandwidth)
}
