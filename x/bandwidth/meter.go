package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cybercongress/cyberd/store"
	"github.com/cybercongress/cyberd/x/bandwidth/types"
	"math"
	"strconv"
)

var _ types.BandwidthMeter = &BaseBandwidthMeter{}

type BaseBandwidthMeter struct {
	// data providers
	accKeeper            auth.AccountKeeper
	stakeProvider        types.AccStakeProvider
	bwKeeper             types.Keeper
	mainKeeper           store.MainKeeper
	blockBandwidthKeeper types.BlockSpentBandwidthKeeper
	paramsKeeper         params.Keeper
	// bw configuration
	msgCost types.MsgBandwidthCost

	// price adjustment fields
	curBlockSpentBandwidth     uint64 //resets every block
	currentCreditPrice         float64
	bandwidthSpent             map[uint64]uint64 // bandwidth spent by blocks
	totalSpentForSlidingWindow uint64
}

func NewBaseMeter(
	mk store.MainKeeper, pk params.Keeper, ak auth.AccountKeeper, bwk types.Keeper,
	bbwk types.BlockSpentBandwidthKeeper, sp types.AccStakeProvider, msgCost types.MsgBandwidthCost,
) *BaseBandwidthMeter {

	return &BaseBandwidthMeter{
		mainKeeper:           mk,
		blockBandwidthKeeper: bbwk,
		accKeeper:            ak,
		bwKeeper:             bwk,
		paramsKeeper:         pk,
		stakeProvider:        sp,
		msgCost:              msgCost,
		bandwidthSpent:       make(map[uint64]uint64),
	}
}

func (m *BaseBandwidthMeter) Load(ctx sdk.Context) {
	paramset := m.GetParamSet(ctx)
	m.totalSpentForSlidingWindow = 0
	m.bandwidthSpent = m.blockBandwidthKeeper.GetValuesForPeriod(ctx, paramset.SlidingWindowSize)
	for _, spentBandwidth := range m.bandwidthSpent {
		m.totalSpentForSlidingWindow += spentBandwidth
	}
	floatBaseCreditPrice, err := strconv.ParseFloat(paramset.BaseCreditPrice, 64)
	if err != nil {
		panic(err)
	}
	m.currentCreditPrice = math.Float64frombits(m.mainKeeper.GetBandwidthPrice(ctx, floatBaseCreditPrice))
	m.curBlockSpentBandwidth = 0
}

func (m *BaseBandwidthMeter) AddToBlockBandwidth(value int64) {
	m.curBlockSpentBandwidth += uint64(value)
}

// Here we move bandwidth window:
// Remove first block of window and add new block to window end
func (m *BaseBandwidthMeter) CommitBlockBandwidth(ctx sdk.Context) {
	m.totalSpentForSlidingWindow += m.curBlockSpentBandwidth

	newWindowEnd := ctx.BlockHeight()
	paramset := m.GetParamSet(ctx)
	windowStart := newWindowEnd - paramset.SlidingWindowSize
	if windowStart < 0 { // check needed cause it will be casted to uint and can cause overflow
		windowStart = 0
	}
	windowStartValue, exists := m.bandwidthSpent[uint64(windowStart)]
	if exists {
		m.totalSpentForSlidingWindow -= windowStartValue
		delete(m.bandwidthSpent, uint64(windowStart))
	}
	m.blockBandwidthKeeper.SetBlockSpentBandwidth(ctx, uint64(ctx.BlockHeight()), m.curBlockSpentBandwidth)
	m.bandwidthSpent[uint64(newWindowEnd)] = m.curBlockSpentBandwidth
	m.curBlockSpentBandwidth = 0
}

func (m *BaseBandwidthMeter) AdjustPrice(ctx sdk.Context) {
	paramset := m.GetParamSet(ctx)
	floatBaseCreditPrice, err := strconv.ParseFloat(paramset.BaseCreditPrice, 64)
	if err != nil {
		panic(err)
	}

	floatShouldBeSpentPerSlidingWindow, err := strconv.ParseFloat(paramset.ShouldBeSpentPerSlidingWindow, 64)
	if err != nil {
		panic(err)
	}

	newPrice := float64(m.totalSpentForSlidingWindow) / floatShouldBeSpentPerSlidingWindow

	if newPrice < 0.01 * floatBaseCreditPrice {
		newPrice = 0.01 * floatBaseCreditPrice
	}

	m.currentCreditPrice = newPrice
	m.mainKeeper.StoreBandwidthPrice(ctx, math.Float64bits(newPrice))
}

func (m *BaseBandwidthMeter) GetTxCost(ctx sdk.Context, tx sdk.Tx) int64 {
	paramset := m.GetParamSet(ctx)
	bandwidthForTx := paramset.TxCost
	for _, msg := range tx.GetMsgs() {
		bandwidthForTx = bandwidthForTx + m.msgCost(ctx, m.paramsKeeper, msg)
	}
	return bandwidthForTx
}

func (m *BaseBandwidthMeter) GetPricedTxCost(ctx sdk.Context, tx sdk.Tx) int64 {
	return int64(float64(m.GetTxCost(ctx, tx)) * m.currentCreditPrice)
}

func (m *BaseBandwidthMeter) GetAccMaxBandwidth(ctx sdk.Context, addr sdk.AccAddress) int64 {
	accStakePercentage := m.stakeProvider.GetAccStakePercentage(ctx, addr)
	paramset := m.GetParamSet(ctx)
	return int64(accStakePercentage * float64(paramset.DesirableNetworkBandwidthForRecoveryPeriod))
}

func (m *BaseBandwidthMeter) GetCurrentAccBandwidth(ctx sdk.Context, address sdk.AccAddress) types.AcсBandwidth {
	accBw := m.bwKeeper.GetAccBandwidth(ctx, address)
	accMaxBw := m.GetAccMaxBandwidth(ctx, address)
	paramset := m.GetParamSet(ctx)
	accBw.UpdateMax(accMaxBw, ctx.BlockHeight(), paramset.RecoveryPeriod)
	return accBw
}

func (m *BaseBandwidthMeter) UpdateAccMaxBandwidth(ctx sdk.Context, address sdk.AccAddress) {
	bw := m.GetCurrentAccBandwidth(ctx, address)
	m.bwKeeper.SetAccBandwidth(ctx, bw)
}

//
// Performs bw consumption for given acc
// To get right number, should be called after tx delivery with bw state obtained prior delivery
//
// Pseudo code:
// bw := getCurrentBw(addr)
// bwCost := deliverTx(tx)
// consumeBw(bw, bwCost)
func (m *BaseBandwidthMeter) ConsumeAccBandwidth(ctx sdk.Context, bw types.AcсBandwidth, amt int64) {
	bw.Consume(amt)
	m.bwKeeper.SetAccBandwidth(ctx, bw)
	bw = m.GetCurrentAccBandwidth(ctx, bw.Address)
	m.bwKeeper.SetAccBandwidth(ctx, bw)
}

func (m *BaseBandwidthMeter) GetCurrentCreditPrice() float64 {
	return m.currentCreditPrice
}

func (m *BaseBandwidthMeter) GetParamSet(ctx sdk.Context) Params {
	subspace, ok := m.paramsKeeper.GetSubspace(DefaultParamspace)
	if !ok {
		panic("bandwidth params subspace is not found")
	}
	var paramset Params
	subspace.GetParamSet(ctx, &paramset)

	return paramset
}
