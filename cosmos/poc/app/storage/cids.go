package storage

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/cybercongress/cyberd/cosmos/poc/app/types"
)

type CidIndexStorage struct {
	ms             MainStorage
	indexKey       *sdk.KVStoreKey
	reversIndexKey *sdk.KVStoreKey
}

func NewCidIndexStorage(ms MainStorage, indexKey *sdk.KVStoreKey, reversIndexKey *sdk.KVStoreKey) CidIndexStorage {
	return CidIndexStorage{
		ms:             ms,
		indexKey:       indexKey,
		reversIndexKey: reversIndexKey,
	}
}

// Return cid number and true, if cid exists
func (cis CidIndexStorage) GetCidIndex(ctx sdk.Context, cid Cid) (CidNumber, bool) {
	cidsIndex := ctx.KVStore(cis.indexKey)
	cidAsBytes := []byte(cid)
	cidIndexAsBytes := cidsIndex.Get(cidAsBytes)
	if cidIndexAsBytes != nil {
		return CidNumber(binary.LittleEndian.Uint64(cidIndexAsBytes)), true
	}
	return 0, false
}

func (cis CidIndexStorage) GetCid(ctx sdk.Context, num CidNumber) Cid {
	index := ctx.KVStore(cis.reversIndexKey)
	cidNumberAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(cidNumberAsBytes, uint64(num))
	cidAsBytes := index.Get(cidNumberAsBytes)
	return Cid(cidAsBytes)
}

// CIDs index is array of all added CIDs, sorted asc by first link time.
//   - for given link, CIDs added in order [CID1, CID2] (if they both new to chain)
// This method performs lookup of CIDs, returns index value, or create and put in index new value if not exists.
func (cis CidIndexStorage) GetOrPutCidNumber(ctx sdk.Context, cid Cid) CidNumber {

	cidsIndex := ctx.KVStore(cis.indexKey)
	cidsReversIndex := ctx.KVStore(cis.reversIndexKey)

	cidAsBytes := []byte(cid)
	cidIndexAsBytes := cidsIndex.Get(cidAsBytes)

	// new cid, get new index
	if cidIndexAsBytes == nil {

		lastIndex := cis.ms.GetCidsCount(ctx)
		lastIndexAsBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(lastIndexAsBytes, lastIndex)

		cidsIndex.Set(cidAsBytes, lastIndexAsBytes)
		cidsReversIndex.Set(lastIndexAsBytes, cidAsBytes)
		cis.ms.SetLastCidIndex(ctx, lastIndexAsBytes)
		return CidNumber(lastIndex)
	}

	return CidNumber(binary.LittleEndian.Uint64(cidIndexAsBytes))
}

// returns all added cids
func (cis CidIndexStorage) GetFullCidsIndex(ctx sdk.Context) map[Cid]CidNumber {

	cidsIndex := ctx.KVStore(cis.indexKey)
	iterator := cidsIndex.Iterator(nil, nil)

	index := make(map[Cid]CidNumber)

	for iterator.Valid() {
		index[Cid(iterator.Key())] = CidNumber(binary.LittleEndian.Uint64(iterator.Value()))
		iterator.Next()
	}
	iterator.Close()
	return index
}
