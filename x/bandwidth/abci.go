package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	//"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cybercongress/cyberd/x/bandwidth/internal/types"
)

var accountsToUpdate = make([]sdk.AccAddress, 0)

// user to recover and update bandwidth for accounts with changed stake
func updateAccMaxBandwidth(ctx sdk.Context, meter types.BandwidthMeter) {
	for _, addr := range accountsToUpdate {
		meter.UpdateAccMaxBandwidth(ctx, addr)
	}
	accountsToUpdate = make([]sdk.AccAddress, 0)
}

// collect all addresses with updated stake
func CollectAddressesWithStakeChange() func(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress) {
	return func(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress) {
		if ctx.IsCheckTx() {
			return
		}
		if from != nil {
			accountsToUpdate = append(accountsToUpdate, from)
		}
		if to != nil {
			accountsToUpdate = append(accountsToUpdate, to)
		}
	}
}

// Used for 2 points:
// 1. Adjust credit price each `AdjustPricePeriod` blocks
// 2. For accounts with updated on current block stake adjust max bandwidth. Why not update on `onCoinsTransfer`?
//  Because for some bound related operations, coins already added/reduced from accounts, but not added to
//  validator\delegator pool.
func EndBlocker(ctx sdk.Context, bk AccountBandwidthKeeper, meter types.BandwidthMeter) {
	params := bk.GetParams(ctx)
	if ctx.BlockHeight() != 0 && ctx.BlockHeight()%params.AdjustPricePeriod == 0 {
		meter.AdjustPrice(ctx)
	}
	meter.CommitBlockBandwidth(ctx)
	updateAccMaxBandwidth(ctx, meter)
}
