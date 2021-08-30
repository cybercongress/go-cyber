package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	ModuleName 					= "energy"
	StoreKey 					= ModuleName
  	RouterKey 					= ModuleName
	QuerierRoute 				= ModuleName

	EnergyPoolName  			= "energy_grid"
)

var (
	RouteKey                     = []byte{0x00}
	RoutedEnergyByDestinationKey = []byte{0x01}
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

