package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cybercongress/go-cyber/x/energy/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, k, legacyQuerierCdc)
		case types.QuerySourceRoutes:
			return querySourceRoutes(ctx, req, k, legacyQuerierCdc)
		case types.QueryDestinationRoutes:
			return queryDestinationRoutes(ctx, req, k, legacyQuerierCdc)
		case types.QueryDestinationRoutedEnergy:
			return queryDestinationRoutedEnergy(ctx, req, k, legacyQuerierCdc)
		case types.QuerySourceRoutedEnergy:
			return querySourceRoutedEnergy(ctx, req, k, legacyQuerierCdc)
		case types.QueryRoute:
			return queryRoute(ctx, req, k, legacyQuerierCdc)
		case types.QueryRoutes:
			return queryRoutes(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown energy query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func querySourceRoutes(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QuerySourceRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	addr, _ := sdk.AccAddressFromBech32(params.Source)

	routes := k.GetSourceRoutes(ctx, addr, 10)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, routes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDestinationRoutes(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryDestinationRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	addr, _ := sdk.AccAddressFromBech32(params.Destination)

	routes := k.GetDestinationRoutes(ctx, addr)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, routes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDestinationRoutedEnergy(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryDestinationRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	addr, _ := sdk.AccAddressFromBech32(params.Destination)

	routedEnergy := k.GetRoutedToEnergy(ctx, addr)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, routedEnergy)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func querySourceRoutedEnergy(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QuerySourceRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	addr, _ := sdk.AccAddressFromBech32(params.Source)

	routedEnergy := k.GetRoutedFromEnergy(ctx, addr)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, routedEnergy)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryRoute(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryRouteRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	src, _ := sdk.AccAddressFromBech32(params.Source)
	dst, _ := sdk.AccAddressFromBech32(params.Destination)

	route, _ := k.GetRoute(ctx, src, dst)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, route)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryRoutes(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	routes := k.GetAllRoutes(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, routes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
