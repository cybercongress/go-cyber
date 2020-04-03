package keeper

import (
	"bytes"
	"sync"

	"github.com/cybercongress/go-cyber/store"
	"github.com/cybercongress/go-cyber/x/link/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"encoding/binary"
	"io"
)

const (
	LinkBytesSize       = uint64(24)
	LinksCountBytesSize = uint64(8)
)

const defaultBufferSize = 65536

var DefaultLinkFilter = func(l types.CompactLink) bool { return true }

type Keeper struct {
	ms      	store.MainKeeper
	storeKey    sdk.StoreKey
	buffer      *bytes.Buffer
	mu 			*sync.Mutex
}

func NewLinkKeeper(ms store.MainKeeper, storeKey sdk.StoreKey) *Keeper {
	return &Keeper{
		storeKey: storeKey,
		ms:       ms,
		buffer:   bytes.NewBuffer(make([]byte, 0, defaultBufferSize)),
		mu:		  new(sync.Mutex),
	}
}

func (lk Keeper) PutLink(ctx sdk.Context, link types.CompactLink) {
	lk.mu.Lock()
	defer lk.mu.Unlock()
	linkAsBytes := link.MarshalBinary()
	if uint64(len(linkAsBytes)) != LinkBytesSize {
		panic("invalid element length")
	}
	lk.buffer.Write(linkAsBytes)
	lk.ms.IncrementLinksCount(ctx)
}

func (lk Keeper) GetAllLinks(ctx sdk.Context) (types.Links, types.Links, error) {
	return lk.GetAllLinksFiltered(ctx, DefaultLinkFilter)
}

func (lk Keeper) GetAllLinksFiltered(ctx sdk.Context, filter types.LinkFilter) (types.Links, types.Links, error) {

	inLinks := make(map[types.CidNumber]types.CidLinks)
	outLinks := make(map[types.CidNumber]types.CidLinks)

	lk.IterateLinks(ctx, func(link types.CompactLink) {
		if filter(link) {
			types.Links(outLinks).Put(link.From(), link.To(), link.Acc())
			types.Links(inLinks).Put(link.To(), link.From(), link.Acc())
		}
	})

	return inLinks, outLinks, nil
}

func (lk Keeper) GetLinksCount(ctx sdk.Context) uint64 {
	return lk.ms.GetLinksCount(ctx)
}

func (lk Keeper) IterateLinks(ctx sdk.Context, process func(link types.CompactLink)) {
	lk.IterateBinaryLinks(ctx, func(bytes []byte) {
		process(types.UnmarshalBinaryLink(bytes))
	})
}

func (lk Keeper) IterateBinaryLinks(ctx sdk.Context, process func(bytes []byte)) {
	store := ctx.KVStore(lk.storeKey)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		value := iterator.Value()
		if (len(value) == 0) { continue }

		links := len(value)/int(LinkBytesSize)

		for i := 0 ; i < links; i++ {
			elementBytes := value[:LinkBytesSize]
			value = value[LinkBytesSize:]
			process(elementBytes)
		}
	}
}

// write links to writer in binary format: <links_count><cid_number_from><cid_number_to><acc_number>...
func (lk Keeper) WriteLinks(ctx sdk.Context, writer io.Writer) (err error) {
	uintAsBytes := make([]byte, 8)
	linksCount := lk.GetLinksCount(ctx)
	binary.LittleEndian.PutUint64(uintAsBytes, linksCount)
	_, err = writer.Write(uintAsBytes)
	if err != nil {
		return
	}

	lk.IterateBinaryLinks(ctx, func(bytes []byte) {
		_, err = writer.Write(bytes)
		if err != nil {
			return
		}
	})

	return nil
}

func (lk Keeper) Commit(ctx sdk.Context) {
	lk.mu.Lock()
	defer func() {
		lk.mu.Unlock()
		lk.buffer.Reset()
	}()

	if lk.buffer.Len() > 0 {
		versionAsBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(versionAsBytes, uint64(ctx.BlockHeight()))
		store := ctx.KVStore(lk.storeKey)

		store.Set(versionAsBytes, lk.buffer.Bytes())
	}

}
