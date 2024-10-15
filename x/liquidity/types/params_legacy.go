package types

import paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

// Deprecated: Type declaration for parameters
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements paramstypes.ParamSet.
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyPoolTypes, &p.PoolTypes, validatePoolTypes),
		paramstypes.NewParamSetPair(KeyMinInitDepositAmount, &p.MinInitDepositAmount, validateMinInitDepositAmount),
		paramstypes.NewParamSetPair(KeyInitPoolCoinMintAmount, &p.InitPoolCoinMintAmount, validateInitPoolCoinMintAmount),
		paramstypes.NewParamSetPair(KeyMaxReserveCoinAmount, &p.MaxReserveCoinAmount, validateMaxReserveCoinAmount),
		paramstypes.NewParamSetPair(KeyPoolCreationFee, &p.PoolCreationFee, validatePoolCreationFee),
		paramstypes.NewParamSetPair(KeySwapFeeRate, &p.SwapFeeRate, validateSwapFeeRate),
		paramstypes.NewParamSetPair(KeyWithdrawFeeRate, &p.WithdrawFeeRate, validateWithdrawFeeRate),
		paramstypes.NewParamSetPair(KeyMaxOrderAmountRatio, &p.MaxOrderAmountRatio, validateMaxOrderAmountRatio),
		paramstypes.NewParamSetPair(KeyUnitBatchHeight, &p.UnitBatchHeight, validateUnitBatchHeight),
		paramstypes.NewParamSetPair(KeyCircuitBreakerEnabled, &p.CircuitBreakerEnabled, validateCircuitBreakerEnabled),
	}
}
