package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cybercongress/go-cyber/x/energy/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, k)
		case types.QuerySourceRoutes:
			return querySourceRoutes(ctx, req, k)
		case types.QueryDestinationRoutes:
			return queryDestinationRoutes(ctx, req, k)
		case types.QueryDestinationRoutedEnergy:
			return queryDestinationRoutedEnergy(ctx, req, k)
		case types.QuerySourceRoutedEnergy:
			return querySourceRoutedEnergy(ctx, req, k)
		case types.QueryRoute:
			return queryRoute(ctx, req, k)
		case types.QueryRoutes:
			return queryRoutes(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown energy query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func querySourceRoutes(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QuerySourceParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	routes := k.GetSourceRoutes(ctx, params.Source, 10)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, routes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDestinationRoutes(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDestinationParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	routes := k.GetDestinationRoutes(ctx, params.Destination)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, routes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDestinationRoutedEnergy(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDestinationParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	routedEnergy := k.GetRoutedToEnergy(ctx, params.Destination)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, routedEnergy)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func querySourceRoutedEnergy(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QuerySourceParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	routedEnergy := k.GetRoutedFromEnergy(ctx, params.Source)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, routedEnergy)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryRoute(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryRouteParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	route, _ := k.GetRoute(ctx, params.Source, params.Destination)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, route)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryRoutes(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	routes := k.GetAllRoutes(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, routes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
