package bandwidth

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v6/x/bandwidth/keeper"
	"github.com/cybercongress/go-cyber/v6/x/bandwidth/types"
)

func EndBlocker(ctx sdk.Context, bm *keeper.BandwidthMeter) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	params := bm.GetParams(ctx)
	if ctx.BlockHeight() != 0 && uint64(ctx.BlockHeight())%params.AdjustPricePeriod == 0 {
		bm.AdjustPrice(ctx) // TODO Add block event for price and load
	}

	bm.CommitBlockBandwidth(ctx) // TODO add block event for committed bandwidth
}
