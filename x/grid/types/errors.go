package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrWrongName       = errorsmod.Register(ModuleName, 2, "length of the name is not valid")
	ErrRouteNotExist   = errorsmod.Register(ModuleName, 3, "the route does not exist")
	ErrRouteExist      = errorsmod.Register(ModuleName, 4, "the route exists")
	ErrWrongValueDenom = errorsmod.Register(ModuleName, 5, "the denom of value is not supported")
	ErrMaxRoutes       = errorsmod.Register(ModuleName, 6, "max routes are exceeded")
	ErrSelfRoute       = errorsmod.Register(ModuleName, 7, "routing to self is not allowed")
)
