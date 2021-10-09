package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrWrongAlias           = sdkerrors.Register(ModuleName, 2, "length of the alias isn't valid")
	ErrRouteNotExist		= sdkerrors.Register(ModuleName, 3, "the route isn't exist")
	ErrRouteExist      		= sdkerrors.Register(ModuleName, 4, "the route is exist")
	ErrWrongValueDenom 		= sdkerrors.Register(ModuleName, 5, "the denom of value isn't supported")
	ErrMaxRoutes       		= sdkerrors.Register(ModuleName, 6, "max routes are exceeded")
	ErrSelfRoute			= sdkerrors.Register(ModuleName, 7, "routing to self is not allowed")
)
