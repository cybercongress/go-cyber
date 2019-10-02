package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"strings"
)

// Default parameter values
const (
	DefaultCalculationPeriod int64 = 10
)

// Parameter keys
var (
	KeyCalculationPeriod = []byte("CalculationPeriod")
)

// Params defines the parameters for the rank module.
type Params struct {
	CalculationPeriod int64
}

// NewParams creates a new Params object
func NewParams(calculationPeriod int64) Params {

	return Params{
		CalculationPeriod: calculationPeriod,
	}
}

// NewDefaultParams returns a default set of parameters.
func NewDefaultParams() Params {
	return Params{
		CalculationPeriod: DefaultCalculationPeriod,
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
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("CalculationPeriod: %d\n", p.CalculationPeriod))

	return sb.String()
}

func (p Params) Validate() error {
	if p.CalculationPeriod < 1 {
		return fmt.Errorf("invalid calculation period: %d less then 1", p.CalculationPeriod)
	}
	return nil
}
