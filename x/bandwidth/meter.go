package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/go-cyber/store"
	"github.com/cybercongress/go-cyber/x/bandwidth/internal/types"
	"math"
	"strconv"
)

var _ types.BandwidthMeter = &BaseBandwidthMeter{}

type BaseBandwidthMeter struct {
	// data providers
	accountKeeper           auth.AccountKeeper
	stakeProvider           types.AccStakeProvider
	accountBaindwidthKeeper AccountBandwidthKeeper
	mainKeeper              store.MainKeeper
	blockBandwidthKeeper    BlockSpentBandwidthKeeper
	// bw configuration
	msgCost types.MsgBandwidthCost

	// price adjustment fields
	curBlockSpentBandwidth     uint64 //resets every block
	currentCreditPrice         float64
	bandwidthSpent             map[uint64]uint64 // bandwidth spent by blocks
	totalSpentForSlidingWindow uint64
	currentBlockSpentKarma     uint64
}

func NewBaseMeter(
	mk store.MainKeeper, ak auth.AccountKeeper, bwk AccountBandwidthKeeper,
	bbwk BlockSpentBandwidthKeeper, sp types.AccStakeProvider, msgCost types.MsgBandwidthCost,
) *BaseBandwidthMeter {

	return &BaseBandwidthMeter{
		mainKeeper:              mk,
		blockBandwidthKeeper:    bbwk,
		accountKeeper:           ak,
		accountBaindwidthKeeper: bwk,
		stakeProvider:           sp,
		msgCost:                 msgCost,
		bandwidthSpent:          make(map[uint64]uint64),
	}
}

func (m *BaseBandwidthMeter) Load(ctx sdk.Context) {
	params := m.accountBaindwidthKeeper.GetParams(ctx)
	m.totalSpentForSlidingWindow = 0
	m.bandwidthSpent = m.blockBandwidthKeeper.GetValuesForPeriod(ctx, params.RecoveryPeriod)
	for _, spentBandwidth := range m.bandwidthSpent {
		m.totalSpentForSlidingWindow += spentBandwidth
	}
	floatBaseCreditPrice, err := strconv.ParseFloat(params.BaseCreditPrice.String(), 64)
	if err != nil {
		panic(err)
	}
	m.currentCreditPrice = 0.01 * math.Float64frombits(m.mainKeeper.GetBandwidthPrice(ctx, floatBaseCreditPrice))
	m.curBlockSpentBandwidth = 0
}

func (m *BaseBandwidthMeter) AddToBlockBandwidth(value int64) {
	m.curBlockSpentBandwidth += uint64(value)
}

func (m *BaseBandwidthMeter) AddToBlockKarma(value int64) {
	m.currentBlockSpentKarma += uint64(value)
}

// Here we move bandwidth window:
// Remove first block of window and add new block to window end
func (m *BaseBandwidthMeter) CommitBlockBandwidth(ctx sdk.Context) {
	m.totalSpentForSlidingWindow += m.curBlockSpentBandwidth

	newWindowEnd := ctx.BlockHeight()
	params := m.accountBaindwidthKeeper.GetParams(ctx)
	windowStart := newWindowEnd - params.RecoveryPeriod
	if windowStart < 0 { // check needed cause it will be casted to uint and can cause overflow
		windowStart = 0
	}
	// If recovery period will be increased via governance extended windows will not be accessible
	// todo If recover period will be decreased via governance need to clean garbage values
	windowStartValue, exists := m.bandwidthSpent[uint64(windowStart)]
	if exists {
		m.totalSpentForSlidingWindow -= windowStartValue
		delete(m.bandwidthSpent, uint64(windowStart))
	}
	m.blockBandwidthKeeper.SetBlockSpentBandwidth(ctx, uint64(ctx.BlockHeight()), m.curBlockSpentBandwidth)
	m.bandwidthSpent[uint64(newWindowEnd)] = m.curBlockSpentBandwidth
	m.curBlockSpentBandwidth = 0
}

func (m *BaseBandwidthMeter) CommitTotalKarma(ctx sdk.Context) {
	if m.currentBlockSpentKarma != uint64(0) {
		currentKarma := m.mainKeeper.GetSpentKarma(ctx)
		m.mainKeeper.StoreSpentKarma(ctx, currentKarma + m.currentBlockSpentKarma)
	}
	m.currentBlockSpentKarma = 0
}

