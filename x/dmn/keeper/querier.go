package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cybercongress/go-cyber/x/dmn/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, k, legacyQuerierCdc)
		case types.QueryThought:
			return queryThought(ctx, req, k, legacyQuerierCdc)
		case types.QueryThoughtStats:
			return queryThoughtStats(ctx, req, k, legacyQuerierCdc)
		case types.QueryThoughts:
			return queryThoughts(ctx, req, k, legacyQuerierCdc)
		case types.QueryThoughtsStats:
			return queryThoughtsStats(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown dmn query endpoint")
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

func queryThought(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryThoughtParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	thought, found := k.GetThought(ctx, params.Program, params.Name)
	if !found {
		return nil, types.ErrThoughtNotExist
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, thought)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryThoughtStats(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryThoughtParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	thoughtStats, found := k.GetThoughtStats(ctx, params.Program, params.Name)
	if !found {
		return nil, types.ErrThoughtNotExist
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, thoughtStats)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryThoughts(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	thoughts := k.GetAllThoughts(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, thoughts)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryThoughtsStats(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	thoughtsStats := k.GetAllThoughtsStats(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, thoughtsStats)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
