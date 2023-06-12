package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalidAddress      = sdkerrors.Register(ModuleName, 2, "invalid address")
	ErrExceededMaxThoughts = sdkerrors.Register(ModuleName, 3, "exceeded max thoughts")
	ErrBadCallData         = sdkerrors.Register(ModuleName, 4, "bad call data")
	ErrBadGasPrice         = sdkerrors.Register(ModuleName, 5, "bad gas price")
	ErrBadTrigger          = sdkerrors.Register(ModuleName, 6, "bad trigger")
	ErrBadName         	   = sdkerrors.Register(ModuleName, 7, "bad name")
	ErrThoughtNotExist     = sdkerrors.Register(ModuleName, 8, "thought does not exist")
	ErrConvertTrigger      = sdkerrors.Register(ModuleName, 9, "cannot convert trigger")
)
