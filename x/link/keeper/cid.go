package keeper

import (
	"encoding/binary"

	"github.com/cybercongress/go-cyber/util"
	//"github.com/cybercongress/go-cyber/store"
	//"github.com/cybercongress/go-cyber/util"
	"github.com/cybercongress/go-cyber/x/link/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	//"encoding/binary"
	//"errors"
	"io"
)

const (
	CidLengthBytesSize = uint64(1)
	CidNumberBytesSize = uint64(8)
	CidCountBytesSize  = uint64(8)
)


// Return cid number and true, if cid exists
func (gk GraphKeeper) GetCidNumber(ctx sdk.Context, cid types.Cid) (types.CidNumber, bool) {
	store := ctx.KVStore(gk.key)
	cidIndexAsBytes := store.Get(types.CidStoreKey(cid))

	if cidIndexAsBytes != nil {
		return types.CidNumber(types.BigEndianToUint64(cidIndexAsBytes)), true
	}

	return 0, false
}

func (gk GraphKeeper) GetCid(ctx sdk.Context, num types.CidNumber) types.Cid {
	store := ctx.KVStore(gk.key)
	cidAsBytes := store.Get(types.CidReverseStoreKey(num))

	return types.Cid(cidAsBytes)
}

// CIDs index is array of all added CIDs, sorted asc by first link time.
// for given link, CIDs added in order [CID1, CID2] (if they both new to chain)
// This method performs lookup of CIDs, returns index value, or create and put in index new value if not exists.
func (gk GraphKeeper) GetOrPutCidNumber(ctx sdk.Context, cid types.Cid) types.CidNumber {
	store := ctx.KVStore(gk.key)
	cidIndexAsBytes := store.Get(types.CidStoreKey(cid))

	if cidIndexAsBytes == nil {
		lastIndex := gk.GetCidsCount(ctx)
		store.Set(types.CidStoreKey(cid), sdk.Uint64ToBigEndian(lastIndex))
		store.Set(types.CidReverseStoreKey(types.CidNumber(lastIndex)), []byte(cid))
		store.Set(types.LastCidNumber, sdk.Uint64ToBigEndian(lastIndex))

		return types.CidNumber(lastIndex)
	}

	return types.CidNumber(types.BigEndianToUint64(cidIndexAsBytes))
}

func (gk GraphKeeper) GetCidsCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(gk.key)
	lastIndexAsBytes := store.Get(types.LastCidNumber)

	if lastIndexAsBytes == nil {
		return 0
	}

	return types.BigEndianToUint64(lastIndexAsBytes) + 1
}

func (gk GraphKeeper) SetLastCidIndex(ctx sdk.Context, lastCidIndex uint64) {
	store := ctx.KVStore(gk.key)
	store.Set(types.LastCidNumber, sdk.Uint64ToBigEndian(lastCidIndex))
}

// NOTICE: use only for state import. Don't forget to set right cid count after
func (gk GraphKeeper) PutCid(ctx sdk.Context, cid types.Cid, cidNumber types.CidNumber) {
	store := ctx.KVStore(gk.key)

	//fmt.Println("[GraphKeeper] PutCid: ", cid, " ",cidNumber)
	store.Set(types.CidStoreKey(cid), sdk.Uint64ToBigEndian(uint64(cidNumber)))
	store.Set(types.CidReverseStoreKey(cidNumber), []byte(cid))
}

func (gk GraphKeeper) IterateCids(ctx sdk.Context, process func(types.Cid, types.CidNumber)) {
	store := ctx.KVStore(gk.key)
	iterator := sdk.KVStorePrefixIterator(store, types.CidStoreKeyPrefix)
	defer iterator.Close()

	for iterator.Valid() {
		// [1:0] because we have prefix []byte{0x01} for cids
		process(types.Cid(iterator.Key()[1:]), types.CidNumber(types.BigEndianToUint64(iterator.Value())))
		iterator.Next()
	}
}

// write CIDs to writer in binary format: <n><cid1_size><cid1><cid1_number><cid2_size><cid2><cid2_number>....<cidn_size><cidn><cidn_number>
func (gk GraphKeeper) WriteCids(ctx sdk.Context, writer io.Writer) (err error) {
	uintAsBytes := make([]byte, 8) //common bytes array to convert uints

	cidsCount := gk.GetCidsCount(ctx)
	binary.LittleEndian.PutUint64(uintAsBytes, cidsCount)
	_, err = writer.Write(uintAsBytes)
	if err != nil {
		return
	}

	gk.IterateCids(ctx, func(cid types.Cid, number types.CidNumber) {
		cidLength := len(cid)
		if cidLength > 255 {
			//err = errors.New("cid length cannot be over 255")
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
		//fmt.Println("CIDs: ", cidLength, " ", cid, " ",number)
	})
	return
}

func (gk GraphKeeper) LoadFromReader(ctx sdk.Context, reader io.Reader) (err error) {
	cidCountBytes, err := util.ReadExactlyNBytes(reader, CidCountBytesSize)
	if err != nil {
		return
	}
	cidCount := binary.LittleEndian.Uint64(cidCountBytes)

	// Read all CIDs with their numbers
	for i := uint64(0); i < cidCount; i++ {
		cidLengthBytes, err := util.ReadExactlyNBytes(reader, CidLengthBytesSize)
		if err != nil {
			return err
		}

		cidBytes, err := util.ReadExactlyNBytes(reader, uint64(cidLengthBytes[0]))
		if err != nil {
			return err
		}
		cid := types.Cid(cidBytes)
		cidNumberBytes, err := util.ReadExactlyNBytes(reader, CidNumberBytesSize)
		if err != nil {
			return err
		}
		cidNumber := types.CidNumber(binary.LittleEndian.Uint64(cidNumberBytes))
		//fmt.Println("CIDs: ", uint64(cidLengthBytes[0]), " ", cid, " ",cidNumber)
		gk.PutCid(ctx, cid, cidNumber)
	}

	lastCidIndex := cidCount - 1
	gk.SetLastCidIndex(ctx, lastCidIndex)
	//fmt.Println("CIDs count: ", gk.GetCidsCount(ctx))

	return
}
