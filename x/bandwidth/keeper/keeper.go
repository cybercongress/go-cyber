package keeper

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cybercongress/go-cyber/x/bandwidth/types"
)

//var _ exported.BandwidthMeter = &BaseBandwidthMeter{}
type MsgBandwidthCost func(ctx sdk.Context, params types.Params, msg sdk.Msg) uint64

type BandwidthMeter struct {
	stakeProvider           types.AccountStakeProvider
	msgCost MsgBandwidthCost

	storeKey         			sdk.StoreKey
	paramSpace					paramtypes.Subspace

	currentBlockSpentBandwidth  uint64
	currentCreditPrice          float64
	bandwidthSpentByBlock       map[uint64]uint64
	totalSpentForSlidingWindow  uint64
}

func NewBandwidthMeter(
	key sdk.StoreKey,
	asp types.AccountStakeProvider,
	mc MsgBandwidthCost,
	paramSpace paramtypes.Subspace,
) *BandwidthMeter {

	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &BandwidthMeter{
		storeKey: 				 key,
		stakeProvider:           asp,
		msgCost:                 mc,
		paramSpace:     		 paramSpace,
		bandwidthSpentByBlock:   make(map[uint64]uint64),
	}
}

func (bm BandwidthMeter) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (bm BandwidthMeter) GetParams(ctx sdk.Context) (params types.Params) {
	bm.paramSpace.GetParamSet(ctx, &params)
	return params
}

func (bm BandwidthMeter) SetParams(ctx sdk.Context, params types.Params) {
	bm.paramSpace.SetParamSet(ctx, &params)
}

//______________________________________________________________________

func (m *BandwidthMeter) Load(ctx sdk.Context) {
	params := m.GetParams(ctx)
	m.totalSpentForSlidingWindow = 0
	m.bandwidthSpentByBlock = m.GetValuesForPeriod(ctx, params.RecoveryPeriod)
	for _, spentBandwidth := range m.bandwidthSpentByBlock {
		m.totalSpentForSlidingWindow += spentBandwidth
	}
	floatBaseCreditPrice, err := strconv.ParseFloat(params.BaseCreditPrice.String(), 64)
	if err != nil {
		panic(err)
	}
	m.currentCreditPrice = 0.01 * math.Float64frombits(m.GetBandwidthPrice(ctx, floatBaseCreditPrice))
	m.currentBlockSpentBandwidth = 0
}

//____________________________________________________________________________

// TODO use from sdk after version bump
func BigEndianToUint64(bz []byte) uint64 {
	if len(bz) == 0 {
		return 0
	}

	return binary.BigEndian.Uint64(bz)
}

func (bm BandwidthMeter) GetBandwidthPrice(ctx sdk.Context, basePrice float64) uint64 {
	store := ctx.KVStore(bm.storeKey)
	priceAsBytes := store.Get(types.LastBandwidthPrice)
	if priceAsBytes == nil {
		priceAsBytes = sdk.Uint64ToBigEndian(math.Float64bits(basePrice))
	}
	return BigEndianToUint64(priceAsBytes)
}

func (bm BandwidthMeter) StoreBandwidthPrice(ctx sdk.Context, price uint64) {
	store := ctx.KVStore(bm.storeKey)
	store.Set(types.LastBandwidthPrice, sdk.Uint64ToBigEndian(price))
}

//____________________________________________________________________________

func (m *BandwidthMeter) AddToBlockBandwidth(value uint64) {
	m.currentBlockSpentBandwidth += value
}

// Here we move bandwidth window:
// Remove first block of window and add new block to window end
func (m *BandwidthMeter) CommitBlockBandwidth(ctx sdk.Context) {
	m.totalSpentForSlidingWindow += m.currentBlockSpentBandwidth

	newWindowEnd := ctx.BlockHeight()
	params := m.GetParams(ctx)
	windowStart := newWindowEnd - int64(params.RecoveryPeriod)
	if windowStart < 0 { // check needed cause it will be casted to uint and can cause overflow
		windowStart = 0
	}

	windowStartValue, exists := m.bandwidthSpentByBlock[uint64(windowStart)]
	if exists {
		m.totalSpentForSlidingWindow -= windowStartValue
		delete(m.bandwidthSpentByBlock, uint64(windowStart))
	}
	m.SetBlockBandwidth(ctx, uint64(ctx.BlockHeight()), m.currentBlockSpentBandwidth)
	m.bandwidthSpentByBlock[uint64(newWindowEnd)] = m.currentBlockSpentBandwidth
	m.currentBlockSpentBandwidth = 0
}

