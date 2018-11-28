package storage

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var lastCidNumberKey = []byte("cyberd_last_cid_number")
var linksCountKey = []byte("cyberd_links_count")
var lastAppHashKey = []byte("cyberd_app_hash")

type MainStorage struct {
	key *sdk.KVStoreKey
}

func NewMainStorage(key *sdk.KVStoreKey) MainStorage {
	return MainStorage{key: key}
}

// returns overall added cids count
func (ms MainStorage) GetCidsCount(ctx sdk.Context) uint64 {

	mainStore := ctx.KVStore(ms.key)
	lastIndexAsBytes := mainStore.Get(lastCidNumberKey)

	if lastIndexAsBytes == nil {
		return 0
	}

	return binary.LittleEndian.Uint64(lastIndexAsBytes) + 1
}

func (ms MainStorage) GetLinksCount(ctx sdk.Context) uint64 {
	mainStore := ctx.KVStore(ms.key)
	linksCountAsBytes := mainStore.Get(linksCountKey)

	if linksCountAsBytes == nil {
		return 0
	}
	return binary.LittleEndian.Uint64(linksCountAsBytes)
}

func (ms MainStorage) IncrementLinksCount(ctx sdk.Context) {
	mainStore := ctx.KVStore(ms.key)
	linksCount := ms.GetLinksCount(ctx) + 1
	linksCountAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(linksCountAsBytes, linksCount)
	mainStore.Set(linksCountKey, linksCountAsBytes)

}

func (ms MainStorage) SetLastCidIndex(ctx sdk.Context, cidsCount []byte) {

	mainStore := ctx.KVStore(ms.key)
	mainStore.Set(lastCidNumberKey, cidsCount)
}

func (ms MainStorage) GetAppHash(ctx sdk.Context) []byte {
	store := ctx.KVStore(ms.key)
	return store.Get(lastAppHashKey)
}

func (ms MainStorage) StoreAppHash(ctx sdk.Context, hash []byte) {

	store := ctx.KVStore(ms.key)
	store.Set(lastAppHashKey, hash)
}