func (m *BaseBandwidthMeter) AdjustPrice(ctx sdk.Context) {
	params := m.accountBaindwidthKeeper.GetParams(ctx)
	floatBaseCreditPrice, err := strconv.ParseFloat(params.BaseCreditPrice.String(), 64)
	if err != nil {
		panic(err)
	}
	newPrice := float64(m.totalSpentForSlidingWindow) / float64(params.DesirableBandwidth)

	if newPrice < 0.01 * floatBaseCreditPrice {
		newPrice = 0.01 * floatBaseCreditPrice
	}

	m.currentCreditPrice = newPrice
	m.mainKeeper.StoreBandwidthPrice(ctx, math.Float64bits(newPrice))
}

func (m *BaseBandwidthMeter) GetTxCost(ctx sdk.Context, tx sdk.Tx) int64 {
	params := m.accountBaindwidthKeeper.GetParams(ctx)
	bandwidthForTx := params.TxCost
	for _, msg := range tx.GetMsgs() {
		bandwidthForTx = bandwidthForTx + m.msgCost(ctx, params, msg)
	}
	return bandwidthForTx
}

func (m *BaseBandwidthMeter) GetMaxBlockBandwidth(ctx sdk.Context) uint64 {
	params := m.accountBaindwidthKeeper.GetParams(ctx)
	maxBlockBandwidth := params.MaxBlockBandwidth
	return maxBlockBandwidth
}

func (m *BaseBandwidthMeter) GetPricedTxCost(ctx sdk.Context, tx sdk.Tx) int64 {
	return int64(float64(m.GetTxCost(ctx, tx)) * m.currentCreditPrice)
}

func (m *BaseBandwidthMeter) GetCurBlockSpentBandwidth(ctx sdk.Context) uint64 {
	return m.curBlockSpentBandwidth
}

func (m *BaseBandwidthMeter) GetPricedLinksCost(ctx sdk.Context, tx sdk.Tx) int64 {
	params := m.accountBaindwidthKeeper.GetParams(ctx)
	usedBandwidth := int64(0)
	for _, msg := range tx.GetMsgs() {
		if msg.Type() == "link" {
			usedBandwidth = usedBandwidth + m.msgCost(ctx, params, msg)
		}
	}
	return int64(float64(usedBandwidth) * m.currentCreditPrice)
}

func (m *BaseBandwidthMeter) GetAccMaxBandwidth(ctx sdk.Context, addr sdk.AccAddress) int64 {
	accStakePercentage := m.stakeProvider.GetAccStakePercentage(ctx, addr)
	params := m.accountBaindwidthKeeper.GetParams(ctx)
	return int64(accStakePercentage * float64(params.DesirableBandwidth))
}

func (m *BaseBandwidthMeter) GetCurrentAccBandwidth(ctx sdk.Context, address sdk.AccAddress) types.AcсountBandwidth {
	accBw := m.accountBaindwidthKeeper.GetAccountBandwidth(ctx, address)
	accMaxBw := m.GetAccMaxBandwidth(ctx, address)
	params := m.accountBaindwidthKeeper.GetParams(ctx)
	accBw.UpdateMax(accMaxBw, ctx.BlockHeight(), params.RecoveryPeriod)
	return accBw
}

func (m *BaseBandwidthMeter) UpdateAccMaxBandwidth(ctx sdk.Context, address sdk.AccAddress) {
	bw := m.GetCurrentAccBandwidth(ctx, address)
	m.accountBaindwidthKeeper.SetAccountBandwidth(ctx, bw)
}

//
// Performs bw consumption for given acc
// To get right number, should be called after tx delivery with bw state obtained prior delivery
//
// Pseudo code:
// bw := getCurrentBw(addr)
// bwCost := deliverTx(tx)
// consumeBw(bw, bwCost)
func (m *BaseBandwidthMeter) ConsumeAccBandwidth(ctx sdk.Context, bw types.AcсountBandwidth, amt int64) {
	bw.Consume(amt)
	m.accountBaindwidthKeeper.SetAccountBandwidth(ctx, bw)
	bw = m.GetCurrentAccBandwidth(ctx, bw.Address)
	m.accountBaindwidthKeeper.SetAccountBandwidth(ctx, bw)
}

func (m *BaseBandwidthMeter) UpdateLinkedBandwidth(ctx sdk.Context, bw types.AcсountBandwidth, amt int64) {
	bw.AddLinked(amt)
	m.accountBaindwidthKeeper.SetAccountBandwidth(ctx, bw)
}


func (m *BaseBandwidthMeter) GetCurrentCreditPrice() float64 {
	return m.currentCreditPrice
}
