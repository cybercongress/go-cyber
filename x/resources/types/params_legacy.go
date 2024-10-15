package types

import paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

var (
	KeyMaxSlots                   = []byte("MaxSlots")
	KeyHalvingPeriodVoltBlocks    = []byte("HalvingPeriodVoltBlocks")
	KeyHalvingPeriodAmpereBlocks  = []byte("HalvingPeriodAmpereBlocks")
	KeyBaseInvestmintPeriodVolt   = []byte("BaseInvestmintPeriodVolt")
	KeyBaseInvestmintPeriodAmpere = []byte("BaseInvestmintPeriodAmpere")
	KeyBaseInvestmintAmountVolt   = []byte("BaseInvestmintAmountVolt")
	KeyBaseInvestmintAmountAmpere = []byte("BaseInvestmintAmountAmpere")
	KeyMinInvestmintPeriod        = []byte("MinInvestmintPeriod")
)

// Deprecated: Type declaration for parameters
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyMaxSlots, &p.MaxSlots, validateMaxSlots),
		paramstypes.NewParamSetPair(KeyHalvingPeriodVoltBlocks, &p.HalvingPeriodVoltBlocks, validateHalvingPeriodVoltBlocks),
		paramstypes.NewParamSetPair(KeyHalvingPeriodAmpereBlocks, &p.HalvingPeriodAmpereBlocks, validateHalvingPeriodAmpereBlocks),
		paramstypes.NewParamSetPair(KeyBaseInvestmintPeriodVolt, &p.BaseInvestmintPeriodVolt, validateBaseInvestmintPeriodVolt),
		paramstypes.NewParamSetPair(KeyBaseInvestmintPeriodAmpere, &p.BaseInvestmintPeriodAmpere, validateBaseInvestmintPeriodAmpere),
		paramstypes.NewParamSetPair(KeyBaseInvestmintAmountVolt, &p.BaseInvestmintAmountVolt, validateBaseInvestmintAmountVolt),
		paramstypes.NewParamSetPair(KeyBaseInvestmintAmountAmpere, &p.BaseInvestmintAmountAmpere, validateBaseInvestmintAmountAmpere),
		paramstypes.NewParamSetPair(KeyMinInvestmintPeriod, &p.MinInvestmintPeriod, validateMinInvestmintPeriod),
	}
}
