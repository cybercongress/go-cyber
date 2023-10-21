package wasm

import (
	"encoding/json"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/cybercongress/go-cyber/v2/x/grid/keeper"
	"github.com/cybercongress/go-cyber/v2/x/grid/types"
)

var _ MsgParserInterface = MsgParser{}

//--------------------------------------------------

type MsgParserInterface interface {
	Parse(contractAddr sdk.AccAddress, msg wasmvmtypes.CosmosMsg) ([]sdk.Msg, error)
	ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error)
}

type MsgParser struct{}

func NewWasmMsgParser() MsgParser {
	return MsgParser{}
}

func (MsgParser) Parse(_ sdk.AccAddress, _ wasmvmtypes.CosmosMsg) ([]sdk.Msg, error) {
	return nil, nil
}

type CosmosMsg struct {
	CreateEnergyRoute   *types.MsgCreateRoute   `json:"create_energy_route,omitempty"`
	EditEnergyRoute     *types.MsgEditRoute     `json:"edit_energy_route,omitempty"`
	EditEnergyRouteName *types.MsgEditRouteName `json:"edit_energy_route_name,omitempty"`
	DeleteEnergyRoute   *types.MsgDeleteRoute   `json:"delete_energy_route,omitempty"`
}

func (MsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var sdkMsg CosmosMsg
	err := json.Unmarshal(data, &sdkMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to parse link custom msg")
	}

	switch {
	case sdkMsg.CreateEnergyRoute != nil:
		return []sdk.Msg{sdkMsg.CreateEnergyRoute}, sdkMsg.CreateEnergyRoute.ValidateBasic()
	case sdkMsg.EditEnergyRoute != nil:
		return []sdk.Msg{sdkMsg.EditEnergyRoute}, sdkMsg.EditEnergyRoute.ValidateBasic()
	case sdkMsg.EditEnergyRouteName != nil:
		return []sdk.Msg{sdkMsg.EditEnergyRouteName}, sdkMsg.EditEnergyRouteName.ValidateBasic()
	case sdkMsg.DeleteEnergyRoute != nil:
		return []sdk.Msg{sdkMsg.DeleteEnergyRoute}, sdkMsg.DeleteEnergyRoute.ValidateBasic()
	default:
		return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Energy")
	}
}

//--------------------------------------------------

type QuerierInterface interface {
	Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

type Querier struct {
	keeper.Keeper
}

func NewWasmQuerier(keeper keeper.Keeper) Querier {
	return Querier{keeper}
}

func (Querier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) { return nil, nil }

type CosmosQuery struct {
	SourceRoutes            *QuerySourceParams      `json:"source_routes,omitempty"`
	SourceRoutedEnergy      *QuerySourceParams      `json:"source_routed_energy,omitempty"`
	DestinationRoutedEnergy *QueryDestinationParams `json:"destination_routed_energy,omitempty"`
	Route                   *QueryRouteParams       `json:"route,omitempty"`
}

type Route struct {
	Source      string            `json:"source"`
	Destination string            `json:"destination"`
	Name        string            `json:"name"`
	Value       wasmvmtypes.Coins `json:"value"`
}

type Routes []Route

type QuerySourceParams struct {
	Source string `json:"source"`
}

type QueryDestinationParams struct {
	Destination string `json:"destination"`
}

type QueryRouteParams struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type RoutesResponse struct {
	Routes Routes `json:"routes"`
}

type RoutedEnergyResponse struct {
	Value wasmvmtypes.Coins `json:"value"`
}

type RouteResponse struct {
	Route Route `json:"route"`
}

func (querier Querier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	switch {
	case query.SourceRoutes != nil:
		source, _ := sdk.AccAddressFromBech32(query.SourceRoutes.Source)
		routes := querier.Keeper.GetSourceRoutes(ctx, source, 16)

		bz, err = json.Marshal(RoutesResponse{
			Routes: convertCyberRoutesToWasmRoutes(routes),
		})
	case query.SourceRoutedEnergy != nil:
		source, _ := sdk.AccAddressFromBech32(query.SourceRoutedEnergy.Source)
		value := querier.Keeper.GetRoutedFromEnergy(ctx, source)

		bz, err = json.Marshal(RoutedEnergyResponse{
			Value: wasmkeeper.ConvertSdkCoinsToWasmCoins(value),
		})
	case query.DestinationRoutedEnergy != nil:
		destination, _ := sdk.AccAddressFromBech32(query.DestinationRoutedEnergy.Destination)
		value := querier.Keeper.GetRoutedToEnergy(ctx, destination)

		bz, err = json.Marshal(RoutedEnergyResponse{
			Value: wasmkeeper.ConvertSdkCoinsToWasmCoins(value),
		})
	case query.Route != nil:
		source, _ := sdk.AccAddressFromBech32(query.Route.Source)
		destination, _ := sdk.AccAddressFromBech32(query.Route.Destination)
		route, found := querier.Keeper.GetRoute(ctx, source, destination)
		if !found {
			return nil, sdkerrors.ErrInvalidRequest
		}

		bz, err = json.Marshal(RouteResponse{
			Route: convertCyberRouteToWasmRoute(route),
		})
	default:
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Liquidity variant"}
	}

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func convertCyberRoutesToWasmRoutes(routes types.Routes) Routes {
	converted := make(Routes, len(routes))
	for i, c := range routes {
		converted[i] = convertCyberRouteToWasmRoute(c)
	}
	return converted
}

func convertCyberRouteToWasmRoute(route types.Route) Route {
	return Route{
		route.Source,
		route.Destination,
		route.Name,
		wasmkeeper.ConvertSdkCoinsToWasmCoins(route.Value),
	}
}
