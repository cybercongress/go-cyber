package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNotEnoughBandwidth        = sdkerrors.Register(ModuleName, 2, "not enough personal bandwidth")
	ErrExceededMaxBlockBandwidth = sdkerrors.Register(ModuleName, 3, "exceeded max block bandwidth")
)
