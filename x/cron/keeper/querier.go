package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cybercongress/go-cyber/x/cron/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, k, legacyQuerierCdc)
		case types.QueryJob:
			return queryJob(ctx, req, k, legacyQuerierCdc)
		case types.QueryJobStats:
			return queryJobStats(ctx, req, k, legacyQuerierCdc)
		case types.QueryJobs:
			return queryJobs(ctx, req, k, legacyQuerierCdc)
		case types.QueryJobsStats:
			return queryJobsStats(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown cron query endpoint")
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

func queryJob(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryJobParamsRequest

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	ct, _ := sdk.AccAddressFromBech32(params.Contract)
	cr, _ := sdk.AccAddressFromBech32(params.Creator)
	job, _ := k.GetJob(ctx, ct, cr, params.Label)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, job)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryJobStats(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryJobParamsRequest

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	ct, _ := sdk.AccAddressFromBech32(params.Contract)
	cr, _ := sdk.AccAddressFromBech32(params.Creator)
	jobStats, _ := k.GetJobStats(ctx, ct, cr, params.Label)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, jobStats)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryJobs(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	jobs := k.GetAllJobs(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, jobs)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryJobsStats(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	jobs := k.GetAllJobsStats(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, jobs)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
