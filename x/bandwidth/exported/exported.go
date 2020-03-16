package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	types "github.com/cybercongress/go-cyber/x/bandwidth/internal/types"
)

type BaseAccountBandwidthKeeper interface {
	SetAccountBandwidth(ctx sdk.Context, bandwidth types.AcсountBandwidth)
	GetAccountBandwidth(ctx sdk.Context, address sdk.AccAddress) (bw types.AcсountBandwidth)

	SetParams(ctx sdk.Context, params types.Params)
	GetParams(ctx sdk.Context) (params types.Params)
}

type BaseBlockSpentBandwidthKeeper interface {
	SetBlockSpentBandwidth(ctx sdk.Context, blockNumber uint64, value uint64)
	GetValuesForPeriod(ctx sdk.Context, period int64) map[uint64]uint64
}

type Meter interface {
	GetCurrentCreditPrice() float64
}