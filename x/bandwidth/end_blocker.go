package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/app/storage"
	"math"
)

func EndBlocker(ctx sdk.Context, totalSpentForPeriod uint64, ms storage.MainStorage) uint64 {
	if ctx.BlockHeight() != 0 && ctx.BlockHeight()%AdjustPricePeriod == 0 {
		averageSpentValue := MaxNetworkBandwidth / (RecoveryPeriod / AdjustPricePeriod)
		newPrice := float64(totalSpentForPeriod) / float64(averageSpentValue)

		ms.StoreBandwidthPrice(ctx, math.Float64bits(newPrice))
		ms.StoreSpentBandwidth(ctx, 0)
		totalSpentForPeriod = 0
	}

	ms.StoreSpentBandwidth(ctx, totalSpentForPeriod)
	return totalSpentForPeriod
}
