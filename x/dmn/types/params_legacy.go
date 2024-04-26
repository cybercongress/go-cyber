package types

import paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

var (
	KeyMaxSlots = []byte("MaxSlots")
	KeyMaxGas   = []byte("MaxGas")
	KeyFeeTTL   = []byte("FeeTTL")
)

// Deprecated: Type declaration for parameters
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyMaxSlots, &p.MaxSlots, validateMaxSlots),
		paramstypes.NewParamSetPair(KeyMaxGas, &p.MaxGas, validateMaxGas),
		paramstypes.NewParamSetPair(KeyFeeTTL, &p.FeeTtl, validateFeeTTL),
	}
}
