package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/store"
	. "github.com/cybercongress/cyberd/x/link/types"
)

type LinkKeeper interface {
	PutLink(ctx sdk.Context, link CompactLink)
	IsLinkExist(ctx sdk.Context, link CompactLink) bool
	GetAllLinks(ctx sdk.Context) (Links, Links, error)
	GetLinksCount(ctx sdk.Context) uint64
	Iterate(ctx sdk.Context, process func(link CompactLink))
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

	inLinks := make(map[CidNumber]CidLinks)
	outLinks := make(map[CidNumber]CidLinks)

	lk.Iterate(ctx, func(link CompactLink) {
		Links(outLinks).Put(link.From(), link.To(), link.Acc())
		Links(inLinks).Put(link.To(), link.From(), link.Acc())
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

