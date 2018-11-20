package storage

import (
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/cybercongress/cyberd/cosmos/poc/app/types"
)

type LinksStorage struct {
	cdc *codec.Codec
	key *sdk.KVStoreKey
}

func NewLinksStorage(key *sdk.KVStoreKey, cdc *codec.Codec) LinksStorage {
	return LinksStorage{
		key: key,
		cdc: cdc,
	}
}

func (ls LinksStorage) AddLink(ctx sdk.Context, link Link) {
	store := ctx.KVStore(ls.key)
	linkAsBytes := marshalLink(link)
	store.Set(linkAsBytes, []byte{})
}

func (ls LinksStorage) IsLinkExist(ctx sdk.Context, link Link) bool {
	store := ctx.KVStore(ls.key)
	linkAsBytes := marshalLink(link)
	return store.Get(linkAsBytes) != nil
}

func (ls LinksStorage) GetAllLinks(ctx sdk.Context) (Links, Links, error) {

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
		AccountNumber(binary.LittleEndian.Uint64(b[16:24])),
	)
}

func marshalLink(l Link) []byte {
	b := make([]byte, 24)
	binary.LittleEndian.PutUint64(b[0:8], uint64(l.From()))
	binary.LittleEndian.PutUint64(b[8:16], uint64(l.To()))
	binary.LittleEndian.PutUint64(b[16:24], uint64(l.Acc()))
	return b
}
