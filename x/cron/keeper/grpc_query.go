package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cybercongress/go-cyber/x/cron/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(goCtx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) Job(goCtx context.Context, request *types.QueryJobParamsRequest) (*types.QueryJobResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if request.Program == "" {
		return nil, status.Errorf(codes.InvalidArgument, "program address cannot be empty")
	}

	if request.Label == "" {
		return nil, status.Errorf(codes.InvalidArgument, "job label cannot be empty")
	}

	program, err := sdk.AccAddressFromBech32(request.Program)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	job, found := k.GetJob(ctx, program, request.Label)
	if !found {
		return nil, status.Errorf(codes.NotFound, "job with program %s and label %s not found", request.Program, request.Label)
	}

	return &types.QueryJobResponse{Job: job}, nil
}

func (k Keeper) JobStats(goCtx context.Context, request *types.QueryJobParamsRequest) (*types.QueryJobStatsResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if request.Program == "" {
		return nil, status.Errorf(codes.InvalidArgument, "program address cannot be empty")
	}

	if request.Label == "" {
		return nil, status.Errorf(codes.InvalidArgument, "job label cannot be empty")
	}

	program, err := sdk.AccAddressFromBech32(request.Program)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	jobStats, found := k.GetJobStats(ctx, program, request.Label)
	if !found {
		return nil, status.Errorf(codes.NotFound, "job stats with program %s and label %s not found", request.Program, request.Label)
	}

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