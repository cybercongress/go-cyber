package storage

import (
	b "encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math"
)

type RankStorage struct {
	ms  MainStorage
	key *sdk.KVStoreKey
}

func NewRankStorage(ms MainStorage, key *sdk.KVStoreKey) RankStorage {
	return RankStorage{
		ms:  ms,
		key: key,
	}
}

// CIDs index is array of all added CIDs, sorted asc by first link time.
//   - for given link, CIDs added in order [CID1, CID2] (if they both new to chain)
// This method performs lookup of CIDs, returns index value, or create and put in index new value if not exists.
func (rs RankStorage) StoreFullRank(ctx sdk.Context, ranks []float64) {

	store := ctx.KVStore(rs.key)

	for i, rank := range ranks {
		indexAsBytes := make([]byte, 8)
		rankAsBytes := make([]byte, 8)
		b.LittleEndian.PutUint64(indexAsBytes, uint64(i))
		b.LittleEndian.PutUint64(rankAsBytes, math.Float64bits(rank))
		store.Set(indexAsBytes, rankAsBytes)
	}
}

// returns all added cids
func (rs RankStorage) GetFullRank(ctx sdk.Context) []float64 {

	cidsCount := rs.ms.GetCidsCount(ctx)
	ranks := make([]float64, cidsCount)

	store := ctx.KVStore(rs.key)
	iterator := store.Iterator(nil, nil)

	for iterator.Valid() {
		itVal := iterator.Value()
		itKey := iterator.Key()
		rankBits := b.LittleEndian.Uint64(itVal)
		rank := math.Float64frombits(rankBits)
		key := b.LittleEndian.Uint64(itKey)

		ranks[key] = rank
		iterator.Next()
	}
	iterator.Close()
	return ranks
}
