package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// NOTE ABCI code 1 is unique. https://github.com/cosmos/cosmos-sdk/issues/5662
	ErrNotEnoughBandwidth = sdkerrors.Register(ModuleName, 2, "not enough personal bandwidth")
	ErrExceededMaxBlockBandwidth = sdkerrors.Register(ModuleName, 3, "exceeded max block bandwidth")
)