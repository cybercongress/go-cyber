package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrNotEnoughBandwidth        = errorsmod.Register(ModuleName, 2, "not enough personal bandwidth")
	ErrExceededMaxBlockBandwidth = errorsmod.Register(ModuleName, 3, "exceeded max block bandwidth")
)
