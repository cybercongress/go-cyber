package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ctypes "github.com/cybercongress/go-cyber/types"
)

const (
	DefaultParamspace             = ModuleName
	DefaultMaxSlots               = uint32(8)
	DefaultHalvingPeriodVolt      = uint32(9000000)
	DefaultHalvingPeriodAmpere    = uint32(9000000)
	DefaultInvestmintPeriodVolt   = uint32(2592000)
	DefaultInvestmintPeriodAmpere = uint32(2592000)
	DefaultMinInvestmintPeriod    = uint32(86400)
)

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

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func DefaultParams() Params {
	return Params{
		MaxSlots:                   DefaultMaxSlots,
		HalvingPeriodVoltBlocks:    DefaultHalvingPeriodVolt,
		HalvingPeriodAmpereBlocks:  DefaultHalvingPeriodAmpere,
		BaseInvestmintPeriodVolt:   DefaultInvestmintPeriodVolt,
		BaseInvestmintPeriodAmpere: DefaultInvestmintPeriodAmpere,
		BaseInvestmintAmountVolt:   ctypes.NewSCybCoin(ctypes.Mega * 1000),
		BaseInvestmintAmountAmpere: ctypes.NewSCybCoin(ctypes.Mega * 100),
		MinInvestmintPeriod:        DefaultMinInvestmintPeriod,
	}
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

func (p Params) Validate() error {
	if err := validateMaxSlots(p.MaxSlots); err != nil {
		return err
	}
	if err := validateHalvingPeriodVoltBlocks(p.HalvingPeriodVoltBlocks); err != nil {
		return err
	}
	if err := validateHalvingPeriodAmpereBlocks(p.HalvingPeriodAmpereBlocks); err != nil {
		return err
	}
	if err := validateBaseInvestmintPeriodVolt(p.BaseInvestmintPeriodVolt); err != nil {
		return err
	}
	if err := validateBaseInvestmintPeriodAmpere(p.BaseInvestmintPeriodAmpere); err != nil {
		return err
	}
	if err := validateBaseInvestmintAmountVolt(p.BaseInvestmintAmountVolt); err != nil {
		return err
	}
	if err := validateBaseInvestmintAmountAmpere(p.BaseInvestmintAmountAmpere); err != nil {
		return err
	}
	if err := validateMinInvestmintPeriod(p.MinInvestmintPeriod); err != nil {
		return err
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
		return fmt.Errorf("max resources slots must be less or equal to 16: %d", v)
	}

	return nil
}

func validateHalvingPeriodVoltBlocks(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 6000000 {
		return fmt.Errorf("base halving period for Volt must be more than 6000000 blocks: %d", v)
	}

	return nil
}

func validateHalvingPeriodAmpereBlocks(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 6000000 {
		return fmt.Errorf("base halving period for Ampere must be more than 6000000 blocks: %d", v)
	}

	return nil
}

func validateBaseInvestmintPeriodVolt(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 604800 {
		return fmt.Errorf("base investmint period for Volt must be more than 604800 seconds: %d", v)
	}

	return nil
}

func validateBaseInvestmintPeriodAmpere(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 604800 {
		return fmt.Errorf("base investmint period for Ampere must be more than 604800 seconds: %d", v)
	}

	return nil
}

func validateBaseInvestmintAmountVolt(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsLT(ctypes.NewSCybCoin(ctypes.Mega * 10)) {
		return fmt.Errorf("base investmint amount for Volt must be more than 10000000: %d", v)
	}

	return nil
}

func validateBaseInvestmintAmountAmpere(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsLT(ctypes.NewSCybCoin(ctypes.Mega * 10)) {
		return fmt.Errorf("base investmint amount for Ampere must be more than 10000000: %d", v)
	}

	return nil
}

func validateMinInvestmintPeriod(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 86400 {
		return fmt.Errorf("min investmint period must be more than 86400 seconds: %d", v)
	}

	return nil
}
