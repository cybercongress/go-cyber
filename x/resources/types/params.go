package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ctypes "github.com/cybercongress/go-cyber/types"
)

const (
	DefaultParamspace = ModuleName
	DefaultMaxSlots  = uint32(8)
	DefaultBaseVestingTime = 3600
)

var (
	KeyMaxSlots   		   = []byte("MaxSlots")
	KeyBaseVestingTime     = []byte("BaseVestingTime")
	KeyBaseVestingResource = []byte("BaseVestingResource")
)

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func DefaultParams() Params {
	return Params{
		MaxSlots: DefaultMaxSlots,
		BaseVestingTime: DefaultBaseVestingTime,
		BaseVestingResource: ctypes.NewCybCoin(ctypes.Mega*10),
	}
}

func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyMaxSlots, &p.MaxSlots, validateMaxSlots),
		paramstypes.NewParamSetPair(KeyBaseVestingTime, &p.BaseVestingTime, validateBaseVestingTime),
		paramstypes.NewParamSetPair(KeyBaseVestingResource, &p.BaseVestingResource, validateBaseVestingResource),
	}
}

func (p Params) Validate() error {
	if err := validateMaxSlots(p.MaxSlots); err != nil {
		return err
	}
	if err := validateBaseVestingTime(p.BaseVestingTime); err != nil {
		return err
	}
	if err := validateBaseVestingResource(p.BaseVestingResource); err != nil {
		return err
	}

	return nil
}

func validateBaseVestingTime(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 100 {
		return fmt.Errorf("base vesting time must be more than 100: %d", v)
	}

	return nil
}

func validateBaseVestingResource(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsLT(ctypes.NewCybCoin(ctypes.Mega*10)) {
		return fmt.Errorf("base vesting resource must be more than 10M: %d", v)
	}

	return nil
}

func validateMaxSlots(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max entries must be positive: %d", v)
	}

	if v > 16 {
		return fmt.Errorf("max entries must be less than 16: %d", v)
	}

	return nil
}
