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


type QueryRoutesParams struct {
	Page, Limit   int
}

func NewQueryRoutesParams(page, limit int) QueryRoutesParams {
	return QueryRoutesParams{page, limit}
}

type QueryRouteParams struct {
	Source, Destination   sdk.AccAddress
}

func NewQueryRouteParams(source, destination sdk.AccAddress) QueryRouteParams {
	return QueryRouteParams{source, destination}
}

type QuerySourceParams struct {
	Source   sdk.AccAddress
}

func NewQuerySourceParams(source sdk.AccAddress) QuerySourceParams {
	return QuerySourceParams{source}
}

type QueryDestinationParams struct {
	Destination   sdk.AccAddress
}

func NewQueryDestinationParams(destination sdk.AccAddress) QueryDestinationParams {
	return QueryDestinationParams{destination}
}