package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultParamspace = ModuleName
)

var (
	KeyTxCost             = []byte("TxCost")
	KeyLinkCost 		  = []byte("LinkCost")
	KeyRecoveryPeriod     = []byte("RecoveryPeriod")
	KeyAdjustPricePeriod  = []byte("AdjustPricePeriod")
	KeyBaseCreditPrice    = []byte("BaseCreditPrice")
	KeyDesirableBandwidth = []byte("DesirableBandwidth")
	KeyMaxBlockBandwidth  = []byte("MaxBlockBandwidth")
)

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func DefaultParams() Params {
	return Params{
		TxCost:             uint64(300),
		LinkCost:           uint64(100),
		RecoveryPeriod:     uint64(100),
		AdjustPricePeriod:  uint64(5),
		BaseCreditPrice:    sdk.NewDecWithPrec(5,1),
		DesirableBandwidth: uint64(500000),
		MaxBlockBandwidth:  uint64(10000),
	}
}

func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyTxCost, &p.TxCost, validateTxCost),
		paramstypes.NewParamSetPair(KeyLinkCost, &p.LinkCost, validateLinkCost),
		paramstypes.NewParamSetPair(KeyRecoveryPeriod, &p.RecoveryPeriod, validateRecoveryPeriod),
		paramstypes.NewParamSetPair(KeyAdjustPricePeriod, &p.AdjustPricePeriod, validateAdjustPricePeriod),
		paramstypes.NewParamSetPair(KeyBaseCreditPrice, &p.BaseCreditPrice, validateBaseCreditPrice),
		paramstypes.NewParamSetPair(KeyDesirableBandwidth, &p.DesirableBandwidth, validateDesirableBandwidth),
		paramstypes.NewParamSetPair(KeyMaxBlockBandwidth, &p.MaxBlockBandwidth, validateMaxBlockBandwidth),
	}
}

func (p Params) Validate() error {
	if err := validateTxCost(p.TxCost); err != nil {
		return err
	}
	if err := validateLinkCost(p.LinkCost); err != nil {
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
	if err := validateMaxBlockBandwidth(p.MaxBlockBandwidth); err != nil {
		return err
	}

	return nil
}

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

func validateLinkCost(i interface{}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= uint64(10) {
		return fmt.Errorf("link msg cost too low: %d", v)
	}

	return nil
}

func validateRecoveryPeriod(i interface{}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= uint64(10) {
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
	//v, ok := i.(sdk.Dec)
	//
	//if !ok {
	//	return fmt.Errorf("invalid parameter type: %T", i)
	//}

	//if v.LT(sdk.OneDec()) {
	//	return fmt.Errorf("base credit price too low: %s", v)
	//}

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

func validateMaxBlockBandwidth(i interface{}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= uint64(100) {
		return fmt.Errorf("max block bandwidth too low: %d", v)
	}

	return nil
}