package types

import paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

var KeyMaxRoutes = []byte("MaxRoutes")

// Deprecated: Type declaration for parameters
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyMaxRoutes, &p.MaxRoutes, validateMaxRoutes),
	}
}
