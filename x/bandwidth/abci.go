package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/app/storage"
	"math"
)

// Used for 2 points:
// 1. Adjust credit price each `AdjustPricePeriod` blocks
// 2. For accs with updated on current block stake adjust max bandwidth. Why not update on `onCoinsTransfer`?
//  Cuz for some bound related operations, coins already added/reduced from acc, but not added to
//  validator\delegator pool.
func EndBlocker(
	ctx sdk.Context, totalSpentForPeriod uint64, curPrice float64, ms storage.MainStorage,
) (newPrice float64, totalSpent uint64) {

	newPrice = curPrice
	if ctx.BlockHeight() != 0 && ctx.BlockHeight()%AdjustPricePeriod == 0 {
		averageSpentValue := MaxNetworkBandwidth / (RecoveryPeriod / AdjustPricePeriod)
		newPrice = float64(totalSpentForPeriod) / float64(averageSpentValue)

		if newPrice < 0.1*float64(BaseCreditPrice) {
			newPrice = 0.1 * float64(BaseCreditPrice)
		}

		ms.StoreBandwidthPrice(ctx, math.Float64bits(newPrice))
		ms.StoreSpentBandwidth(ctx, 0)
		totalSpentForPeriod = 0
	}

	ms.StoreSpentBandwidth(ctx, totalSpentForPeriod)
	return newPrice, totalSpentForPeriod
}
