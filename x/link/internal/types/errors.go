package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// NOTE ABCI code 1 is unique. https://github.com/cosmos/cosmos-sdk/issues/5662
	ErrCyberlinkExist = sdkerrors.Register(ModuleName, 2, "your cyberlink already exists")
	ErrInvalidCid = sdkerrors.Register(ModuleName, 3, "invalid cid")
	ErrZeroLinks = sdkerrors.Register(ModuleName, 4, "no links found")
	//ErrCidNotFound = sdkerrors.Register(ModuleName, 5, "cid not found")
)