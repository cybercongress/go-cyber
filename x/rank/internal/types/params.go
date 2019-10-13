package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"strconv"
	"strings"
)

// Default parameter values
const (
	DefaultCalculationPeriod int64  = 10
	DefaultDampingFactor   	 string = "0.85"
	DefaultTolerance		 string = "0.001"
)

// Parameter keys
var (
	KeyCalculationPeriod = []byte("CalculationPeriod")
	KeyDampingFactor     = []byte("DampingFactor")
	KeyTolerance		 = []byte("Tolerance")
)

// Params defines the parameters for the rank module.
type Params struct {
	CalculationPeriod int64
	DampingFactor 	  string
	Tolerance		  string
}

// NewParams creates a new Params object
func NewParams(
	calculationPeriod int64,
	dampingFactor string,
	tolerance string) Params {

	return Params{
		CalculationPeriod: calculationPeriod,
		DampingFactor:     dampingFactor,
		Tolerance:		   tolerance,
	}
}

// NewDefaultParams returns a default set of parameters.
func NewDefaultParams() Params {
	return Params{
		CalculationPeriod: DefaultCalculationPeriod,
		DampingFactor:	   DefaultDampingFactor,
		Tolerance:         DefaultTolerance,
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
	if p.CalculationPeriod < 1 {
		return fmt.Errorf("invalid calculation period: %d less then 1", p.CalculationPeriod)
	}
	if strconv.ParseFloat(p.DampingFactor) > float64(1.0) {

	}
	return nil
}
