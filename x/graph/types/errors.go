package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrCyberlinkExist  = errorsmod.Register(ModuleName, 2, "your cyberlink already exists")
	ErrZeroLinks       = errorsmod.Register(ModuleName, 3, "cyberlinks not found")
	ErrSelfLink        = errorsmod.Register(ModuleName, 4, "loop cyberlink not allowed")
	ErrInvalidParticle = errorsmod.Register(ModuleName, 5, "invalid particle")
	ErrCidNotFound     = errorsmod.Register(ModuleName, 6, "particle not found")
	ErrCidVersion      = errorsmod.Register(ModuleName, 7, "unsupported cid version")
	ErrZeroPower       = errorsmod.Register(ModuleName, 8, "neuron has zero power")
)
