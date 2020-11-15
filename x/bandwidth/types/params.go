package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultParamspace = ModuleName
)

var (
	KeyTxCost             = []byte("TxCost")
	KeyLinkMsgCost 		  = []byte("LinkMsgCost")
	KeyNonLinkMsgCost     = []byte("NonLinkMsgCost")
	KeyRecoveryPeriod     = []byte("RecoveryPeriod")
	KeyAdjustPricePeriod  = []byte("AdjustPricePeriod")
	KeyBaseCreditPrice    = []byte("BaseCreditPrice")
	KeyDesirableBandwidth = []byte("DesirableBandwidth")
	KeyMaxBlockBandwidth  = []byte("MaxBlockBandwidth")
)

// TODO move to proto
type Params struct {
	TxCost             uint64   `json:"tx_cost" yaml:"tx_cost"`
	LinkMsgCost        uint64   `json:"link_msg_cost" yaml:"link_msg_cost"`
	NonLinkMsgCost     uint64   `json:"non_link_msg_cost" yaml:"non_link_msg_cost"`
	RecoveryPeriod     uint64   `json:"recovery_period" yaml:"recovery_period"`
	AdjustPricePeriod  uint64   `json:"adjust_price_period" yaml:"adjust_price_period"`
	BaseCreditPrice    sdk.Dec  `json:"base_credit_price" yaml:"base_credit_price"`
	DesirableBandwidth uint64   `json:"desirable_bandwidth" yaml:"desirable_bandwidth"`
	MaxBlockBandwidth  uint64   `json:"max_block_bandwidth" yaml:"max_block_bandwidth"`
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func DefaultParams() Params {
	return Params{
		TxCost:             uint64(3000),
		LinkMsgCost:        uint64(1000),
		NonLinkMsgCost:     uint64(5000),
		RecoveryPeriod:     uint64(700),
		AdjustPricePeriod:  uint64(5),
		BaseCreditPrice:    sdk.NewDec(50),
		DesirableBandwidth: uint64(200000),
		MaxBlockBandwidth:  uint64(10000),
	}
}

func (p Params) String() string {

	return fmt.Sprintf(`Bandwidth params:
  LinkMsgCost:        %d
  TxCost:			  %d
  NonLinkMsgCost:	  %d
  RecoveryPeriod:     %d
  AdjustPricePeriod:  %d
  BaseCreditPrice:    %d
  DesirableBandwidth: %d
  MaxBlockBandidth:   %d
`,
		p.LinkMsgCost, p.RecoveryPeriod, p.AdjustPricePeriod,
		p.BaseCreditPrice, p.DesirableBandwidth, p.MaxBlockBandwidth,
		p.TxCost, p.NonLinkMsgCost,
	)
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyTxCost, &p.TxCost, validateTxCost),
		paramtypes.NewParamSetPair(KeyLinkMsgCost, &p.LinkMsgCost, validateLinkMsgCost),
		paramtypes.NewParamSetPair(KeyNonLinkMsgCost, &p.NonLinkMsgCost, validateNonLinkMsgCost),
		paramtypes.NewParamSetPair(KeyRecoveryPeriod, &p.RecoveryPeriod, validateRecoveryPeriod),
		paramtypes.NewParamSetPair(KeyAdjustPricePeriod, &p.AdjustPricePeriod, validateAdjustPricePeriod),
		paramtypes.NewParamSetPair(KeyBaseCreditPrice, &p.BaseCreditPrice, validateBaseCreditPrice),
		paramtypes.NewParamSetPair(KeyDesirableBandwidth, &p.DesirableBandwidth, validateDesirableBandwidth),
		paramtypes.NewParamSetPair(KeyMaxBlockBandwidth, &p.MaxBlockBandwidth, validateMaxBlockBanwidth),
	}
}

func (p Params) ValidateBasic() error {
	if err := validateTxCost(p.TxCost); err != nil {
		return err
	}
	if err := validateLinkMsgCost(p.TxCost); err != nil {
		return err
	}
	if err := validateNonLinkMsgCost(p.NonLinkMsgCost); err != nil {
		return err
	}
	if err := validateRecoveryPeriod(p.RecoveryPeriod); err != nil {
		return err
	}
	if err := validateAdjustPricePeriod(p.AdjustPricePeriod); err != nil {
		return err
	}
	if err := validateBaseCreditPrice(p.BaseCreditPrice); err != nil {
		return err
	}
	if err := validateDesirableBandwidth(p.DesirableBandwidth); err != nil {
		return err
	}
	if err := validateMaxBlockBanwidth(p.MaxBlockBandwidth); err != nil {
		return err
	}

	return nil
}

// TODO improve validations for parameters
func validateTxCost(i interface{}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= uint64(10) {
		return fmt.Errorf("tx cost too low: %d", v)
	}

	return nil
}

func validateLinkMsgCost(i interface{}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= uint64(10) {
		return fmt.Errorf("link msg cost too low: %d", v)
	}

	return nil
}

func validateNonLinkMsgCost(i interface{}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= uint64(10) {
		return fmt.Errorf("non link msg too low: %d", v)
	}

	return nil
}

func validateRecoveryPeriod(i interface{}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= uint64(100) {
		return fmt.Errorf("recovery period too low: %d", v)
	}

	return nil
}

func validateAdjustPricePeriod(i interface{}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= uint64(2) {
		return fmt.Errorf("adjust price period too low: %d", v)
	}

	return nil
}

func validateBaseCreditPrice(i interface{}) error {
	v, ok := i.(sdk.Dec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.OneDec()) {
		return fmt.Errorf("base credit price too low: %s", v)
	}

	return nil
}

func validateDesirableBandwidth(i interface{}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= uint64(10000) {
		return fmt.Errorf("desirable bandwidth too low: %d", v)
	}

	return nil
}

func validateMaxBlockBanwidth(i interface{}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= uint64(100) {
		return fmt.Errorf("max block bandwidth too low: %d", v)
	}

	return nil
}