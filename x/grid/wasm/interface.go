package wasm

import (
	"encoding/json"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	wasmplugins "github.com/cybercongress/go-cyber/v2/plugins"
	"github.com/cybercongress/go-cyber/v2/x/grid/keeper"
	"github.com/cybercongress/go-cyber/v2/x/grid/types"
)

var _ WasmMsgParserInterface = WasmMsgParser{}

//--------------------------------------------------

type WasmMsgParserInterface interface {
	Parse(contractAddr sdk.AccAddress, msg wasmvmtypes.CosmosMsg) ([]sdk.Msg, error)
	ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error)
}

type WasmMsgParser struct{}

func NewWasmMsgParser() WasmMsgParser {
	return WasmMsgParser{}
}

func (WasmMsgParser) Parse(_ sdk.AccAddress, _ wasmvmtypes.CosmosMsg) ([]sdk.Msg, error) {
	return nil, nil
}

type CosmosMsg struct {
	CreateEnergyRoute   *types.MsgCreateRoute   `json:"create_energy_route,omitempty"`
	EditEnergyRoute     *types.MsgEditRoute     `json:"edit_energy_route,omitempty"`
	EditEnergyRouteName *types.MsgEditRouteName `json:"edit_energy_route_name,omitempty"`
	DeleteEnergyRoute   *types.MsgDeleteRoute   `json:"delete_energy_route,omitempty"`
}

func (WasmMsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var sdkMsg CosmosMsg
	err := json.Unmarshal(data, &sdkMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to parse link custom msg")
	}

	if sdkMsg.CreateEnergyRoute != nil {
		return []sdk.Msg{sdkMsg.CreateEnergyRoute}, sdkMsg.CreateEnergyRoute.ValidateBasic()
	} else if sdkMsg.EditEnergyRoute != nil {
		return []sdk.Msg{sdkMsg.EditEnergyRoute}, sdkMsg.EditEnergyRoute.ValidateBasic()
	} else if sdkMsg.EditEnergyRouteName != nil {
		return []sdk.Msg{sdkMsg.EditEnergyRouteName}, sdkMsg.EditEnergyRouteName.ValidateBasic()
	} else if sdkMsg.EditEnergyRoute != nil {
		return []sdk.Msg{sdkMsg.EditEnergyRoute}, sdkMsg.EditEnergyRoute.ValidateBasic()
	} else if sdkMsg.EditEnergyRouteName != nil {
		return []sdk.Msg{sdkMsg.EditEnergyRouteName}, sdkMsg.EditEnergyRouteName.ValidateBasic()
	} else if sdkMsg.DeleteEnergyRoute != nil {
		return []sdk.Msg{sdkMsg.DeleteEnergyRoute}, sdkMsg.DeleteEnergyRoute.ValidateBasic()
	}

	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Energy")
}

//--------------------------------------------------

type WasmQuerierInterface interface {
	Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

type WasmQuerier struct {
	keeper.Keeper
}

func NewWasmQuerier(keeper keeper.Keeper) WasmQuerier {
	return WasmQuerier{keeper}
}

func (WasmQuerier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) { return nil, nil }

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

func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	if query.SourceRoutes != nil {
		source, _ := sdk.AccAddressFromBech32(query.SourceRoutes.Source)
		routes := querier.Keeper.GetSourceRoutes(ctx, source, 16)

		bz, err = json.Marshal(RoutesResponse{
			Routes: convertCyberRoutesToWasmRoutes(routes),
		})
	} else if query.SourceRoutedEnergy != nil {
		source, _ := sdk.AccAddressFromBech32(query.SourceRoutedEnergy.Source)
		value := querier.Keeper.GetRoutedFromEnergy(ctx, source)

		bz, err = json.Marshal(RoutedEnergyResponse{
			Value: wasmplugins.ConvertSdkCoinsToWasmCoins(value),
		})
	} else if query.DestinationRoutedEnergy != nil {
		destination, _ := sdk.AccAddressFromBech32(query.DestinationRoutedEnergy.Destination)
		value := querier.Keeper.GetRoutedToEnergy(ctx, destination)

		bz, err = json.Marshal(RoutedEnergyResponse{
			Value: wasmplugins.ConvertSdkCoinsToWasmCoins(value),
		})
	} else if query.Route != nil {
		source, _ := sdk.AccAddressFromBech32(query.Route.Source)
		destination, _ := sdk.AccAddressFromBech32(query.Route.Destination)
		route, found := querier.Keeper.GetRoute(ctx, source, destination)
		if !found {
			return nil, sdkerrors.ErrInvalidRequest
		}

		bz, err = json.Marshal(RouteResponse{
			Route: convertCyberRouteToWasmRoute(route),
		})
	} else {
		return nil, sdkerrors.ErrInvalidRequest
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
		wasmplugins.ConvertSdkCoinsToWasmCoins(route.Value),
	}
}
