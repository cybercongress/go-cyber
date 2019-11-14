package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// Default parameter values
const (
	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName

	MinLinkMsgCost       = 1
	MinRecoveryPeriod    = 100
	MinAdjustPricePeriod = 1
	MinTxCost            = 1
	MinNonLinkMsgCost    = 1
	MinDesirableBandwidth = 10000
	MinBlockBandwidth    = 100
)

// Parameter keys
var (
	// Bandwidth cost of specific messages and tx itself
	KeyLinkMsgCost 		  = []byte("LinkMsgCost")
	// Number of blocks to recover full bandwidth
	KeyRecoveryPeriod     = []byte("RecoveryPeriod")
	// Number of blocks before next adjust price
	KeyAdjustPricePeriod  = []byte("AdjustPricePeriod")
	KeyBaseCreditPrice    = []byte("BaseCreditPrice")
	// Maximum bandwidth of network
	KeyDesirableBandwidth = []byte("DesirableBandwidth")
	KeyMaxBlockBandwidth = []byte("MaxBlockBandwidth")
	KeyTxCost             = []byte("TxCost")
	KeyNonLinkMsgCost     = []byte("NonLinkMsgCost")
)

// Params defines the parameters for the bandwidth module.
type Params struct {
	LinkMsgCost                   int64  `json:"link_msg_cost" yaml:"link_msg_cost"`
	RecoveryPeriod                int64  `json:"recovery_period" yaml:"recovery_period"`
	AdjustPricePeriod             int64  `json:"adjust_price_period" yaml:"adjust_price_period"`
	BaseCreditPrice               sdk.Dec `json:"base_credit_price" yaml:"base_credit_price"`
	DesirableBandwidth            int64  `json:"desirable_bandwidth" yaml:"desirable_bandwidth"`
	MaxBlockBandwidth			  uint64  `json:"max_block_bandwidth" yaml:"max_block_bandwidth"`
	TxCost                        int64  `json:"tx_cost" yaml:"tx_cost"`
	NonLinkMsgCost                int64  `json:"non_link_msg_cost" yaml:"non_link_msg_cost"`
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of bandwidth module's parameters.
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{KeyLinkMsgCost, &p.LinkMsgCost},
		{KeyRecoveryPeriod, &p.RecoveryPeriod},
		{KeyAdjustPricePeriod, &p.AdjustPricePeriod},
		{KeyBaseCreditPrice, &p.BaseCreditPrice},
		{KeyDesirableBandwidth, &p.DesirableBandwidth},
		{KeyMaxBlockBandwidth, &p.MaxBlockBandwidth},
		{KeyTxCost, &p.TxCost},
		{KeyNonLinkMsgCost, &p.NonLinkMsgCost},
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
	sb.WriteString(fmt.Sprintf("DesirableBandwidth: %d\n", p.DesirableBandwidth))
	sb.WriteString(fmt.Sprintf("DesirableBandwidth: %d\n", p.DesirableBandwidth))
	sb.WriteString(fmt.Sprintf("TxCost: %d\n", p.TxCost))
	sb.WriteString(fmt.Sprintf("NonLinkMsgCost: %d\n", p.NonLinkMsgCost))

	return sb.String()
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if p.LinkMsgCost < MinLinkMsgCost {
		return fmt.Errorf("invalid link msg cost: %d, can not be less then %d", p.LinkMsgCost, MinLinkMsgCost)
	}
	if p.RecoveryPeriod < MinRecoveryPeriod {
		return fmt.Errorf("invalid recovery period: %d, can not be less then %d", p.RecoveryPeriod, MinRecoveryPeriod)
	}
	if p.AdjustPricePeriod < MinAdjustPricePeriod {
		return fmt.Errorf("invalid adjust price period: %d, can not be less then %d", p.AdjustPricePeriod, MinAdjustPricePeriod)
	}
	if p.BaseCreditPrice.LT(sdk.OneDec()) {
		return fmt.Errorf("base credit price parameter must be >= 1, is %s", p.BaseCreditPrice)
	}
	if p.BaseCreditPrice.GT(sdk.NewDec(100)) {
		return fmt.Errorf("base credit price parameter must be <= 100, is %s", p.BaseCreditPrice)
	}
	if p.DesirableBandwidth < MinDesirableBandwidth {
		return fmt.Errorf("invalid desirable bandwidth: %d, can not be less then %d", p.DesirableBandwidth, MinDesirableBandwidth)
	}
	if p.MaxBlockBandwidth < MinBlockBandwidth {
		return fmt.Errorf("invalid max block bandwidth: %d, can not be less then %d", p.MaxBlockBandwidth, MinBlockBandwidth)
	}
	if p.TxCost < MinTxCost {
		return fmt.Errorf("invalid tx cost: %d, can not be less then %d", p.TxCost, MinTxCost)
	}
	if p.NonLinkMsgCost < MinNonLinkMsgCost {
		return fmt.Errorf("invalid non link msg cost: %d, can not be less then %d", p.NonLinkMsgCost, MinNonLinkMsgCost)
	}
	return nil
}

// NewParams creates a new Params object
func NewParams(
	linkMsgCost int64,
	recoveryPeriod int64,
	adjustPricePeriod int64,
	baseCreditPrice sdk.Dec,
	desirableBandwidth int64,
	maxBlockBandwidth uint64,
	txCost int64,
	nonLinkMsgCost int64) Params {

	return Params{
		LinkMsgCost:        linkMsgCost,
		RecoveryPeriod:     recoveryPeriod,
		AdjustPricePeriod:  adjustPricePeriod,
		BaseCreditPrice:    baseCreditPrice,
		DesirableBandwidth: desirableBandwidth,
		MaxBlockBandwidth:  maxBlockBandwidth,
		TxCost:             txCost,
		NonLinkMsgCost:     nonLinkMsgCost,
	}
}

// NewDefaultParams returns a default set of parameters.
func NewDefaultParams() Params {
	return Params{
		LinkMsgCost:        int64(100),
		RecoveryPeriod:     int64(18000),
		AdjustPricePeriod:  int64(10),
		BaseCreditPrice:    sdk.NewDec(1),
		DesirableBandwidth: int64(200000000),
		MaxBlockBandwidth:  uint64(200000000*10/18000),
		TxCost:             int64(300),
		NonLinkMsgCost:     int64(500),
	}
}
