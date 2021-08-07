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
	cdc           codec.BinaryMarshaler
	storeKey      sdk.StoreKey
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
		cdc:                   cdc,
		storeKey:              key,
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

func (m *BandwidthMeter) LoadState(ctx sdk.Context) {
	params := m.GetParams(ctx)
	m.totalSpentForSlidingWindow = 0
	m.bandwidthSpentByBlock = m.GetValuesForPeriod(ctx, params.RecoveryPeriod)
	for _, spentBandwidth := range m.bandwidthSpentByBlock {
		m.totalSpentForSlidingWindow += spentBandwidth
	}
	m.currentCreditPrice = m.GetBandwidthPrice(ctx, params.BasePrice)
	m.currentBlockSpentBandwidth = 0
}

func (bm BandwidthMeter) GetBandwidthPrice(ctx sdk.Context, basePrice sdk.Dec) sdk.Dec {
	store := ctx.KVStore(bm.storeKey)
	priceAsBytes := store.Get(types.LastBandwidthPrice)
	if priceAsBytes == nil {
		return basePrice
	}
	var price types.Price
	bm.cdc.MustUnmarshalBinaryBare(priceAsBytes, &price)
	return price.Price
}

func (bm BandwidthMeter) StoreBandwidthPrice(ctx sdk.Context, price sdk.Dec) {
	store := ctx.KVStore(bm.storeKey)
	store.Set(types.LastBandwidthPrice, bm.cdc.MustMarshalBinaryBare(&types.Price{Price: price}))
}

func (bm BandwidthMeter) GetDesirableBandwidth(ctx sdk.Context) uint64 {
	store := ctx.KVStore(bm.storeKey)
	bandwidthAsBytes := store.Get(types.DesirableBandwidth)
	if bandwidthAsBytes == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bandwidthAsBytes)
}

func (bm BandwidthMeter) StoreDesirableBandwidth(ctx sdk.Context, bandwidth uint64) {
	store := ctx.KVStore(bm.storeKey)
	store.Set(types.DesirableBandwidth, sdk.Uint64ToBigEndian(bandwidth))
}

func (bm BandwidthMeter) AddToDesirableBandwidth(ctx sdk.Context, toAdd uint64) {
	current := bm.GetDesirableBandwidth(ctx)
	store := ctx.KVStore(bm.storeKey)
	store.Set(types.DesirableBandwidth, sdk.Uint64ToBigEndian(current+toAdd))
}

func (m *BandwidthMeter) AddToBlockBandwidth(value uint64) {
	m.currentBlockSpentBandwidth += value
}

// Here we move bandwidth window:
// Remove first block of window and add new block to window end
func (m *BandwidthMeter) CommitBlockBandwidth(ctx sdk.Context) {
	params := m.GetParams(ctx)
	defer func() {
		if m.currentBlockSpentBandwidth > 0 {
			m.Logger(ctx).Info("Block", "bandwidth", m.currentBlockSpentBandwidth)
			m.Logger(ctx).Info("Window", "bandwidth", m.totalSpentForSlidingWindow)
		}

		telemetry.SetGauge(float32(m.currentBlockSpentBandwidth), types.ModuleName, "block_bandwidth")
		telemetry.SetGauge(float32(m.totalSpentForSlidingWindow), types.ModuleName, "window_bandwidth")
		m.currentBlockSpentBandwidth = 0
	}()

	m.totalSpentForSlidingWindow += m.currentBlockSpentBandwidth

	newWindowEnd := ctx.BlockHeight()
	windowStart := newWindowEnd - int64(params.RecoveryPeriod)
	if windowStart < 0 {
		windowStart = 0
	}

	// clean window slot in in-memory
	windowStartValue, exists := m.bandwidthSpentByBlock[uint64(windowStart)]
	if exists {
		m.totalSpentForSlidingWindow -= windowStartValue
		delete(m.bandwidthSpentByBlock, uint64(windowStart))
	}

	// clean window slot in storage
	store := ctx.KVStore(m.storeKey)
	if store.Has(types.BlockStoreKey(uint64(windowStart))) {
		store.Delete(types.BlockStoreKey(uint64(windowStart)))
	}

	m.SetBlockBandwidth(ctx, uint64(ctx.BlockHeight()), m.currentBlockSpentBandwidth)
	m.bandwidthSpentByBlock[uint64(newWindowEnd)] = m.currentBlockSpentBandwidth
}

