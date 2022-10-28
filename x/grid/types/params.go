package types

import (
	"fmt"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultParamspace = ModuleName
	DefaultMaxRoutes  = uint32(8)
)

var (
	KeyMaxRoutes = []byte("MaxRoutes")
)

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func DefaultParams() Params {
	return Params{
		MaxRoutes: DefaultMaxRoutes,
	}
}

func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyMaxRoutes, &p.MaxRoutes, validateMaxRoutes),
	}
}

func (p Params) Validate() error {
	if err := validateMaxRoutes(p.MaxRoutes); err != nil {
		return err
	}

	return nil
}

func validateMaxRoutes(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max routes must be positive: %d", v)
	}

	if v > 16 {
		return fmt.Errorf("max routes must be less or equal than 16: %d", v)
	}

	return nil
}
