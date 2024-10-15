package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	ModuleName = "grid"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	GridPoolName = "energy_grid"
)

var (
	RouteKey                     = []byte{0x00}
	RoutedEnergyByDestinationKey = []byte{0x01}
	ParamsKey                    = []byte{0x02}
)

func GetRoutedEnergyByDestinationKey(dst sdk.AccAddress) []byte {
	return append(RoutedEnergyByDestinationKey, dst.Bytes()...)
}

func GetRouteKey(src sdk.AccAddress, dst sdk.AccAddress) []byte {
	return append(GetRoutesKey(src), dst.Bytes()...)
}

func GetRoutesKey(src sdk.AccAddress) []byte {
	return append(RouteKey, src.Bytes()...)
}