func (m *BandwidthMeter) GetCurrentBlockSpentBandwidth(ctx sdk.Context) uint64 {
	return m.currentBlockSpentBandwidth
}

func (m *BandwidthMeter) GetCurrentNetworkLoad(ctx sdk.Context) sdk.Dec {
	return sdk.NewDec(int64(m.totalSpentForSlidingWindow)).QuoInt64(int64(m.GetDesirableBandwidth(ctx)))
}

func (m *BandwidthMeter) GetMaxBlockBandwidth(ctx sdk.Context) uint64 {
	params := m.GetParams(ctx)
	maxBlockBandwidth := params.MaxBlockBandwidth
	return maxBlockBandwidth
}

func (m *BandwidthMeter) GetCurrentCreditPrice() sdk.Dec {
	return m.currentCreditPrice
}

func (m *BandwidthMeter) AdjustPrice(ctx sdk.Context) {
	params := m.GetParams(ctx)

	desirableBandwidth := m.GetDesirableBandwidth(ctx)
	if desirableBandwidth != 0 {
		telemetry.SetGauge(float32(m.totalSpentForSlidingWindow)/float32(desirableBandwidth), types.ModuleName, "load")

		newPrice := sdk.NewDec(int64(m.totalSpentForSlidingWindow)).QuoInt64(int64(desirableBandwidth))
		m.Logger(ctx).Info("Load", "value", newPrice.String())
		if newPrice.LT(params.BasePrice) {
			newPrice = params.BasePrice
		}
		m.Logger(ctx).Info("Price", "value", newPrice.String())

		m.currentCreditPrice = newPrice
		m.StoreBandwidthPrice(ctx, newPrice)
	}
}

func (m *BandwidthMeter) GetTotalCyberlinksCost(ctx sdk.Context, tx sdk.Tx) (uint64) {
	bandwidthForTx := uint64(0)
	for _, msg := range tx.GetMsgs() {
		linkMsg := msg.(*gtypes.MsgCyberlink)
		bandwidthForTx = bandwidthForTx + uint64(len(linkMsg.Links)) * 1000
	}
	return bandwidthForTx
}

func (m *BandwidthMeter) GetPricedTotalCyberlinksCost(ctx sdk.Context, tx sdk.Tx) uint64 {
	return uint64(m.currentCreditPrice.Mul(sdk.NewDec(int64(m.GetTotalCyberlinksCost(ctx, tx)))).RoundInt64())
}

// Performs bw consumption for given acc
// To get right number, should be called after tx delivery with bw state obtained prior delivery
// Pseudo code:
// bw := getCurrentBw(addr)
// bwCost := deliverTx(tx)
// consumeBw(bw, bwCost)
func (m *BandwidthMeter) ConsumeAccountBandwidth(ctx sdk.Context, bw types.AccountBandwidth, amt uint64) error {
	err := bw.Consume(amt); if err != nil {
		return err
	}
	m.SetAccountBandwidth(ctx, bw)
	// TODO test and remove next lines
	//addr, _ := sdk.AccAddressFromBech32(bw.Address)
	//bw = m.GetCurrentAccountBandwidth(ctx, addr)
	//m.SetAccountBandwidth(ctx, bw)
	return nil
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
	return uint64(accStakePercentage * float64(m.GetDesirableBandwidth(ctx)))
}

func (m *BandwidthMeter) UpdateAccountMaxBandwidth(ctx sdk.Context, address sdk.AccAddress) {
	bw := m.GetCurrentAccountBandwidth(ctx, address)
	m.SetAccountBandwidth(ctx, bw)
}
