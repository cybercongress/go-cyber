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
	GetAllLinks(ctx sdk.Context) (Links, Links, error)
	GetAllLinksFiltered(ctx sdk.Context, filter LinkFilter) (Links, Links, error)
	GetLinksCount(ctx sdk.Context) uint64
	Iterate(ctx sdk.Context, process func(link CompactLink))
	WriteLinks(ctx sdk.Context, writer io.Writer) (err error)
	LoadFromReader(ctx sdk.Context, reader io.Reader) (err error)
	Commit(blockHeight uint64) (err error)
}

type BaseLinkKeeper struct {
	ms      store.MainKeeper
	key     *sdk.KVStoreKey
	storage store.Storage
}

func NewBaseLinkKeeper(ms store.MainKeeper, key *sdk.KVStoreKey) LinkKeeper {
	storage, err := store.NewBaseStorage(key.Name(), util.RootifyPath("data/"), LinkBytesSize)
	if err != nil {
		panic("Failed to load links DB")
	}
	return BaseLinkKeeper{
		key:     key,
		ms:      ms,
		storage: storage,
	}
}

func (lk BaseLinkKeeper) PutLink(ctx sdk.Context, link CompactLink) {
	if ctx.BlockHeight() == lk.storage.LastVersion() {
		return
	}
	linkAsBytes := link.MarshalBinary()
	err := lk.storage.Put(linkAsBytes)
	if err != nil {
		panic(err)
	}
	lk.ms.IncrementLinksCount(ctx)
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

	err := lk.storage.IterateTillVersion(func(bytes []byte) {
		process(UnmarshalBinaryLink(bytes))
	}, ctx.BlockHeight())

	if err != nil {
		panic("Error during links iterating")
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

func (lk BaseLinkKeeper) Commit(blockHeight uint64) error {
	return lk.storage.Commit(blockHeight)
}
