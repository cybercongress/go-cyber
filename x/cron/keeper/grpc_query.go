package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/cron/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(goCtx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) Job(goCtx context.Context, request *types.QueryJobParamsRequest) (*types.QueryJobResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(request.Contract)
	cr, _ := sdk.AccAddressFromBech32(request.Creator)
	job, _ := k.GetJob(ctx, ct, cr, request.Label)

	return &types.QueryJobResponse{Job: job}, nil
}

func (k Keeper) JobStats(goCtx context.Context, request *types.QueryJobParamsRequest) (*types.QueryJobStatsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(request.Contract)
	cr, _ := sdk.AccAddressFromBech32(request.Creator)
	jobStats, _ := k.GetJobStats(ctx, ct, cr, request.Label)

	return &types.QueryJobStatsResponse{JobStats: jobStats}, nil
}

func (k Keeper) Jobs(goCtx context.Context, _ *types.QueryJobsRequest) (*types.QueryJobsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	jobs := k.GetAllJobs(ctx)

	return &types.QueryJobsResponse{Jobs: jobs}, nil
}

func (k Keeper) JobsStats(goCtx context.Context, _ *types.QueryJobsStatsRequest) (*types.QueryJobsStatsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	jobsStats := k.GetAllJobsStats(ctx)

	return &types.QueryJobsStatsResponse{JobsStats: jobsStats}, nil
}