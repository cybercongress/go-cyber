package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ctypes "github.com/cybercongress/go-cyber/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cybercongress/go-cyber/x/resources/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, k, legacyQuerierCdc)
		case types.QueryInvestmint:
			return queryInvestmint(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown Resources query endpoint")
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

func queryInvestmint(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryInvestmintParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if params.Amount.Denom != ctypes.SCYB {
		return nil, sdkerrors.Wrap(types.ErrInvalidBaseResource, params.Amount.String())
	}

	if params.Resource != ctypes.VOLT && params.Resource != ctypes.AMPERE {
		return nil, sdkerrors.Wrap(types.ErrResourceNotExist, params.Resource)
	}

	routes := k.CalculateInvestmint(ctx, params.Amount, params.Resource, params.Length)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, routes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
