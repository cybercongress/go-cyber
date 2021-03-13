package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrEmptySourceAddr      = sdkerrors.Register(ModuleName, 2, "empty source address")
	ErrEmptyDestinationAddr = sdkerrors.Register(ModuleName, 3, "empty destination address")
	ErrWrongAlias           = sdkerrors.Register(ModuleName, 4, "wrong alias")
	ErrReverseRoute			= sdkerrors.Register(ModuleName, 5, "no reverse routes")
	ErrRouteNotExist		= sdkerrors.Register(ModuleName, 6, "route not exist")
	ErrRouteExist			= sdkerrors.Register(ModuleName, 7, "route exist")
	ErrWrongDenom 			= sdkerrors.Register(ModuleName, 8, "wrong denom")
	ErrMaxRoutes 			= sdkerrors.Register(ModuleName, 9, "max routes exceeded")
	ErrSelfRoute			= sdkerrors.Register(ModuleName, 10, "self routes not allowed")
)
