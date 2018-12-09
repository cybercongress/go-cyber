package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BandwidthHandler func(ctx sdk.Context, tx sdk.Tx) sdk.Error
