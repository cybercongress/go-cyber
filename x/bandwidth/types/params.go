package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"strings"
)

// Default parameter values
const (
	DefaultLinkMsgCost                                int64   = 100
	DefaultRecoveryPeriod                             int64   = 18000
	DefaultAdjustPricePeriod                          int64   = 10
	DefaultBaseCreditPrice                            string = "1.0"
	DefaultDesirableNetworkBandwidthForRecoveryPeriod int64   = 200000000
	DefaultTxCost                                     int64   = 300
	DefaultNonLinkMsgCost                             int64   = 500
	DefaultSlidingWindowSize                          int64   = 18000
	DefaultShouldBeSpentPerSlidingWindow              string = "200000000"
)

// Parameter keys
var (
	// Bandwidth cost of specific messages and tx itself
	KeyLinkMsgCost = []byte("LinkMsgCost")
	// Number of blocks to recover full bandwidth
	KeyRecoveryPeriod = []byte("RecoveryPeriod")
	// Number of blocks before next adjust price
	KeyAdjustPricePeriod = []byte("AdjustPricePeriod")
	KeyBaseCreditPrice   = []byte("BaseCreditPrice")
	// Maximum bandwidth of network
	KeyDesirableNetworkBandwidthForRecoveryPeriod = []byte("DesirableNetworkBandwidthForRecoveryPeriod")
	KeyTxCost                                     = []byte("TxCost")
	KeyNonLinkMsgCost                             = []byte("NonLinkMsgCost")
	KeySlidingWindowSize                          = []byte("SlidingWindowSize")
	KeyShouldBeSpentPerSlidingWindow              = []byte("ShouldBeSpentPerSlidingWindow")
)

// Params defines the parameters for the bandwidth module.
//BaseCreditPrice and ShouldBeSpentPerSlidingWindow are strings because
// `amino:"unsafe"` tag is not working for now:
// https://github.com/tendermint/go-amino/issues/230
type Params struct {
	LinkMsgCost                                int64
	RecoveryPeriod                             int64
	AdjustPricePeriod                          int64
	BaseCreditPrice                            string
	DesirableNetworkBandwidthForRecoveryPeriod int64
	TxCost                                     int64
	NonLinkMsgCost                             int64
	SlidingWindowSize                          int64
	ShouldBeSpentPerSlidingWindow              string
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of bandwidth module's parameters.
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{KeyLinkMsgCost, &p.LinkMsgCost},
		{KeyRecoveryPeriod, &p.RecoveryPeriod},
		{KeyAdjustPricePeriod, &p.AdjustPricePeriod},
		{KeyBaseCreditPrice, &p.BaseCreditPrice},
		{KeyDesirableNetworkBandwidthForRecoveryPeriod, &p.DesirableNetworkBandwidthForRecoveryPeriod},
		{KeyTxCost, &p.TxCost},
		{KeyNonLinkMsgCost, &p.NonLinkMsgCost},
		{KeySlidingWindowSize, &p.SlidingWindowSize},
		{KeyShouldBeSpentPerSlidingWindow, &p.ShouldBeSpentPerSlidingWindow},
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("LinkMsgCost: %d\n", p.LinkMsgCost))
	sb.WriteString(fmt.Sprintf("RecoveryPeriod: %d\n", p.RecoveryPeriod))
	sb.WriteString(fmt.Sprintf("AdjustPricePeriod: %d\n", p.AdjustPricePeriod))
	sb.WriteString(fmt.Sprintf("BaseCreditPrice: %d\n", p.BaseCreditPrice))
	sb.WriteString(fmt.Sprintf("DesirableNetworkBandwidthForRecoveryPeriod: %d\n", p.DesirableNetworkBandwidthForRecoveryPeriod))
	sb.WriteString(fmt.Sprintf("TxCost: %d\n", p.TxCost))
	sb.WriteString(fmt.Sprintf("NonLinkMsgCost: %d\n", p.NonLinkMsgCost))
	sb.WriteString(fmt.Sprintf("SlidingWindowSize: %d\n", p.SlidingWindowSize))
	sb.WriteString(fmt.Sprintf("ShouldBeSpentPerSlidingWindow: %d\n", p.ShouldBeSpentPerSlidingWindow))

	return sb.String()
}

// NewParams creates a new Params object
func NewParams(
	linkMsgCost int64,
	recoveryPeriod int64,
	adjustPricePeriod int64,
	baseCreditPrice string,
	desirableNetworkBandwidthForRecoveryPeriod int64,
	txCost int64,
	nonLinkMsgCost int64,
	slidingWindowSize int64,
	shouldBeSpentPerSlidingWindow string) Params {

	return Params{
		LinkMsgCost:       linkMsgCost,
		RecoveryPeriod:    recoveryPeriod,
		AdjustPricePeriod: adjustPricePeriod,
		BaseCreditPrice:   baseCreditPrice,
		DesirableNetworkBandwidthForRecoveryPeriod: desirableNetworkBandwidthForRecoveryPeriod,
		TxCost:                        txCost,
		NonLinkMsgCost:                nonLinkMsgCost,
		SlidingWindowSize:             slidingWindowSize,
		ShouldBeSpentPerSlidingWindow: shouldBeSpentPerSlidingWindow,
	}
}

// NewDefaultParams returns a default set of parameters.
func NewDefaultParams() Params {
	return Params{
		LinkMsgCost:       DefaultLinkMsgCost,
		RecoveryPeriod:    DefaultRecoveryPeriod,
		AdjustPricePeriod: DefaultAdjustPricePeriod,
		BaseCreditPrice:   DefaultBaseCreditPrice,
		DesirableNetworkBandwidthForRecoveryPeriod: DefaultDesirableNetworkBandwidthForRecoveryPeriod,
		TxCost:                        DefaultTxCost,
		NonLinkMsgCost:                DefaultNonLinkMsgCost,
		SlidingWindowSize:             DefaultSlidingWindowSize,
		ShouldBeSpentPerSlidingWindow: DefaultShouldBeSpentPerSlidingWindow,
	}
}

// ParamKeyTable for bandwidth module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}
