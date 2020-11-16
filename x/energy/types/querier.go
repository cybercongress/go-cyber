package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryParams 			     = "params"
	QuerySourceRoutes   		 = "source_routes"
	QueryDestinationRoutes 	     = "destination_routes"
	QuerySourceRoutedEnergy      = "source_routed_energy"
	QueryDestinationRoutedEnergy = "destination_routed_energy"
	QueryRoute      			 = "route"
	QueryRoutes      			 = "routes"
)

type QueryDestinationParams struct {
	 Destination sdk.AccAddress
}

func NewQueryDestinationParams(dst sdk.AccAddress) QueryDestinationParams {
	return QueryDestinationParams{
		Destination: dst,
	}
}

type QuerySourceParams struct {
	Source sdk.AccAddress
}

func NewQuerySourceParams(src sdk.AccAddress) QuerySourceParams {
	return QuerySourceParams{
		Source: src,
	}
}

type QueryRouteParams struct {
	Source 		sdk.AccAddress
	Destination sdk.AccAddress
}

func NewQueryRouteParams(src, dst sdk.AccAddress) QueryRouteParams {
	return QueryRouteParams{
		Source: 	 src,
		Destination: dst,
	}
}