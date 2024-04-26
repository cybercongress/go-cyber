package types

import (
	"fmt"
)

const (
	DefaultMaxSlots uint32 = 4
	DefaultMaxGas   uint32 = 2000000
	DefaultFeeTTL   uint32 = 50
)

func DefaultParams() Params {
	return Params{
		DefaultMaxSlots,
		DefaultMaxGas,
		DefaultFeeTTL,
	}
}

func NewParams(
	maxSlots uint32,
	maxGas uint32,
	feeTtl uint32,
) Params {
	return Params{
		maxSlots,
		maxGas,
		feeTtl,
	}
}

func (p Params) Validate() error {
	if err := validateMaxSlots(p.MaxSlots); err != nil {
		return err
	}
	if err := validateMaxGas(p.MaxGas); err != nil {
		return err
	}
	if err := validateFeeTTL(p.FeeTtl); err != nil {
		return err
	}

	return nil
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
