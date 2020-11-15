package keeper

import (
	"encoding/binary"

	"github.com/cybercongress/go-cyber/x/link/types"

	"io"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	LinkBytesSize       = uint64(24)
	LinksCountBytesSize = uint64(8)
)

type GraphKeeper struct {
	key  sdk.StoreKey
}

func NewKeeper(storeKey sdk.StoreKey) GraphKeeper {
	return GraphKeeper{
		key: storeKey,
	}
}

func (gk GraphKeeper) PutLink(ctx sdk.Context, link types.CompactLink) {
	linkAsBytes := link.MarshalBinary()

	store := ctx.KVStore(gk.key)
	store.Set(types.CyberlinksStoreKey(gk.GetLinksCount(ctx)), linkAsBytes)

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
			types.Links(outLinks).Put(link.From(), link.To(), link.Account())
			types.Links(inLinks).Put(link.To(), link.From(), link.Account())
		}
	})

	return inLinks, outLinks, nil
}

func (gk GraphKeeper) IterateLinks(ctx sdk.Context, process func(link types.CompactLink)) {
	gk.IterateBinaryLinks(ctx, func(bytes []byte) {
		process(types.UnmarshalBinaryLink(bytes))
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
	return types.BigEndianToUint64(linksCountAsBytes)
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

	gk.IterateBinaryLinks(ctx, func(bytes []byte) {
		//fmt.Println("link: ", types.UnmarshalBinaryLink(bytes))
		_, err = writer.Write(bytes)
		if err != nil {
			return
		}
	})

	return nil
}