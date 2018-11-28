package storage

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cbd "github.com/cybercongress/cyberd/cosmos/poc/app/types"
)

type LinksStorage struct {
	ms  MainStorage
	key *sdk.KVStoreKey
}

func NewLinksStorage(ms MainStorage, key *sdk.KVStoreKey) LinksStorage {
	return LinksStorage{
		key: key,
		ms:  ms,
	}
}

func (ls LinksStorage) AddLink(ctx sdk.Context, link cbd.Link) {
	store := ctx.KVStore(ls.key)
	linkAsBytes := marshalLink(link)
	store.Set(linkAsBytes, []byte{})
	ls.ms.IncrementLinksCount(ctx)
}

func (ls LinksStorage) IsLinkExist(ctx sdk.Context, link cbd.Link) bool {
	store := ctx.KVStore(ls.key)
	linkAsBytes := marshalLink(link)
	return store.Get(linkAsBytes) != nil
}

func (ls LinksStorage) GetAllLinks(ctx sdk.Context) (cbd.Links, cbd.Links, error) {

	inLinks := make(map[cbd.CidNumber]cbd.CidLinks)
	outLinks := make(map[cbd.CidNumber]cbd.CidLinks)

	store := ctx.KVStore(ls.key)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for iterator.Valid() {
		linkAsBytes := iterator.Key()
		link := unmarshalLink(linkAsBytes)
		cbd.Links(outLinks).Put(link.From(), link.To(), link.Acc())
		cbd.Links(inLinks).Put(link.To(), link.From(), link.Acc())
		iterator.Next()
	}
	return inLinks, outLinks, nil
}

func unmarshalLink(b []byte) cbd.Link {
	return cbd.NewLink(
		cbd.CidNumber(binary.LittleEndian.Uint64(b[0:8])),
		cbd.CidNumber(binary.LittleEndian.Uint64(b[8:16])),
		cbd.AccountNumber(binary.LittleEndian.Uint64(b[16:24])),
	)
}

func marshalLink(l cbd.Link) []byte {
	b := make([]byte, 24)
	binary.LittleEndian.PutUint64(b[0:8], uint64(l.From()))
	binary.LittleEndian.PutUint64(b[8:16], uint64(l.To()))
	binary.LittleEndian.PutUint64(b[16:24], uint64(l.Acc()))
	return b
}
