package storage

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (ls LinksStorage) AddLink(ctx sdk.Context, link LinkedCids) error {

	store := ctx.KVStore(ls.key)

	linkAsBytes, err := ls.cdc.MarshalBinaryBare(link)
	if err != nil {
		return err
	}

	store.Set(linkAsBytes, []byte{})
	return nil
}

func (ls LinksStorage) GetAllLinks(ctx sdk.Context) (map[CidNumber]CidLinks, map[CidNumber]CidLinks, error) {

	inLinks := make(map[CidNumber]CidLinks)
	outLinks := make(map[CidNumber]CidLinks)

	store := ctx.KVStore(ls.key)
	iterator := store.Iterator(nil, nil)

	var link LinkedCids
	for iterator.Valid() {
		linkAsBytes := iterator.Key()
		err := ls.cdc.UnmarshalBinaryBare(linkAsBytes, &link)
		if err != nil {
			iterator.Close()
			return nil, nil, err
		}

		CidsLinks(outLinks).Put(link.FromCid, link.ToCid, link.Creator)
		CidsLinks(inLinks).Put(link.ToCid, link.FromCid, link.Creator)
		iterator.Next()
	}
	iterator.Close()
	return inLinks, outLinks, nil
}