func (m *BandwidthMeter) GetCurrentBlockSpentBandwidth(ctx sdk.Context) uint64 {
	return m.currentBlockSpentBandwidth
}

func (m *BandwidthMeter) GetCurrentNetworkLoad(ctx sdk.Context) float64 {
	params := m.GetParams(ctx)
	return float64(m.totalSpentForSlidingWindow) / float64(params.DesirableBandwidth)
}

func (m *BandwidthMeter) GetMaxBlockBandwidth(ctx sdk.Context) uint64 {
	params := m.GetParams(ctx)
	maxBlockBandwidth := params.MaxBlockBandwidth
	return maxBlockBandwidth
}

//____________________________________________________________________________

func (m *BandwidthMeter) GetCurrentCreditPrice() float64 {
	return m.currentCreditPrice
}

func (m *BandwidthMeter) AdjustPrice(ctx sdk.Context) {
	params := m.GetParams(ctx)
	floatBaseCreditPrice, err := strconv.ParseFloat(params.BaseCreditPrice.String(), 64)
	if err != nil {
		panic(err)
	}
	newPrice := float64(m.totalSpentForSlidingWindow) / float64(params.DesirableBandwidth)

	if newPrice < 0.01 * floatBaseCreditPrice {
		newPrice = 0.01 * floatBaseCreditPrice
	}

	m.currentCreditPrice = newPrice
	m.StoreBandwidthPrice(ctx, math.Float64bits(newPrice))
}

//____________________________________________________________________________

func (m *BandwidthMeter) GetTxCost(ctx sdk.Context, tx sdk.Tx) uint64 {
	params := m.GetParams(ctx)
	bandwidthForTx := params.TxCost
	for _, msg := range tx.GetMsgs() {
		bandwidthForTx = bandwidthForTx + m.msgCost(ctx, params, msg)
	}
	return bandwidthForTx
}

func (m *BandwidthMeter) GetPricedTxCost(ctx sdk.Context, tx sdk.Tx) uint64 {
	return uint64(float64(m.GetTxCost(ctx, tx)) * m.currentCreditPrice)
}

//____________________________________________________________________________

// Performs bw consumption for given acc
// To get right number, should be called after tx delivery with bw state obtained prior delivery
// Pseudo code:
// bw := getCurrentBw(addr)
// bwCost := deliverTx(tx)
// consumeBw(bw, bwCost)
func (m *BandwidthMeter) ConsumeAccountBandwidth(ctx sdk.Context, bw types.AcсountBandwidth, amt uint64) {
	bw.Consume(amt)
	m.SetAccountBandwidth(ctx, bw)
	bw = m.GetCurrentAccountBandwidth(ctx, bw.Address)
	m.SetAccountBandwidth(ctx, bw)
}

func (m *BandwidthMeter) GetCurrentAccountBandwidth(ctx sdk.Context, address sdk.AccAddress) types.AcсountBandwidth {
	accBw := m.GetAccountBandwidth(ctx, address)
	accMaxBw := m.GetAccountMaxBandwidth(ctx, address)
	params := m.GetParams(ctx)
	accBw.UpdateMax(accMaxBw, uint64(ctx.BlockHeight()), params.RecoveryPeriod)
	return accBw
}

func (m *BandwidthMeter) GetAccountMaxBandwidth(ctx sdk.Context, addr sdk.AccAddress) uint64 {
	accStakePercentage := m.stakeProvider.GetAccountStakePercentage(ctx, addr)
	params := m.GetParams(ctx)
	return uint64(accStakePercentage * float64(params.DesirableBandwidth))
}

func (m *BandwidthMeter) UpdateAccountMaxBandwidth(ctx sdk.Context, address sdk.AccAddress) {
	bw := m.GetCurrentAccountBandwidth(ctx, address)
	m.SetAccountBandwidth(ctx, bw)
}

//____________________________________________________________________________
