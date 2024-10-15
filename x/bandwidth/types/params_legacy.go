package types

import paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

var (
	KeyRecoveryPeriod    = []byte("RecoveryPeriod")
	KeyAdjustPricePeriod = []byte("AdjustPricePeriod")
	KeyBasePrice         = []byte("BasePrice")
	KeyBaseLoad          = []byte("BaseLoad")
	KeyMaxBlockBandwidth = []byte("MaxBlockBandwidth")
)

// Deprecated: Type declaration for parameters
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyRecoveryPeriod, &p.RecoveryPeriod, validateRecoveryPeriod),
		paramstypes.NewParamSetPair(KeyAdjustPricePeriod, &p.AdjustPricePeriod, validateAdjustPricePeriod),
		paramstypes.NewParamSetPair(KeyBasePrice, &p.BasePrice, validateBasePrice),
		paramstypes.NewParamSetPair(KeyBaseLoad, &p.BaseLoad, validateBaseLoad),
		paramstypes.NewParamSetPair(KeyMaxBlockBandwidth, &p.MaxBlockBandwidth, validateMaxBlockBandwidth),
	}
}
