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
	DefaultHalvingPeriodVolt      = uint32(512)
	DefaultHalvingPeriodAmpere    = uint32(512)
	DefaultInvestmintPeriodVolt   = uint32(1024)
	DefaultInvestmintPeriodAmpere = uint32(1024)
	DefaultMinInvestmintPeriodSec = uint32(300)
)

var (
	KeyMaxSlots   		          = []byte("MaxSlots")
	KeyBaseHalvingPeriodVolt      = []byte("HalvingPeriodVoltPeriod")
	KeyBaseHalvingPeriodAmpere    = []byte("HalvingPeriodAmperePeriod")
	KeyBaseInvestmintPeriodVolt   = []byte("BaseInvestmintPeriodVolt")
	KeyBaseInvestmintPeriodAmpere = []byte("BaseInvestmintPeriodAmpere")
	KeyBaseInvestmintAmountVolt   = []byte("BaseInvestmintAmountVolt")
	KeyBaseInvestmintAmountAmpere = []byte("BaseInvestmintAmountAmpere")
	KeyMinInvestmintPeriodSec     = []byte("MinInvestmintPeriod")
)

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})

}
func DefaultParams() Params {
	return Params{
		MaxSlots: 		 			DefaultMaxSlots,
		HalvingPeriodVoltBlocks:    DefaultHalvingPeriodVolt,
		HalvingPeriodAmpereBlocks:  DefaultHalvingPeriodAmpere,
		BaseInvestmintPeriodVolt:   DefaultInvestmintPeriodVolt,
		BaseInvestmintPeriodAmpere: DefaultInvestmintPeriodAmpere,
		BaseInvestmintAmountVolt:   ctypes.NewCybCoin(ctypes.Mega*10),
		BaseInvestmintAmountAmpere: ctypes.NewCybCoin(ctypes.Mega*10),
		MinInvestmintPeriod:        DefaultMinInvestmintPeriodSec,
	}
}

func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyMaxSlots, &p.MaxSlots, validateMaxSlots),
		paramstypes.NewParamSetPair(KeyBaseHalvingPeriodVolt, &p.HalvingPeriodVoltBlocks, validateBaseHalvingPeriodVolt),
		paramstypes.NewParamSetPair(KeyBaseHalvingPeriodAmpere, &p.HalvingPeriodAmpereBlocks, validateBaseHalvingPeriodAmpere),
		paramstypes.NewParamSetPair(KeyBaseInvestmintPeriodVolt, &p.BaseInvestmintPeriodVolt, validateBaseInvestmintPeriodVolt),
		paramstypes.NewParamSetPair(KeyBaseInvestmintPeriodAmpere, &p.BaseInvestmintPeriodAmpere, validateBaseInvestmintPeriodAmpere),
		paramstypes.NewParamSetPair(KeyBaseInvestmintAmountVolt, &p.BaseInvestmintAmountVolt, validateBaseInvestmintAmountVolt),
		paramstypes.NewParamSetPair(KeyBaseInvestmintAmountAmpere, &p.BaseInvestmintAmountAmpere, validateBaseInvestmintAmountAmpere),
		paramstypes.NewParamSetPair(KeyMinInvestmintPeriodSec, &p.MinInvestmintPeriod, validateMinInvestmintPeriodSec),
	}
}

func (p Params) Validate() error {
	if err := validateMaxSlots(p.MaxSlots); err != nil {
		return err
	}
	if err := validateBaseHalvingPeriodVolt(p.HalvingPeriodVoltBlocks); err != nil {
		return err
	}
	if err := validateBaseHalvingPeriodAmpere(p.HalvingPeriodAmpereBlocks); err != nil {
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
	if err := validateMinInvestmintPeriodSec(p.MinInvestmintPeriod); err != nil {
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

func validateBaseHalvingPeriodVolt(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// TODO set production value
	if v < 100 {
		return fmt.Errorf("base halving period for Volt must be more than 100 blocks: %d", v)
	}

	return nil
}

func validateBaseHalvingPeriodAmpere(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// TODO set production value
	if v < 100 {
		return fmt.Errorf("base halving period for Ampere must be more than 100 blocks: %d", v)
	}

	return nil
}

func validateBaseInvestmintPeriodVolt(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// TODO set production value
	if v < 100 {
		return fmt.Errorf("base investmint period for Volt must be more than 100 blocks: %d", v)
	}

	return nil
}

func validateBaseInvestmintPeriodAmpere(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// TODO set production value
	if v < 100 {
		return fmt.Errorf("base investmint period for Ampere must be more than 100 blocks: %d", v)
	}

	return nil
}

func validateBaseInvestmintAmountVolt(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsLT(ctypes.NewCybCoin(ctypes.Mega*10)) {
		return fmt.Errorf("base investmint amount for Volt must be more than 10000000: %d", v)
	}

	return nil
}

func validateBaseInvestmintAmountAmpere(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsLT(ctypes.NewCybCoin(ctypes.Mega*10)) {
		return fmt.Errorf("base investmint amount for Ampere must be more than 10000000: %d", v)
	}

	return nil
}

func validateMinInvestmintPeriodSec(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// TODO set production value
	if v < 100 {
		return fmt.Errorf("min investmint period must be more than 100 seconds: %d", v)
	}

	return nil
}
