package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	ctypes "github.com/cybercongress/go-cyber/types"
)

const (
	DefaultParamspace = ModuleName

	DefaultMaxRoutes   uint16 = 10

	DefaultEnergyDenom string = ctypes.CYB
)

var (
	KeyMaxRoutes   = []byte("MaxRoutes")
	KeyEnergyDenom = []byte("EnergyDenom")
)

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

type Params struct {
	MaxRoutes   uint16 `json:"max_routes" yaml:"max_routes"`
	EnergyDenom string `json:"energy_denom" yaml:"energy_denom"`
}

func NewParams(maxRoutes uint16, energyDenom string) Params {
	return Params{
		MaxRoutes:   maxRoutes,
		EnergyDenom: energyDenom,
	}
}

func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Max Routes:    %d
  Energy Denom:  %s`, p.MaxRoutes, p.EnergyDenom)
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyMaxRoutes, &p.MaxRoutes, validateMaxRoutes),
		params.NewParamSetPair(KeyEnergyDenom, &p.EnergyDenom, validateEnergyDenom),
	}
}

func DefaultParams() Params {
	return NewParams(DefaultMaxRoutes, DefaultEnergyDenom)
}

func validateMaxRoutes(i interface{}) error {
	v, ok := i.(uint16)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max entries must be positive: %d", v)
	}

	return nil
}

func validateEnergyDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("bond denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}
