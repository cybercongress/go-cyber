package storage

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var lastCidNumber = []byte("cyberd_cids_count")
var lastAppHash = []byte("cyberd_app_hash")

type MainStorage struct {
	key *sdk.KVStoreKey
}

func NewMainStorage(key *sdk.KVStoreKey) MainStorage {
	return MainStorage{key: key}
}

// returns overall added cids count
func (ms MainStorage) GetCidsCount(ctx sdk.Context) uint64 {

	mainStore := ctx.KVStore(ms.key)
	lastIndexAsBytes := mainStore.Get(lastCidNumber)

	if lastIndexAsBytes == nil {
		return 0
	}

	return binary.LittleEndian.Uint64(lastIndexAsBytes) + 1
}

func (ms MainStorage) SetLastCidIndex(ctx sdk.Context, cidsCount []byte) {

	mainStore := ctx.KVStore(ms.key)
	mainStore.Set(lastCidNumber, cidsCount)
}

func (ms MainStorage) GetAppHash(ctx sdk.Context) []byte {

	store := ctx.KVStore(ms.key)
	return store.Get(lastAppHash)
}

func (ms MainStorage) StoreAppHash(ctx sdk.Context, hash []byte) {

	store := ctx.KVStore(ms.key)
	store.Set(lastAppHash, hash)
}
