package types

import (
	"fmt"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultParamspace = ModuleName

	DefaultMaxSlots   uint32 = 4
	DefaultMaxGas     uint32 = 2000000
	DefaultFeeTTL     uint32 = 50
)

var (
	KeyMaxSlots   = []byte("MaxSlots")
	KeyMaxGas     = []byte("MaxGas")
	KeyFeeTTL     = []byte("FeeTTL")
)

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func DefaultParams() Params {
	return Params{
		DefaultMaxSlots,
		DefaultMaxGas,
		DefaultFeeTTL,
	}
}

func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyMaxSlots, &p.MaxSlots, validateMaxSlots),
		paramstypes.NewParamSetPair(KeyMaxGas, &p.MaxGas, validateMaxGas),
		paramstypes.NewParamSetPair(KeyFeeTTL, &p.FeeTtl, validateFeeTTL),
	}
}


func validateMaxSlots(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 4 {
		return fmt.Errorf("max slots must be equal or more than 4: %d", v)
	}

	return nil
}

func validateMaxGas(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 2000000 {
		return fmt.Errorf("max gas must be equal or more than 2000000: %d", v)
	}

	return nil
}

func validateFeeTTL(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("fee ttl must be positive: %d", v)
	}

	return nil
}
