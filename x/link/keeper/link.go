package keeper

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/store"
	"github.com/cybercongress/cyberd/util"
	. "github.com/cybercongress/cyberd/x/link/types"
	"io"
)

const (
	LinkBytesSize       = uint64(24)
	LinksCountBytesSize = uint64(8)
)

type LinkFilter func(l CompactLink) bool

var DefaultLinkFilter = func(l CompactLink) bool { return true }

type LinkKeeper interface {
	PutLink(ctx sdk.Context, link CompactLink)
	IsLinkExist(ctx sdk.Context, link CompactLink) bool
	GetAllLinks(ctx sdk.Context) (Links, Links, error)
	GetAllLinksFiltered(ctx sdk.Context, filter LinkFilter) (Links, Links, error)
	GetLinksCount(ctx sdk.Context) uint64
	Iterate(ctx sdk.Context, process func(link CompactLink))
	WriteLinks(ctx sdk.Context, writer io.Writer) (err error)
	LoadFromReader(ctx sdk.Context, reader io.Reader) (err error)
}

type BaseLinkKeeper struct {
	ms  store.MainKeeper
	key *sdk.KVStoreKey
}

func NewBaseLinkKeeper(ms store.MainKeeper, key *sdk.KVStoreKey) LinkKeeper {
	return BaseLinkKeeper{
		key: key,
		ms:  ms,
	}
}

func (lk BaseLinkKeeper) PutLink(ctx sdk.Context, link CompactLink) {
	store := ctx.KVStore(lk.key)
	linkAsBytes := link.MarshalBinary()
	store.Set(linkAsBytes, []byte{})
	lk.ms.IncrementLinksCount(ctx)
}

func (lk BaseLinkKeeper) IsLinkExist(ctx sdk.Context, link CompactLink) bool {
	store := ctx.KVStore(lk.key)
	linkAsBytes := link.MarshalBinary()
	return store.Get(linkAsBytes) != nil
}

func (lk BaseLinkKeeper) GetAllLinks(ctx sdk.Context) (Links, Links, error) {
	return lk.GetAllLinksFiltered(ctx, DefaultLinkFilter)
}

func (lk BaseLinkKeeper) GetAllLinksFiltered(ctx sdk.Context, filter LinkFilter) (Links, Links, error) {

	inLinks := make(map[CidNumber]CidLinks)
	outLinks := make(map[CidNumber]CidLinks)

	lk.Iterate(ctx, func(link CompactLink) {
		if filter(link) {
			Links(outLinks).Put(link.From(), link.To(), link.Acc())
			Links(inLinks).Put(link.To(), link.From(), link.Acc())
		}
	})

	return inLinks, outLinks, nil
}

func (lk BaseLinkKeeper) GetLinksCount(ctx sdk.Context) uint64 {
	return lk.ms.GetLinksCount(ctx)
}

func (lk BaseLinkKeeper) Iterate(ctx sdk.Context, process func(link CompactLink)) {
	iterator := ctx.KVStore(lk.key).Iterator(nil, nil)
	defer iterator.Close()

	for iterator.Valid() {
		process(UnmarshalBinaryLink(iterator.Key()))
		iterator.Next()
	}
}

// write links to writer in binary format: <links_count><cid_number_from><cid_number_to><acc_number>...
func (lk BaseLinkKeeper) WriteLinks(ctx sdk.Context, writer io.Writer) (err error) {
	uintAsBytes := make([]byte, 8) //common bytes array to convert uints

	linksCount := lk.GetLinksCount(ctx)
	binary.LittleEndian.PutUint64(uintAsBytes, linksCount)
	_, err = writer.Write(uintAsBytes)
	if err != nil {
		return
	}

	lk.Iterate(ctx, func(link CompactLink) {
		_, _ = writer.Write(link.MarshalBinary())
	})
	return
}

func (lk BaseLinkKeeper) LoadFromReader(ctx sdk.Context, reader io.Reader) (err error) {
	linksCountBytes, err := util.ReadExactlyNBytes(reader, LinksCountBytesSize)
	if err != nil {
		return
	}
	linksCount := binary.LittleEndian.Uint64(linksCountBytes)

	for j := uint64(0); j < linksCount; j++ {
		linkBytes, err := util.ReadExactlyNBytes(reader, LinkBytesSize)
		if err != nil {
			return err
		}
		lk.PutLink(ctx, UnmarshalBinaryLink(linkBytes))
	}
	return
}
