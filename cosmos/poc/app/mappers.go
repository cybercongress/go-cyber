package app

import (
	"encoding/binary"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

var cidsCount = []byte("cids_count")

type LinksStorage struct {
	cdc *wire.Codec
	key sdk.StoreKey
}

func (iom LinksStorage) AddLink(ctx sdk.Context, user sdk.AccAddress, cid string, linkedCid string) error {
	store := ctx.KVStore(iom.key)

	cidLinksAsBytes := store.Get([]byte(cid))

	var cidLinks CidLinks

	if cidLinksAsBytes == nil {
		cidLinks = make(CidLinks)
	} else {
		err := json.Unmarshal(cidLinksAsBytes, &cidLinks)
		if err != nil {
			return err
		}
	}

	users := cidLinks[linkedCid]
	if users == nil {
		users = make(map[string]struct{})
	}
	users[user.String()] = struct{}{}
	cidLinks[linkedCid] = users

	cidLinksAsBytes, err := json.Marshal(cidLinks)
	if err != nil {
		return err
	}

	store.Set([]byte(cid), cidLinksAsBytes)
	return nil
}

func (iom LinksStorage) GetLinks(ctx sdk.Context, cid string) (CidLinks, error) {
	store := ctx.KVStore(iom.key)

	cidLinksAsBytes := store.Get([]byte(cid))

	if cidLinksAsBytes == nil {
		return make(CidLinks), nil
	} else {
		cidLinks := CidLinks{}
		err := json.Unmarshal(cidLinksAsBytes, &cidLinks)
		return cidLinks, err
	}
}

type CidIndexStorage struct {
	cdc          *wire.Codec
	mainStoreKey sdk.StoreKey
	indexKey     sdk.StoreKey
}

// CIDs index is array of all added CIDs, sorted asc by first link time.
//   - for given link, CIDs added in order [CID1, CID2] (if they both new to chain)
// This method performs lookup of CIDs, returns index value, or create and put in index new value if not exists.
func (im CidIndexStorage) GetOrPutCidIndex(ctx sdk.Context, cid string) []byte {

	cidsIndex := ctx.KVStore(im.indexKey)

	cidAsBytes := []byte(cid)
	cidIndex := cidsIndex.Get(cidAsBytes)

	// new cid, get new index
	if cidIndex == nil {
		mainStore := ctx.KVStore(im.mainStoreKey)
		lastIndexAsBytes := mainStore.Get(cidsCount)

		// first link ever
		if lastIndexAsBytes == nil {
			zero := [8]byte{0, 0, 0, 0, 0, 0, 0, 0}
			lastIndexAsBytes = zero[:] // 0
		}

		lastIndex := binary.LittleEndian.Uint64(lastIndexAsBytes)
		lastIndex++
		binary.LittleEndian.PutUint64(lastIndexAsBytes, lastIndex)
		mainStore.Set(cidsCount, lastIndexAsBytes)
		return lastIndexAsBytes
	}

	return cidIndex
}
