package keeper

import (
	"encoding/binary"
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/cybercongress/cyberd/x/bandwidth/internal/types"
)

type BaseAccountBandwidthKeeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	paramSpace params.Subspace
}

func NewAccountBandwidthKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace) BaseAccountBandwidthKeeper {
	return BaseAccountBandwidthKeeper{
		cdc:        cdc,
		storeKey:   key,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

func (k BaseAccountBandwidthKeeper) SetAccountBandwidth(ctx sdk.Context, bandwidth types.AcсountBandwidth) {
	bwBytes, _ := json.Marshal(bandwidth)
	ctx.KVStore(k.storeKey).Set(bandwidth.Address, bwBytes)
}

func (bk BaseAccountBandwidthKeeper) GetAccountBandwidth(ctx sdk.Context, addr sdk.AccAddress) (bw types.AcсountBandwidth) {
	bwBytes := ctx.KVStore(bk.storeKey).Get(addr)
	if bwBytes == nil {
		return types.AcсountBandwidth{
			Address:          addr,
			RemainedValue:    0,
			LastUpdatedBlock: ctx.BlockHeight(),
			MaxValue:         0,
		}
	}
	err := json.Unmarshal(bwBytes, &bw)
	if err != nil {
		panic("bandwidth index is broken")
	}
	return
}

func (bk BaseAccountBandwidthKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	bk.paramSpace.GetParamSet(ctx, &params)
	return params
}

func (bk BaseAccountBandwidthKeeper) SetParams(ctx sdk.Context, params types.Params) {
	bk.paramSpace.SetParamSet(ctx, &params)
}

type BaseBlockSpentBandwidthKeeper struct {
	storeKey sdk.StoreKey
}

func NewBlockSpentBandwidthKeeper(key sdk.StoreKey) BaseBlockSpentBandwidthKeeper {
	return BaseBlockSpentBandwidthKeeper{
		storeKey: key,
	}
}

func (bk BaseBlockSpentBandwidthKeeper) SetBlockSpentBandwidth(ctx sdk.Context, blockNumber uint64, value uint64) {
	keyAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(keyAsBytes, blockNumber)
	valueAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(valueAsBytes, value)

	store := ctx.KVStore(bk.storeKey)
	store.Set(keyAsBytes, valueAsBytes)
}

func (bk BaseBlockSpentBandwidthKeeper) GetValuesForPeriod(ctx sdk.Context, period int64) map[uint64]uint64 {

	store := ctx.KVStore(bk.storeKey)

	windowStart := ctx.BlockHeight() - period + 1
	if windowStart <= 0 {
		windowStart = 1
	}

	key := make([]byte, 8)
	result := make(map[uint64]uint64)
	for blockNumber := windowStart; blockNumber <= ctx.BlockHeight(); blockNumber++ {
		binary.LittleEndian.PutUint64(key, uint64(blockNumber))
		value := binary.LittleEndian.Uint64(store.Get(key))
		result[uint64(blockNumber)] = value
	}

	return result
}
