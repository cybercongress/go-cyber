package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/v2/x/graph/types"
)

// In order to calculate the flow of amperes through cyberlinks created by given agents than need to compute neurons out-degree
// transient store used to sync values to in-memory at the end block
func (gk *GraphKeeper) LoadNeudeg(rankCtx sdk.Context, freshCtx sdk.Context) {
	iterator := sdk.KVStorePrefixIterator(rankCtx.KVStore(gk.key), types.NeudegStoreKeyPrefix)
	for ; iterator.Valid(); iterator.Next() {
		acc := sdk.BigEndianToUint64(iterator.Key()[1:])
		gk.rankNeudeg[acc] = sdk.BigEndianToUint64(iterator.Value())
	}
	iterator.Close()

	iterator = sdk.KVStorePrefixIterator(freshCtx.KVStore(gk.key), types.NeudegStoreKeyPrefix)
	for ; iterator.Valid(); iterator.Next() {
		acc := sdk.BigEndianToUint64(iterator.Key()[1:])
		gk.neudeg[acc] = sdk.BigEndianToUint64(iterator.Value())
	}
	iterator.Close()
}

func (gk GraphKeeper) IncrementNeudeg(ctx sdk.Context, accNumber uint64) {
	store := ctx.KVStore(gk.key)
	neudeg := gk.GetNeudeg(ctx, accNumber) + 1
	store.Set(types.NeudegStoreKey(accNumber), sdk.Uint64ToBigEndian(neudeg))

	store = ctx.TransientStore(gk.tkey)
	tneudeg := gk.GetTNeudeg(ctx, accNumber) + 1
	store.Set(types.NeudegTStoreKey(accNumber), sdk.Uint64ToBigEndian(tneudeg))
}

func (gk GraphKeeper) GetNeudeg(ctx sdk.Context, accNumber uint64) uint64 {
	store := ctx.KVStore(gk.key)
	neudeg := store.Get(types.NeudegStoreKey(accNumber))
	if neudeg == nil {
		return 0
	}
	return sdk.BigEndianToUint64(neudeg)
}

func (gk GraphKeeper) GetTNeudeg(ctx sdk.Context, accNumber uint64) uint64 {
	store := ctx.TransientStore(gk.tkey)
	neudeg := store.Get(types.NeudegTStoreKey(accNumber))
	if neudeg == nil {
		return 0
	}
	return sdk.BigEndianToUint64(neudeg)
}

func (gk *GraphKeeper) GetNeudegs() map[uint64]uint64 {
	return gk.rankNeudeg
}

func (gk *GraphKeeper) UpdateMemNeudegs(ctx sdk.Context) {
	iterator := sdk.KVStorePrefixIterator(ctx.TransientStore(gk.tkey), types.NeudegTStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		acc := sdk.BigEndianToUint64(iterator.Key()[1:])
		gk.neudeg[acc] = gk.neudeg[acc] + sdk.BigEndianToUint64(iterator.Value())
	}
}

func (gk *GraphKeeper) UpdateRankNeudegs() {
	for acc := range gk.neudeg {
		gk.rankNeudeg[acc] = gk.neudeg[acc]
	}
}
