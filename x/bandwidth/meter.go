package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/cyberd/store"
	"github.com/cybercongress/cyberd/x/bandwidth/types"
	"math"
)

var _ types.BandwidthMeter = &BaseBandwidthMeter{}

type BaseBandwidthMeter struct {
	// data providers
	accKeeper            auth.AccountKeeper
	stakeProvider        types.AccStakeProvider
	bwKeeper             types.Keeper
	mainKeeper           store.MainKeeper
	blockBandwidthKeeper types.BlockSpentBandwidthKeeper

	// bw configuration
	msgCost types.MsgBandwidthCost

	// price adjustment fields
	curBlockSpentBandwidth     uint64 //resets every block
	currentCreditPrice         float64
	bandwidthSpent             map[uint64]uint64 // bandwidth spent by blocks
	totalSpentForSlidingWindow uint64
}

func NewBaseMeter(
	mainKeeper store.MainKeeper, ak auth.AccountKeeper, sp types.AccStakeProvider, bwKeeper types.Keeper,
	msgCost types.MsgBandwidthCost, blockBandwidthKeeper types.BlockSpentBandwidthKeeper,
) *BaseBandwidthMeter {

	return &BaseBandwidthMeter{
		mainKeeper:           mainKeeper,
		blockBandwidthKeeper: blockBandwidthKeeper,
		accKeeper:            ak,
		stakeProvider:        sp,
		bwKeeper:             bwKeeper,
		msgCost:              msgCost,
		bandwidthSpent:       make(map[uint64]uint64),
	}
}

func (m *BaseBandwidthMeter) Load(ctx sdk.Context) {
	m.totalSpentForSlidingWindow = 0
	m.bandwidthSpent = m.blockBandwidthKeeper.GetValuesForPeriod(ctx, SlidingWindowSize)
	for _, spentBandwidth := range m.bandwidthSpent {
		m.totalSpentForSlidingWindow += spentBandwidth
	}
	m.currentCreditPrice = math.Float64frombits(m.mainKeeper.GetBandwidthPrice(ctx, BaseCreditPrice))
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
	windowStart := newWindowEnd - SlidingWindowSize
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

	newPrice := float64(m.totalSpentForSlidingWindow) / float64(ShouldBeSpentPerSlidingWindow)

	if newPrice < 0.01*BaseCreditPrice {
		newPrice = 0.01 * BaseCreditPrice
	}

	m.currentCreditPrice = newPrice
	m.mainKeeper.StoreBandwidthPrice(ctx, math.Float64bits(newPrice))
}

func (m *BaseBandwidthMeter) GetTxCost(tx sdk.Tx) int64 {
	bandwidthForTx := TxCost
	for _, msg := range tx.GetMsgs() {
		bandwidthForTx = bandwidthForTx + m.msgCost(msg)
	}
	return bandwidthForTx
}

func (m *BaseBandwidthMeter) GetPricedTxCost(tx sdk.Tx) int64 {
	return int64(float64(m.GetTxCost(tx)) * m.currentCreditPrice)
}

func (m *BaseBandwidthMeter) GetAccMaxBandwidth(ctx sdk.Context, addr sdk.AccAddress) int64 {
	accStakePercentage := m.stakeProvider.GetAccStakePercentage(ctx, addr)
	return int64(accStakePercentage * float64(DesirableNetworkBandwidthForRecoveryPeriod))
}

func (m *BaseBandwidthMeter) GetCurrentAccBandwidth(ctx sdk.Context, address sdk.AccAddress) types.AcсBandwidth {
	accBw := m.bwKeeper.GetAccBandwidth(ctx, address)
	accMaxBw := m.GetAccMaxBandwidth(ctx, address)
	accBw.UpdateMax(accMaxBw, ctx.BlockHeight(), RecoveryPeriod)
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
