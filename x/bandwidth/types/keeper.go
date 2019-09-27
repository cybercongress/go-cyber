package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper interface {
	SetAccBandwidth(ctx sdk.Context, bandwidth AcсBandwidth)
	GetAccBandwidth(ctx sdk.Context, address sdk.AccAddress) (bw AcсBandwidth)

	SetParams(ctx sdk.Context, params Params)
	GetParams(ctx sdk.Context) (params Params)
}

type BlockSpentBandwidthKeeper interface {
	SetBlockSpentBandwidth(ctx sdk.Context, blockNumber uint64, value uint64)
	GetValuesForPeriod(ctx sdk.Context, period int64) map[uint64]uint64
}
