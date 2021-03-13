package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/bandwidth/keeper"
)

var accountsToUpdate = make([]sdk.AccAddress, 0)

// user to recover and update bandwidth for accounts with changed stake
func updateAccountMaxBandwidth(ctx sdk.Context, meter *keeper.BandwidthMeter) {
	for _, addr := range accountsToUpdate {
		meter.UpdateAccountMaxBandwidth(ctx, addr)
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
func EndBlocker(ctx sdk.Context, bm *keeper.BandwidthMeter) {

	params := bm.GetParams(ctx)
	if ctx.BlockHeight() != 0 && uint64(ctx.BlockHeight())%params.AdjustPricePeriod == 0 {
		bm.AdjustPrice(ctx)
	}

	bm.CommitBlockBandwidth(ctx)
	updateAccountMaxBandwidth(ctx, bm)

	// TODO Sandbbox to remove
	//store := ctx.KVStore(bm.StoreKey)
	//iterator := sdk.KVStorePrefixIterator(store, types.AccountStoreKeyPrefix)
	//defer iterator.Close()
	//
	//for iterator.Valid() {
	//	var ab types.AccountBandwidth
	//	bm.Cdc.UnmarshalBinaryBare(iterator.Value(), &ab)
	//	fmt.Println("Acc:", ab)
	//	iterator.Next()
	//}
}
