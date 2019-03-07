package keeper

import (
	"encoding/binary"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cbdio "github.com/cybercongress/cyberd/io"
	"github.com/cybercongress/cyberd/store"
	. "github.com/cybercongress/cyberd/x/link/types"
	"io"
)

const (
	CidLengthBytesSize = uint64(1)
	CidNumberBytesSize = uint64(8)
	CidCountBytesSize  = uint64(8)
)

type CidNumberKeeper interface {
	GetCidNumber(ctx sdk.Context, cid Cid) (CidNumber, bool)
	GetCid(ctx sdk.Context, num CidNumber) Cid
	GetOrPutCidNumber(ctx sdk.Context, cid Cid) CidNumber
	GetFullCidsNumbers(ctx sdk.Context) map[Cid]CidNumber
	GetCidsCount(ctx sdk.Context) uint64
	PutCid(ctx sdk.Context, cid Cid, cidNumber CidNumber)
	Iterate(ctx sdk.Context, process func(Cid, CidNumber))
	WriteCids(ctx sdk.Context, writer io.Writer) (err error)
	LoadFromReader(ctx sdk.Context, reader io.Reader) (err error)
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

// WARNING: use only for state import. Don't forget to set right cid count after
func (k BaseCidNumberKeeper) PutCid(ctx sdk.Context, cid Cid, cidNumber CidNumber) {
	cidsIndex := ctx.KVStore(k.key)
	cidsReverseIndex := ctx.KVStore(k.reverseKey)

	cidNumberAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(cidNumberAsBytes, uint64(cidNumber))

	cidsIndex.Set([]byte(cid), cidNumberAsBytes)
	cidsReverseIndex.Set(cidNumberAsBytes, []byte(cid))
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

		lastIndex := k.GetCidsCount(ctx)
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
	index := make(map[Cid]CidNumber)
	k.Iterate(ctx, func(cid Cid, number CidNumber) {
		index[cid] = number
	})
	return index
}

func (k BaseCidNumberKeeper) Iterate(ctx sdk.Context, process func(Cid, CidNumber)) {
	iterator := ctx.KVStore(k.key).Iterator(nil, nil)
	defer iterator.Close()

	for iterator.Valid() {
		process(Cid(iterator.Key()), CidNumber(binary.LittleEndian.Uint64(iterator.Value())))
		iterator.Next()
	}
}

func (k BaseCidNumberKeeper) GetCidsCount(ctx sdk.Context) uint64 {
	return k.ms.GetCidsCount(ctx)
}

// write CIDs to writer in binary format: <n><cid1_size><cid1><cid1_number><cid2_size><cid2><cid2_number>....<cidn_size><cidn><cidn_number>
func (k BaseCidNumberKeeper) WriteCids(ctx sdk.Context, writer io.Writer) (err error) {
	uintAsBytes := make([]byte, 8) //common bytes array to convert uints

	cidsCount := k.GetCidsCount(ctx)
	binary.LittleEndian.PutUint64(uintAsBytes, cidsCount)
	_, err = writer.Write(uintAsBytes)
	if err != nil {
		return
	}

	k.Iterate(ctx, func(cid Cid, number CidNumber) {
		cidLength := len(cid)
		if cidLength > 255 {
			err = errors.New("cid length cannot be over 255")
			return
		}

		_, err = writer.Write([]byte{byte(cidLength)})
		if err != nil {
			return
		}
		_, err = writer.Write([]byte(cid))
		if err != nil {
			return
		}
		binary.LittleEndian.PutUint64(uintAsBytes, uint64(number))
		_, err = writer.Write(uintAsBytes)
		if err != nil {
			return
		}
	})
	return
}

func (k BaseCidNumberKeeper) LoadFromReader(ctx sdk.Context, reader io.Reader) (err error) {
	cidCountBytes, err := cbdio.ReadExactlyNBytes(reader, CidCountBytesSize)
	if err != nil {
		return
	}
	cidCount := binary.LittleEndian.Uint64(cidCountBytes)

	// Read all CIDs with their numbers
	for i := uint64(0); i < cidCount; i++ {
		cidLengthBytes, err := cbdio.ReadExactlyNBytes(reader, CidLengthBytesSize)
		if err != nil {
			return err
		}

		cidBytes, err := cbdio.ReadExactlyNBytes(reader, uint64(cidLengthBytes[0]))
		if err != nil {
			return err
		}
		cid := Cid(cidBytes)
		cidNumberBytes, err := cbdio.ReadExactlyNBytes(reader, CidNumberBytesSize)
		if err != nil {
			return err
		}
		cidNumber := CidNumber(binary.LittleEndian.Uint64(cidNumberBytes))
		k.PutCid(ctx, cid, cidNumber)
	}

	lastCidIndex := cidCount - 1
	lastCidIndexBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(lastCidIndexBytes, lastCidIndex)
	k.ms.SetLastCidIndex(ctx, lastCidIndexBytes)

	return
}
