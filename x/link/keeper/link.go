package keeper

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/store"
	. "github.com/cybercongress/cyberd/types"
	. "github.com/cybercongress/cyberd/x/link/types"
)

type LinkKeeper interface {
	PutLink(ctx sdk.Context, link Link)
	IsLinkExist(ctx sdk.Context, link Link) bool
	GetAllLinks(ctx sdk.Context) (Links, Links, error)
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

func (ls BaseLinkKeeper) PutLink(ctx sdk.Context, link Link) {
	store := ctx.KVStore(ls.key)
	linkAsBytes := marshalLink(link)
	store.Set(linkAsBytes, []byte{})
	ls.ms.IncrementLinksCount(ctx)
}

func (ls BaseLinkKeeper) IsLinkExist(ctx sdk.Context, link Link) bool {
	store := ctx.KVStore(ls.key)
	linkAsBytes := marshalLink(link)
	return store.Get(linkAsBytes) != nil
}

func (ls BaseLinkKeeper) GetAllLinks(ctx sdk.Context) (Links, Links, error) {

	inLinks := make(map[CidNumber]CidLinks)
	outLinks := make(map[CidNumber]CidLinks)

	store := ctx.KVStore(ls.key)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for iterator.Valid() {
		linkAsBytes := iterator.Key()
		link := unmarshalLink(linkAsBytes)
		Links(outLinks).Put(link.From(), link.To(), link.Acc())
		Links(inLinks).Put(link.To(), link.From(), link.Acc())
		iterator.Next()
	}
	return inLinks, outLinks, nil
}

func unmarshalLink(b []byte) Link {
	return NewLink(
		CidNumber(binary.LittleEndian.Uint64(b[0:8])),
		CidNumber(binary.LittleEndian.Uint64(b[8:16])),
		AccNumber(binary.LittleEndian.Uint64(b[16:24])),
	)
}

func marshalLink(l Link) []byte {
	b := make([]byte, 24)
	binary.LittleEndian.PutUint64(b[0:8], uint64(l.From()))
	binary.LittleEndian.PutUint64(b[8:16], uint64(l.To()))
	binary.LittleEndian.PutUint64(b[16:24], uint64(l.Acc()))
	return b
}
