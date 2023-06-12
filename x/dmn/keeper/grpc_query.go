package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/x/dmn/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(goCtx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) Thought(goCtx context.Context, request *types.QueryThoughtParamsRequest) (*types.QueryThoughtResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if request.Program == "" {
		return nil, status.Errorf(codes.InvalidArgument, "program address cannot be empty")
	}

	if request.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "thought name cannot be empty")
	}

	program, err := sdk.AccAddressFromBech32(request.Program)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	thought, found := k.GetThought(ctx, program, request.Name)
	if !found {
		return nil, status.Errorf(codes.NotFound, "thought with program %s and name %s not found", request.Program, request.Name)
	}

	return &types.QueryThoughtResponse{Thought: thought}, nil
}

func (k Keeper) ThoughtStats(goCtx context.Context, request *types.QueryThoughtParamsRequest) (*types.QueryThoughtStatsResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if request.Program == "" {
		return nil, status.Errorf(codes.InvalidArgument, "program address cannot be empty")
	}

	if request.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "thought name cannot be empty")
	}

	program, err := sdk.AccAddressFromBech32(request.Program)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	thoughtStats, found := k.GetThoughtStats(ctx, program, request.Name)
	if !found {
		return nil, status.Errorf(codes.NotFound, "thought stats with program %s and name %s not found", request.Program, request.Name)
	}

	return &types.QueryThoughtStatsResponse{ThoughtStats: thoughtStats}, nil
}

func (k Keeper) Thoughts(goCtx context.Context, _ *types.QueryThoughtsRequest) (*types.QueryThoughtsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	thoughts := k.GetAllThoughts(ctx)

	return &types.QueryThoughtsResponse{Thoughts: thoughts}, nil
}

func (k Keeper) ThoughtsStats(goCtx context.Context, _ *types.QueryThoughtsStatsRequest) (*types.QueryThoughtsStatsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	thoughtsStats := k.GetAllThoughtsStats(ctx)

	return &types.QueryThoughtsStatsResponse{ThoughtsStats: thoughtsStats}, nil
}
