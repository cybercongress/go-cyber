package storage

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var lastCidNumberKey = []byte("cyberd_last_cid_number")
var linksCountKey = []byte("cyberd_links_count")
var lastAppHashKey = []byte("cyberd_app_hash")
var genesisSupplyKey = []byte("cyberd_genesis_supply")
var lastBandwidthPrice = []byte("cyberd_last_bandwidth_price")
var spentBandwidth = []byte("cyberd_spent_bandwidth")
var latestBlockNumber = []byte("cyberd_latest_block_number")

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

func (ms MainStorage) SetGenesisSupply(ctx sdk.Context, supply uint64) {
	mainStore := ctx.KVStore(ms.key)
	supplyAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(supplyAsBytes, supply)
	mainStore.Set(genesisSupplyKey, supplyAsBytes)
}

func (ms MainStorage) GetGenesisSupply(ctx sdk.Context) uint64 {
	mainStore := ctx.KVStore(ms.key)
	supplyAsBytes := mainStore.Get(genesisSupplyKey)
	return binary.LittleEndian.Uint64(supplyAsBytes)
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

func (ms MainStorage) GetBandwidthPrice(ctx sdk.Context, basePrice uint64) uint64 {
	store := ctx.KVStore(ms.key)
	priceAsBytes := store.Get(lastBandwidthPrice)
	if priceAsBytes == nil {
		return basePrice
	}
	return binary.LittleEndian.Uint64(priceAsBytes)
}

func (ms MainStorage) StoreBandwidthPrice(ctx sdk.Context, price uint64) {
	store := ctx.KVStore(ms.key)
	priceAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(priceAsBytes, price)
	store.Set(lastBandwidthPrice, priceAsBytes)
}


func (ms MainStorage) GetSpentBandwidth(ctx sdk.Context) uint64 {
	store := ctx.KVStore(ms.key)
	bandwidthAsBytes := store.Get(spentBandwidth)
	if bandwidthAsBytes == nil {
		return 0
	}
	return binary.LittleEndian.Uint64(bandwidthAsBytes)
}

func (ms MainStorage) StoreSpentBandwidth(ctx sdk.Context, bandwidth uint64) {
	store := ctx.KVStore(ms.key)
	bandwidthAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bandwidthAsBytes, bandwidth)
	store.Set(spentBandwidth, bandwidthAsBytes)
}

func (ms MainStorage) GetLatestBlockNumber(ctx sdk.Context) uint64 {
	store := ctx.KVStore(ms.key)
	numberAsBytes := store.Get(latestBlockNumber)
	if numberAsBytes == nil {
		return 0
	}
	return binary.LittleEndian.Uint64(numberAsBytes)
}

func (ms MainStorage) StoreLatestBlockNumber(ctx sdk.Context, number uint64) {
	store := ctx.KVStore(ms.key)
	numberAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(numberAsBytes, number)
	store.Set(latestBlockNumber, numberAsBytes)
}
