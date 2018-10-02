package storage

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type LinksStorage struct {
	cdc *wire.Codec
	key *sdk.KVStoreKey
}

func NewLinksStorage(key *sdk.KVStoreKey, cdc *wire.Codec) LinksStorage {
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

		//
		// out links
		cidLinks := outLinks[link.FromCid]
		if cidLinks == nil {
			cidLinks = make(CidLinks)
		}
		users := cidLinks[link.ToCid]
		if users == nil {
			users = make(map[AccountNumber]struct{})
		}
		users[link.Creator] = struct{}{}
		cidLinks[link.ToCid] = users
		outLinks[link.FromCid] = cidLinks

		//
		// in links
		cidLinks = inLinks[link.ToCid]
		if cidLinks == nil {
			cidLinks = make(CidLinks)
		}
		users = cidLinks[link.FromCid]
		if users == nil {
			users = make(map[AccountNumber]struct{})
		}
		users[link.Creator] = struct{}{}
		cidLinks[link.FromCid] = users
		inLinks[link.ToCid] = cidLinks

		iterator.Next()
	}
	iterator.Close()
	return inLinks, outLinks, nil
}
