package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrCyberlinkExist 	= sdkerrors.Register(ModuleName, 2, "your cyberlink already exists")
	ErrInvalidCid 		= sdkerrors.Register(ModuleName, 3, "invalid cid")
	ErrZeroLinks 		= sdkerrors.Register(ModuleName, 4, "no links found")
	ErrCidNotFound 		= sdkerrors.Register(ModuleName, 5, "cid not found")
	ErrInvalidData 		= sdkerrors.Register(ModuleName, 6, "invalid data")
	ErrInvalidAccount 	= sdkerrors.Register(ModuleName, 7, "invalid account")
	ErrSelfLink 		= sdkerrors.Register(ModuleName, 8, "self cyberlink")
	ErrZeroPower 		= sdkerrors.Register(ModuleName, 9, "zero power")
)