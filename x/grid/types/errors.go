package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrWrongName       = sdkerrors.Register(ModuleName, 2, "length of the name is not valid")
	ErrRouteNotExist   = sdkerrors.Register(ModuleName, 3, "the route does not exist")
	ErrRouteExist      = sdkerrors.Register(ModuleName, 4, "the route exists")
	ErrWrongValueDenom = sdkerrors.Register(ModuleName, 5, "the denom of value is not supported")
	ErrMaxRoutes       = sdkerrors.Register(ModuleName, 6, "max routes are exceeded")
	ErrSelfRoute       = sdkerrors.Register(ModuleName, 7, "routing to self is not allowed")
)
