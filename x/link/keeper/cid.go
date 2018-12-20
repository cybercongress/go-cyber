package keeper

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/store"
	. "github.com/cybercongress/cyberd/x/link/types"
)

type CidNumberKeeper interface {
	GetCidNumber(ctx sdk.Context, cid Cid) (CidNumber, bool)
	GetCid(ctx sdk.Context, num CidNumber) Cid
	GetOrPutCidNumber(ctx sdk.Context, cid Cid) CidNumber
	GetFullCidsNumbers(ctx sdk.Context) map[Cid]CidNumber
	GetCidsCount(ctx sdk.Context) uint64
}

type BaseCidNumberKeeper struct {
	ms         store.MainKeeper
	key        *sdk.KVStoreKey
	reverseKey *sdk.KVStoreKey
}

func NewBaseCidNumberKeeper(ms store.MainKeeper, key *sdk.KVStoreKey, reverseKey *sdk.KVStoreKey) CidNumberKeeper {
	return BaseCidNumberKeeper{
		ms:         ms,
		key:        key,
		reverseKey: reverseKey,
	}
}

// Return cid number and true, if cid exists
func (k BaseCidNumberKeeper) GetCidNumber(ctx sdk.Context, cid Cid) (CidNumber, bool) {
	cidsIndex := ctx.KVStore(k.key)
	cidAsBytes := []byte(cid)
	cidIndexAsBytes := cidsIndex.Get(cidAsBytes)
	if cidIndexAsBytes != nil {
		return CidNumber(binary.LittleEndian.Uint64(cidIndexAsBytes)), true
	}
	return 0, false
}

func (k BaseCidNumberKeeper) GetCid(ctx sdk.Context, num CidNumber) Cid {
	index := ctx.KVStore(k.reverseKey)
	cidNumberAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(cidNumberAsBytes, uint64(num))
	cidAsBytes := index.Get(cidNumberAsBytes)
	return Cid(cidAsBytes)
}

// CIDs index is array of all added CIDs, sorted asc by first link time.
//   - for given link, CIDs added in order [CID1, CID2] (if they both new to chain)
// This method performs lookup of CIDs, returns index value, or create and put in index new value if not exists.
func (k BaseCidNumberKeeper) GetOrPutCidNumber(ctx sdk.Context, cid Cid) CidNumber {

	cidsIndex := ctx.KVStore(k.key)
	cidsReverseIndex := ctx.KVStore(k.reverseKey)

	cidAsBytes := []byte(cid)
	cidIndexAsBytes := cidsIndex.Get(cidAsBytes)

	// new cid, get new index
	if cidIndexAsBytes == nil {

		lastIndex := k.ms.GetCidsCount(ctx)
		lastIndexAsBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(lastIndexAsBytes, lastIndex)

		cidsIndex.Set(cidAsBytes, lastIndexAsBytes)
		cidsReverseIndex.Set(lastIndexAsBytes, cidAsBytes)
		k.ms.SetLastCidIndex(ctx, lastIndexAsBytes)
		return CidNumber(lastIndex)
	}

	return CidNumber(binary.LittleEndian.Uint64(cidIndexAsBytes))
}

// returns all added cids
func (k BaseCidNumberKeeper) GetFullCidsNumbers(ctx sdk.Context) map[Cid]CidNumber {

	cidsIndex := ctx.KVStore(k.key)
	iterator := cidsIndex.Iterator(nil, nil)

	index := make(map[Cid]CidNumber)

	for iterator.Valid() {
		index[Cid(iterator.Key())] = CidNumber(binary.LittleEndian.Uint64(iterator.Value()))
		iterator.Next()
	}
	iterator.Close()
	return index
}

func (k BaseCidNumberKeeper) GetCidsCount(ctx sdk.Context) uint64 {
	return k.ms.GetCidsCount(ctx)
}
