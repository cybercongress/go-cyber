package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidAddress      = errorsmod.Register(ModuleName, 2, "invalid address")
	ErrExceededMaxThoughts = errorsmod.Register(ModuleName, 3, "exceeded max thoughts")
	ErrBadCallData         = errorsmod.Register(ModuleName, 4, "bad call data")
	ErrBadGasPrice         = errorsmod.Register(ModuleName, 5, "bad gas price")
	ErrBadTrigger          = errorsmod.Register(ModuleName, 6, "bad trigger")
	ErrBadName             = errorsmod.Register(ModuleName, 7, "bad name")
	ErrThoughtNotExist     = errorsmod.Register(ModuleName, 8, "thought does not exist")
	ErrConvertTrigger      = errorsmod.Register(ModuleName, 9, "cannot convert trigger")
)
