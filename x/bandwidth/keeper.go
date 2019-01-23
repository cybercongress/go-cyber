package bandwidth

import (
	"encoding/binary"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/cybercongress/cyberd/x/bandwidth/types"
)

var _ Keeper = BaseAccBandwidthKeeper{}

type BaseAccBandwidthKeeper struct {
	key *sdk.KVStoreKey
}

func NewAccBandwidthKeeper(key *sdk.KVStoreKey) BaseAccBandwidthKeeper {
	return BaseAccBandwidthKeeper{key: key}
}

func (bk BaseAccBandwidthKeeper) SetAccBandwidth(ctx sdk.Context, bandwidth AcсBandwidth) {
	bwBytes, _ := json.Marshal(bandwidth)
	ctx.KVStore(bk.key).Set(bandwidth.Address, bwBytes)
}

func (bk BaseAccBandwidthKeeper) GetAccBandwidth(ctx sdk.Context, addr sdk.AccAddress) (bw AcсBandwidth) {
	bwBytes := ctx.KVStore(bk.key).Get(addr)
	if bwBytes == nil {
		return AcсBandwidth{
			Address:          addr,
			RemainedValue:    0,
			LastUpdatedBlock: ctx.BlockHeight(),
			MaxValue:         0,
		}
	}
	err := json.Unmarshal(bwBytes, &bw)
	if err != nil {
		// should not happen
		panic("bandwidth index is broken")
	}
	return
}

type BaseBlockSpentBandwidthKeeper struct {
	key *sdk.KVStoreKey
}

func NewBlockSpentBandwidthKeeper(key *sdk.KVStoreKey) BaseBlockSpentBandwidthKeeper {
	return BaseBlockSpentBandwidthKeeper{key: key}
}

func (bk BaseBlockSpentBandwidthKeeper) SetBlockSpentBandwidth(ctx sdk.Context, blockNumber uint64, value uint64) {
	keyAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(keyAsBytes, blockNumber)
	valueAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(valueAsBytes, value)
	ctx.KVStore(bk.key).Set(keyAsBytes, valueAsBytes)
}

func (bk BaseBlockSpentBandwidthKeeper) GetValuesForPeriod(ctx sdk.Context, period int64) map[uint64]uint64 {

	startKey := make([]byte, 8)
	binary.LittleEndian.PutUint64(startKey, uint64(ctx.BlockHeight()-period+1))

	endKey := make([]byte, 8)
	binary.LittleEndian.PutUint64(endKey, uint64(ctx.BlockHeight()))

	iterator := ctx.KVStore(bk.key).Iterator(startKey, sdk.PrefixEndBytes(endKey))
	defer iterator.Close()

	result := make(map[uint64]uint64)
	for iterator.Valid() {
		blockNumber := binary.LittleEndian.Uint64(iterator.Key())
		value := binary.LittleEndian.Uint64(iterator.Value())
		result[blockNumber] = value
		iterator.Next()
	}

	return result
}
