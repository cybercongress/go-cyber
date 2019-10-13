package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// Parameter keys
var (
	KeyCalculationPeriod = []byte("CalculationPeriod")
	KeyDampingFactor     = []byte("DampingFactor")
	KeyTolerance		 = []byte("Tolerance")
)

// Params defines the parameters for the rank module.
type Params struct {
	CalculationPeriod int64   `json:"calculation_period" yaml:"calculation_period"`
	DampingFactor 	  sdk.Dec `json:"damping_factor" yaml:"damping_factor"`
	Tolerance		  sdk.Dec `json:"tolerance" yaml:"tolerance"`
}

// NewParams creates a new Params object
func NewParams(calculationPeriod int64, dampingFactor sdk.Dec,
	tolerance sdk.Dec) Params {

	return Params{
		CalculationPeriod: calculationPeriod,
		DampingFactor:     dampingFactor,
		Tolerance:		   tolerance,
	}
}

// NewDefaultParams returns a default set of parameters.
func NewDefaultParams() Params {
	return Params{
		CalculationPeriod: int64(10),
		DampingFactor:	   sdk.NewDecWithPrec(85, 2),
		Tolerance:         sdk.NewDecWithPrec(1, 3),
	}
}

// ParamKeyTable for rank module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of rank module's parameters.
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{KeyCalculationPeriod, &p.CalculationPeriod},
		{KeyDampingFactor, &p.DampingFactor},
		{KeyTolerance, &p.Tolerance},
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("CalculationPeriod: %d\n", p.CalculationPeriod))
	sb.WriteString(fmt.Sprintf("DampingFactor: %d\n", p.DampingFactor))
	sb.WriteString(fmt.Sprintf("Tolerance: %d\n", p.Tolerance))

	return sb.String()
}

func (p Params) Validate() error {
	if p.CalculationPeriod < 2 {
		return fmt.Errorf("invalid calculation period: %d less then 2", p.CalculationPeriod)
	}
	if p.DampingFactor.GTE(sdk.OneDec()) {
		return fmt.Errorf("damping factor parameter must be < 1, is %s", p.DampingFactor.String())
	}
	if p.DampingFactor.LT(sdk.ZeroDec()) {
		return fmt.Errorf("damping factor parameter should be positive, is %s", p.DampingFactor.String())
	}
	if p.Tolerance.GT(sdk.NewDecWithPrec(1, 3)) {
		return fmt.Errorf("tolerance parameter must be <= 0.001, is %s", p.DampingFactor.String())
	}
	if p.Tolerance.LT(sdk.NewDecWithPrec(1, 5)) {
		return fmt.Errorf("tolerance parameter must be >= 0.00001, is %s", p.DampingFactor.String())
	}
	return nil
}
