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
	KeyRecoveryPeriod     = []byte("RecoveryPeriod")
	KeyAdjustPricePeriod  = []byte("AdjustPricePeriod")
	KeyBasePrice          = []byte("BasePrice")
	KeyMaxBlockBandwidth  = []byte("MaxBlockBandwidth")
)

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func DefaultParams() Params {
	return Params{
		RecoveryPeriod:     uint64(100),
		AdjustPricePeriod:  uint64(5),
		BasePrice:          sdk.NewDecWithPrec(5,1),
		MaxBlockBandwidth:  uint64(10000),
	}
}

func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyRecoveryPeriod, &p.RecoveryPeriod, validateRecoveryPeriod),
		paramstypes.NewParamSetPair(KeyAdjustPricePeriod, &p.AdjustPricePeriod, validateAdjustPricePeriod),
		paramstypes.NewParamSetPair(KeyBasePrice, &p.BasePrice, validateBasePrice),
		paramstypes.NewParamSetPair(KeyMaxBlockBandwidth, &p.MaxBlockBandwidth, validateMaxBlockBandwidth),
	}
}

func (p Params) Validate() error {
	if err := validateRecoveryPeriod(p.RecoveryPeriod); err != nil {
		return err
	}
	if err := validateAdjustPricePeriod(p.AdjustPricePeriod); err != nil {
		return err
	}
	if err := validateBasePrice(p.BasePrice); err != nil {
		return err
	}
	if err := validateMaxBlockBandwidth(p.MaxBlockBandwidth); err != nil {
		return err
	}

	return nil
}

func validateRecoveryPeriod(i interface{}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= uint64(50) {
		return fmt.Errorf("recovery period too low: %d", v)
	}

	return nil
}

func validateAdjustPricePeriod(i interface{}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < uint64(5) {
		return fmt.Errorf("adjust price period too low: %d", v)
	}

	return nil
}

func validateBasePrice(i interface{}) error {
	v, ok := i.(sdk.Dec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsPositive() {
		return fmt.Errorf("base price is not positive: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("base price is more than one: %s", v)
	}

	return nil
}

func validateMaxBlockBandwidth(i interface{}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= uint64(1000) {
		return fmt.Errorf("max block bandwidth too low: %d", v)
	}

	return nil
}