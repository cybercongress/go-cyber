package keeper

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"

	ctypes "github.com/cybercongress/go-cyber/v2/types"
	"github.com/cybercongress/go-cyber/v2/x/graph/types"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	LinkBytesSize       = uint64(24)
	LinksCountBytesSize = uint64(8)
)

type GraphKeeper struct {
	key        sdk.StoreKey
	cdc        codec.BinaryCodec
	neudeg     map[uint64]uint64
	rankNeudeg map[uint64]uint64
	tkey       sdk.StoreKey
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	tkey sdk.StoreKey,
) *GraphKeeper {
	return &GraphKeeper{
		cdc:        cdc,
		key:        storeKey,
		tkey:       tkey,
		neudeg:     make(map[uint64]uint64),
		rankNeudeg: make(map[uint64]uint64),
	}
}

func (k GraphKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (gk GraphKeeper) SaveLink(ctx sdk.Context, link types.CompactLink) {
	defer telemetry.IncrCounter(1.0, types.ModuleName, "cyberlinks")

	store := ctx.KVStore(gk.key)
	store.Set(types.CyberlinksStoreKey(types.CyberlinkRawKey(link)), sdk.Uint64ToBigEndian(uint64(ctx.BlockHeight())))

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
			types.Links(outLinks).Put(types.CidNumber(link.From), types.CidNumber(link.To), ctypes.AccNumber(link.Account))
			types.Links(inLinks).Put(types.CidNumber(link.To), types.CidNumber(link.From), ctypes.AccNumber(link.Account))
		}
	})

	return inLinks, outLinks, nil
}

func (gk GraphKeeper) IterateLinks(ctx sdk.Context, process func(link types.CompactLink)) {
	gk.IterateBinaryLinks(ctx, func(key, value []byte) {
		compactLink := types.CompactLink{
			// because first byte is prefix
			From:    sdk.BigEndianToUint64(key[1:9]),
			To:      sdk.BigEndianToUint64(key[17:25]),
			Account: sdk.BigEndianToUint64(key[9:17]),
		}
		process(compactLink)
	})
}

func (gk GraphKeeper) IterateBinaryLinks(ctx sdk.Context, process func(key, value []byte)) {
	store := ctx.KVStore(gk.key)

	iterator := sdk.KVStorePrefixIterator(store, types.CyberlinkStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		process(iterator.Key(), iterator.Value())
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
