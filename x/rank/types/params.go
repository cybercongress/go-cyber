package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyCalculationPeriod = []byte("CalculationPeriod")
	KeyDampingFactor     = []byte("DampingFactor")
	KeyTolerance		 = []byte("Tolerance")
)

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}


// NewDefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		CalculationPeriod: int64(5),
		DampingFactor:	   sdk.NewDecWithPrec(85, 2),
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

	if v <= int64(2) {
		return fmt.Errorf("calculation period too low: %d", v)
	}

	return nil
}

func validateDampingFactor(i interface{}) error {
	v, ok := i.(sdk.Dec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.ZeroDec()) {
		return fmt.Errorf("damping factor should be positive: %s", v)
	}

	if v.GTE(sdk.OneDec()) {
		return fmt.Errorf("damping factor should be < 1, is: %s", v)
	}

	return nil
}

func validateTolerance(i interface{}) error {
	v, ok := i.(sdk.Dec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.NewDecWithPrec(1, 3)) {
		return fmt.Errorf("tolerance too low: %s", v)
	}

	if v.LT(sdk.NewDecWithPrec(1, 5)) {
		return fmt.Errorf("tolerance too big: %s", v)
	}

	return nil
}
