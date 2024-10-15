package types

import paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

var (
	KeyCalculationPeriod = []byte("CalculationPeriod")
	KeyDampingFactor     = []byte("DampingFactor")
	KeyTolerance         = []byte("Tolerance")
)

// Deprecated: Type declaration for parameters
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyCalculationPeriod, &p.CalculationPeriod, validateCalculationPeriod),
		paramstypes.NewParamSetPair(KeyDampingFactor, &p.DampingFactor, validateDampingFactor),
		paramstypes.NewParamSetPair(KeyTolerance, &p.Tolerance, validateTolerance),
	}
}
