package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"

	. "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/graph/types"
	"github.com/tendermint/tendermint/libs/log"

	"io"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	LinkBytesSize       = uint64(24)
	LinksCountBytesSize = uint64(8)
)

type GraphKeeper struct {
	key    sdk.StoreKey
	cdc    codec.BinaryMarshaler
	neudeg map[uint64]uint64
	rankNeudeg map[uint64]uint64
	tkey   sdk.StoreKey
}

func NewKeeper(
	cdc codec.BinaryMarshaler,
	storeKey sdk.StoreKey,
	tkey sdk.StoreKey,
) *GraphKeeper {
	return &GraphKeeper{
		cdc: cdc,
		key: storeKey,
		tkey: tkey,
		neudeg: make(map[uint64]uint64),
		rankNeudeg: make(map[uint64]uint64),
	}
}

func (k GraphKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ------------------------

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
	if neudeg == nil { return 0 }
	return sdk.BigEndianToUint64(neudeg)
}

func (gk GraphKeeper) GetTNeudeg(ctx sdk.Context, accNumber uint64) uint64 {
	store := ctx.TransientStore(gk.tkey)
	neudeg := store.Get(types.NeudegTStoreKey(accNumber))
	if neudeg == nil { return 0 }
	return sdk.BigEndianToUint64(neudeg)
}

func (gk *GraphKeeper) GetNeudegs() map[uint64]uint64 {
	return gk.rankNeudeg
}

func (gk *GraphKeeper) UpdateNeudegs(ctx sdk.Context) {
	iterator := sdk.KVStorePrefixIterator(ctx.TransientStore(gk.tkey), types.NeudegTStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		acc := sdk.BigEndianToUint64(iterator.Key()[1:])
		gk.neudeg[acc] = gk.neudeg[acc] + sdk.BigEndianToUint64(iterator.Value())
	}
}

func (gk *GraphKeeper) FixNeudegs() {
	for acc := range gk.neudeg {
		gk.rankNeudeg[acc] = gk.neudeg[acc]
	}
}
// ------------------------

func (gk GraphKeeper) SaveLink(ctx sdk.Context, link types.CompactLink) {
	defer telemetry.IncrCounter(1.0, types.ModuleName, "cyberlinks")

	store := ctx.KVStore(gk.key)
	store.Set(types.CyberlinksStoreKey(gk.GetLinksCount(ctx)), gk.cdc.MustMarshalBinaryBare(&link))

	gk.IncrementLinksCount(ctx)
}

func (gk GraphKeeper) GetAllLinks(ctx sdk.Context) (types.Links, types.Links, error) {
	return gk.GetAllLinksFiltered(ctx, func(l types.CompactLink) bool { return true })
}

func (gk GraphKeeper) GetAllLinksFiltered(ctx sdk.Context, filter types.LinkFilter) (types.Links, types.Links, error) {

	inLinks := make(map[types.CidNumber]types.CidLinks)
	outLinks := make(map[types.CidNumber]types.CidLinks)

	gk.IterateLinks(ctx, func(link types.CompactLink) {
		if filter(link) {
			types.Links(outLinks).Put(types.CidNumber(link.From), types.CidNumber(link.To), AccNumber(link.Account))
			types.Links(inLinks).Put(types.CidNumber(link.To), types.CidNumber(link.From), AccNumber(link.Account))
		}
	})

	return inLinks, outLinks, nil
}

func (gk GraphKeeper) IterateLinks(ctx sdk.Context, process func(link types.CompactLink)) {
	gk.IterateBinaryLinks(ctx, func(bytes []byte) {
		var compactLink types.CompactLink
		gk.cdc.MustUnmarshalBinaryBare(bytes, &compactLink)
		process(compactLink)
	})
}

func (gk GraphKeeper) IterateBinaryLinks(ctx sdk.Context, process func(bytes []byte)) {
	store := ctx.KVStore(gk.key)

	iterator := sdk.KVStorePrefixIterator(store, types.CyberlinkStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		process(iterator.Value())
	}
}

func (gk GraphKeeper) IncrementLinksCount(ctx sdk.Context) {
	store := ctx.KVStore(gk.key)

	linksCount := gk.GetLinksCount(ctx) + 1
	store.Set(types.LinksCount, sdk.Uint64ToBigEndian(linksCount))
}

func (gk GraphKeeper) GetLinksCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(gk.key)
	linksCountAsBytes := store.Get(types.LinksCount)

	if linksCountAsBytes == nil {
		return 0
	}
	return sdk.BigEndianToUint64(linksCountAsBytes)
}

// write links to writer in binary format: <links_count><cid_number_from><cid_number_to><acc_number>...
func (gk GraphKeeper) WriteLinks(ctx sdk.Context, writer io.Writer) (err error) {
	uintAsBytes := make([]byte, 8)
	linksCount := gk.GetLinksCount(ctx)
	binary.LittleEndian.PutUint64(uintAsBytes, linksCount)
	_, err = writer.Write(uintAsBytes)
	if err != nil {
		return
	}

	gk.IterateLinks(ctx, func(link types.CompactLink) {
		bytes := link.MarshalBinaryLink()
		_, err = writer.Write(bytes)
		if err != nil {
			return
		}
	})

	return nil
}