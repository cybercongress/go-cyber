package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyCalculationPeriod = []byte("CalculationPeriod")
	KeyDampingFactor     = []byte("DampingFactor")
	KeyTolerance         = []byte("Tolerance")
)

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func DefaultParams() Params {
	return Params{
		CalculationPeriod: int64(5),
		DampingFactor:     sdk.NewDecWithPrec(85, 2),
		Tolerance:         sdk.NewDecWithPrec(1, 3),
	}
}

func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyCalculationPeriod, &p.CalculationPeriod, validateCalculationPeriod),
		paramstypes.NewParamSetPair(KeyDampingFactor, &p.DampingFactor, validateDampingFactor),
		paramstypes.NewParamSetPair(KeyTolerance, &p.Tolerance, validateTolerance),
	}
}

func (p Params) ValidateBasic() error {
	if err := validateCalculationPeriod(p.CalculationPeriod); err != nil {
		return err
	}
	if err := validateDampingFactor(p.DampingFactor); err != nil {
		return err
	}
	if err := validateTolerance(p.Tolerance); err != nil {
		return err
	}

	return nil
}

func validateCalculationPeriod(i interface{}) error {
	v, ok := i.(int64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < int64(5) {
		return fmt.Errorf("calculation period should be equal or more than 5 blocks: %d", v)
	}

	return nil
}

func validateDampingFactor(i interface{}) error {
	v, ok := i.(sdk.Dec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LTE(sdk.NewDecWithPrec(7, 1)) {
		return fmt.Errorf("damping factor should be equal or more than 0.7: %s", v)
	}

	if v.GTE(sdk.NewDecWithPrec(9, 1)) {
		return fmt.Errorf("damping factor should be equal or less than 0.9: %s", v)
	}

	return nil
}

func validateTolerance(i interface{}) error {
	v, ok := i.(sdk.Dec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.NewDecWithPrec(1, 3)) {
		return fmt.Errorf("tolerance is too low: %s", v)
	}

	if v.LT(sdk.NewDecWithPrec(1, 5)) {
		return fmt.Errorf("tolerance is too big: %s", v)
	}

	return nil
}
