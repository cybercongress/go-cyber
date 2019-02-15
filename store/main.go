package store

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math"
)

var lastCidNumberKey = []byte("cyberd_last_cid_number")
var linksCountKey = []byte("cyberd_links_count")
var genesisSupplyKey = []byte("cyberd_genesis_supply")
var lastBandwidthPrice = []byte("cyberd_last_bandwidth_price")
var spentBandwidth = []byte("cyberd_spent_bandwidth")
var latestBlockNumber = []byte("cyberd_latest_block_number")
var latestMerkleTree = []byte("cyberd_latest_merkle_tree")
var nextMerkleTree = []byte("cyberd_next_merkle_tree")
var rankCalculationFinished = []byte("cyberd_rank_calc_finished")

type MainKeeper struct {
	key *sdk.KVStoreKey
}

func NewMainKeeper(key *sdk.KVStoreKey) MainKeeper {
	return MainKeeper{key: key}
}

// returns overall added cids count
func (ms MainKeeper) GetCidsCount(ctx sdk.Context) uint64 {

	mainStore := ctx.KVStore(ms.key)
	lastIndexAsBytes := mainStore.Get(lastCidNumberKey)

	if lastIndexAsBytes == nil {
		return 0
	}

	return binary.LittleEndian.Uint64(lastIndexAsBytes) + 1
}

func (ms MainKeeper) GetLinksCount(ctx sdk.Context) uint64 {
	mainStore := ctx.KVStore(ms.key)
	linksCountAsBytes := mainStore.Get(linksCountKey)

	if linksCountAsBytes == nil {
		return 0
	}
	return binary.LittleEndian.Uint64(linksCountAsBytes)
}

func (ms MainKeeper) IncrementLinksCount(ctx sdk.Context) {
	mainStore := ctx.KVStore(ms.key)
	linksCount := ms.GetLinksCount(ctx) + 1
	linksCountAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(linksCountAsBytes, linksCount)
	mainStore.Set(linksCountKey, linksCountAsBytes)

}

func (ms MainKeeper) SetGenesisSupply(ctx sdk.Context, supply uint64) {
	mainStore := ctx.KVStore(ms.key)
	supplyAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(supplyAsBytes, supply)
	mainStore.Set(genesisSupplyKey, supplyAsBytes)
}

func (ms MainKeeper) GetGenesisSupply(ctx sdk.Context) uint64 {
	mainStore := ctx.KVStore(ms.key)
	supplyAsBytes := mainStore.Get(genesisSupplyKey)
	return binary.LittleEndian.Uint64(supplyAsBytes)
}

func (ms MainKeeper) SetLastCidIndex(ctx sdk.Context, cidsCount []byte) {

	mainStore := ctx.KVStore(ms.key)
	mainStore.Set(lastCidNumberKey, cidsCount)
}

func (ms MainKeeper) GetBandwidthPrice(ctx sdk.Context, basePrice float64) uint64 {
	store := ctx.KVStore(ms.key)
	priceAsBytes := store.Get(lastBandwidthPrice)
	if priceAsBytes == nil {
		priceAsBytes = make([]byte, 8)
		binary.LittleEndian.PutUint64(priceAsBytes, math.Float64bits(basePrice))
	}
	return binary.LittleEndian.Uint64(priceAsBytes)
}

func (ms MainKeeper) StoreBandwidthPrice(ctx sdk.Context, price uint64) {
	store := ctx.KVStore(ms.key)
	priceAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(priceAsBytes, price)
	store.Set(lastBandwidthPrice, priceAsBytes)
}

func (ms MainKeeper) GetSpentBandwidth(ctx sdk.Context) uint64 {
	store := ctx.KVStore(ms.key)
	bandwidthAsBytes := store.Get(spentBandwidth)
	if bandwidthAsBytes == nil {
		return 0
	}
	return binary.LittleEndian.Uint64(bandwidthAsBytes)
}

func (ms MainKeeper) StoreSpentBandwidth(ctx sdk.Context, bandwidth uint64) {
	store := ctx.KVStore(ms.key)
	bandwidthAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bandwidthAsBytes, bandwidth)
	store.Set(spentBandwidth, bandwidthAsBytes)
}

func (ms MainKeeper) GetLatestBlockNumber(ctx sdk.Context) uint64 {
	store := ctx.KVStore(ms.key)
	numberAsBytes := store.Get(latestBlockNumber)
	if numberAsBytes == nil {
		return 0
	}
	return binary.LittleEndian.Uint64(numberAsBytes)
}

func (ms MainKeeper) StoreLatestBlockNumber(ctx sdk.Context, number uint64) {
	store := ctx.KVStore(ms.key)
	numberAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(numberAsBytes, number)
	store.Set(latestBlockNumber, numberAsBytes)
}

func (ms MainKeeper) GetLatestMerkleTree(ctx sdk.Context) []byte {
	store := ctx.KVStore(ms.key)
	return store.Get(latestMerkleTree)
}

func (ms MainKeeper) StoreLatestMerkleTree(ctx sdk.Context, treeAsBytes []byte) {
	store := ctx.KVStore(ms.key)
	store.Set(latestMerkleTree, treeAsBytes)
}

func (ms MainKeeper) GetNextMerkleTree(ctx sdk.Context) []byte {
	store := ctx.KVStore(ms.key)
	return store.Get(nextMerkleTree)
}

func (ms MainKeeper) StoreNextMerkleTree(ctx sdk.Context, treeAsBytes []byte) {
	store := ctx.KVStore(ms.key)
	store.Set(nextMerkleTree, treeAsBytes)
}

func (ms MainKeeper) StoreRankCalculationFinished(ctx sdk.Context, finished bool) {
	store := ctx.KVStore(ms.key)
	var byteFlag byte
	if finished {
		byteFlag = 1
	}
	store.Set(rankCalculationFinished, []byte{byteFlag})
}

func (ms MainKeeper) GetRankCalculationFinished(ctx sdk.Context) bool {
	store := ctx.KVStore(ms.key)
	bytes := store.Get(rankCalculationFinished)
	if bytes == nil || bytes[0] == 1 {
		return true
	}
	return false
}
