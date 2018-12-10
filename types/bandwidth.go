package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgBandwidthCost func(msg sdk.Msg) int64

type BandwidthHandler func(ctx sdk.Context, msgCost MsgBandwidthCost, price float64, tx sdk.Tx) sdk.Error
