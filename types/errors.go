package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const CyberCodespace = "cyber"

var (
	ErrCyberlinkExist = sdkerrors.Register(CyberCodespace, 1, "your cyberlink already exists")
	ErrInvalidCid = sdkerrors.Register(CyberCodespace, 2, "invalid cid")
	ErrCidNotFound = sdkerrors.Register(CyberCodespace, 3, "cid not found")
	ErrNotEnoughBandwidth = sdkerrors.Register(CyberCodespace, 4, "not enough personal bandwidth")
	ErrZeroLinks = sdkerrors.Register(CyberCodespace, 5, "no links found")
	ErrExceededMaxBlockBandwidth = sdkerrors.Register(CyberCodespace, 6, "exceeded max block bandwidth")
)