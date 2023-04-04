package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrCyberlinkExist  = sdkerrors.Register(ModuleName, 2, "your cyberlink already exists")
	ErrZeroLinks       = sdkerrors.Register(ModuleName, 3, "cyberlinks not found")
	ErrSelfLink        = sdkerrors.Register(ModuleName, 4, "loop cyberlink not allowed")
	ErrInvalidParticle = sdkerrors.Register(ModuleName, 5, "invalid particle")
	ErrCidNotFound     = sdkerrors.Register(ModuleName, 6, "particle not found")
	ErrCidVersion      = sdkerrors.Register(ModuleName, 7, "unsupported cid version")
	ErrZeroPower       = sdkerrors.Register(ModuleName, 8, "neuron has zero power")
)
