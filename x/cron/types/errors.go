package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalidAddress = sdkerrors.Register(ModuleName, 2, "invalid address")
	ErrEmptyContract  = sdkerrors.Register(ModuleName, 3, "empty contract")
	ErrBadCallData    = sdkerrors.Register(ModuleName, 4, "bad call data")
	ErrBadGasPrice    = sdkerrors.Register(ModuleName, 5, "bad gas price")
	ErrBadTrigger     = sdkerrors.Register(ModuleName, 6, "bad trigger")
	ErrBadLabel       = sdkerrors.Register(ModuleName, 7, "bad label")
	ErrJobNotExist    = sdkerrors.Register(ModuleName, 8, "job not exist")
	ErrNotAuthorized  = sdkerrors.Register(ModuleName, 9, "not authorized")
	ErrConvertTrigger  = sdkerrors.Register(ModuleName, 10, "cannot convert trigger")
	ErrExceededMaxJobs  = sdkerrors.Register(ModuleName, 11, "exceeded max jobs")
	ErrInvalidRoute  = sdkerrors.Register(ModuleName, 12, "invalid route")
)
