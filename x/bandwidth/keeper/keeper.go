package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cybercongress/go-cyber/x/bandwidth/types"
	gtypes "github.com/cybercongress/go-cyber/x/graph/types"
)

type BandwidthMeter struct {
	stakeProvider types.AccountStakeProvider
	Cdc           codec.BinaryMarshaler
	StoreKey      sdk.StoreKey
	paramSpace    paramstypes.Subspace

	currentBlockSpentBandwidth uint64
	currentCreditPrice         sdk.Dec
	bandwidthSpentByBlock      map[uint64]uint64
	totalSpentForSlidingWindow uint64
}

func NewBandwidthMeter(
	cdc codec.BinaryMarshaler,
	key sdk.StoreKey,
	asp types.AccountStakeProvider,
	paramSpace paramstypes.Subspace,
) *BandwidthMeter {

	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &BandwidthMeter{
		Cdc:                   cdc,
		StoreKey:              key,
		stakeProvider:         asp,
		paramSpace:            paramSpace,
		bandwidthSpentByBlock: make(map[uint64]uint64),
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

func (m *BandwidthMeter) LoadState(ctx sdk.Context) {
	params := m.GetParams(ctx)
	m.totalSpentForSlidingWindow = 0
	m.bandwidthSpentByBlock = m.GetValuesForPeriod(ctx, params.RecoveryPeriod)
	for _, spentBandwidth := range m.bandwidthSpentByBlock {
		m.totalSpentForSlidingWindow += spentBandwidth
	}
	m.currentCreditPrice = m.GetBandwidthPrice(ctx, params.BaseCreditPrice)
	m.currentBlockSpentBandwidth = 0
}

//____________________________________________________________________________

func (bm BandwidthMeter) GetBandwidthPrice(ctx sdk.Context, basePrice sdk.Dec) sdk.Dec {
	store := ctx.KVStore(bm.StoreKey)
	priceAsBytes := store.Get(types.LastBandwidthPrice)
	if priceAsBytes == nil {
		return basePrice
	}
	var price types.Price
	bm.Cdc.MustUnmarshalBinaryBare(priceAsBytes, &price)
	return price.Price
}

func (bm BandwidthMeter) StoreBandwidthPrice(ctx sdk.Context, price sdk.Dec) {
	store := ctx.KVStore(bm.StoreKey)
	store.Set(types.LastBandwidthPrice, bm.Cdc.MustMarshalBinaryBare(&types.Price{Price: price}))
}

//____________________________________________________________________________

func (m *BandwidthMeter) AddToBlockBandwidth(value uint64) {
	m.currentBlockSpentBandwidth += value
}

// Here we move bandwidth window:
// Remove first block of window and add new block to window end
func (m *BandwidthMeter) CommitBlockBandwidth(ctx sdk.Context) {
	params := m.GetParams(ctx)
	defer func() {
		telemetry.SetGauge((float32(m.totalSpentForSlidingWindow)/float32(params.DesirableBandwidth)), types.ModuleName, "load")
		telemetry.SetGauge(float32(m.currentBlockSpentBandwidth), types.ModuleName, "block_bandwidth")
		telemetry.SetGauge(float32(m.totalSpentForSlidingWindow), types.ModuleName, "window_bandwidth")
		m.currentBlockSpentBandwidth = 0
	}()

	m.totalSpentForSlidingWindow += m.currentBlockSpentBandwidth

	newWindowEnd := ctx.BlockHeight()
	//params := m.GetParams(ctx)
	windowStart := newWindowEnd - int64(params.RecoveryPeriod)
	if windowStart < 0 {
		windowStart = 0
	}

	windowStartValue, exists := m.bandwidthSpentByBlock[uint64(windowStart)]
	if exists {
		m.totalSpentForSlidingWindow -= windowStartValue
		delete(m.bandwidthSpentByBlock, uint64(windowStart))
	}
	m.SetBlockBandwidth(ctx, uint64(ctx.BlockHeight()), m.currentBlockSpentBandwidth)
	m.bandwidthSpentByBlock[uint64(newWindowEnd)] = m.currentBlockSpentBandwidth
}

func (m *BandwidthMeter) GetCurrentBlockSpentBandwidth(ctx sdk.Context) uint64 {
	return m.currentBlockSpentBandwidth
}

func (m *BandwidthMeter) GetCurrentNetworkLoad(ctx sdk.Context) sdk.Dec {
	params := m.GetParams(ctx)

	return sdk.NewDec(int64(m.totalSpentForSlidingWindow)).QuoInt64(int64(params.DesirableBandwidth))
}

func (m *BandwidthMeter) GetMaxBlockBandwidth(ctx sdk.Context) uint64 {
	params := m.GetParams(ctx)
	maxBlockBandwidth := params.MaxBlockBandwidth
	return maxBlockBandwidth
}

//____________________________________________________________________________

func (m *BandwidthMeter) GetCurrentCreditPrice() sdk.Dec {
	return m.currentCreditPrice
}

func (m *BandwidthMeter) AdjustPrice(ctx sdk.Context) {
	params := m.GetParams(ctx)

	newPrice := sdk.NewDec(int64(m.totalSpentForSlidingWindow)).QuoInt64(int64(params.DesirableBandwidth))
	if newPrice.LT(params.BaseCreditPrice) {
		newPrice = params.BaseCreditPrice
	}

	m.currentCreditPrice = newPrice
	m.StoreBandwidthPrice(ctx, newPrice)
}

//____________________________________________________________________________

func (m *BandwidthMeter) GetTxCost(ctx sdk.Context, tx sdk.Tx) uint64 {
	params := m.GetParams(ctx)
	bandwidthForTx := params.TxCost
	for _, msg := range tx.GetMsgs() {
		linkMsg := msg.(*gtypes.MsgCyberlink)
		bandwidthForTx = bandwidthForTx + uint64(len(linkMsg.Links)) * params.LinkCost
	}
	return bandwidthForTx
}

func (m *BandwidthMeter) GetPricedTxCost(ctx sdk.Context, tx sdk.Tx) uint64 {
	return uint64(m.currentCreditPrice.Mul(sdk.NewDec(int64(m.GetTxCost(ctx, tx)))).RoundInt64())
}

//____________________________________________________________________________

// Performs bw consumption for given acc
// To get right number, should be called after tx delivery with bw state obtained prior delivery
// Pseudo code:
// bw := getCurrentBw(addr)
// bwCost := deliverTx(tx)
// consumeBw(bw, bwCost)
func (m *BandwidthMeter) ConsumeAccountBandwidth(ctx sdk.Context, bw types.AccountBandwidth, amt uint64) {
	bw.Consume(amt)
	m.SetAccountBandwidth(ctx, bw)
	addr, _ := sdk.AccAddressFromBech32(bw.Address)
	bw = m.GetCurrentAccountBandwidth(ctx, addr)
	m.SetAccountBandwidth(ctx, bw)
}

func (m *BandwidthMeter) GetCurrentAccountBandwidth(ctx sdk.Context, address sdk.AccAddress) types.AccountBandwidth {
	accBw := m.GetAccountBandwidth(ctx, address)
	accMaxBw := m.GetAccountMaxBandwidth(ctx, address)
	params := m.GetParams(ctx)
	accBw.UpdateMax(accMaxBw, uint64(ctx.BlockHeight()), params.RecoveryPeriod)
	return accBw
}

func (m *BandwidthMeter) GetAccountMaxBandwidth(ctx sdk.Context, addr sdk.AccAddress) uint64 {
	accStakePercentage := m.stakeProvider.GetAccountStakePercentageVolt(ctx, addr)
	params := m.GetParams(ctx)
	return uint64(accStakePercentage * float64(params.DesirableBandwidth))
}

func (m *BandwidthMeter) UpdateAccountMaxBandwidth(ctx sdk.Context, address sdk.AccAddress) {
	bw := m.GetCurrentAccountBandwidth(ctx, address)
	m.SetAccountBandwidth(ctx, bw)
}

//____________________________________________________________________________
