package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/deep-foundation/deep-chain/x/grid/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(goCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) SourceRoutes(goCtx context.Context, request *types.QuerySourceRequest) (*types.QueryRoutesResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if request.Source == "" {
		return nil, status.Errorf(codes.InvalidArgument, "source address cannot be empty")
	}

	addr, err := sdk.AccAddressFromBech32(request.Source)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	routes := k.GetSourceRoutes(ctx, addr, 16)

	return &types.QueryRoutesResponse{Routes: routes}, nil
}

// DestinationRoutes TODO add pagination
func (k Keeper) DestinationRoutes(goCtx context.Context, request *types.QueryDestinationRequest) (*types.QueryRoutesResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if request.Destination == "" {
		return nil, status.Errorf(codes.InvalidArgument, "destination address cannot be empty")
	}

	addr, err := sdk.AccAddressFromBech32(request.Destination)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	routes := k.GetDestinationRoutes(ctx, addr)

	return &types.QueryRoutesResponse{Routes: routes}, nil
}

func (k Keeper) DestinationRoutedEnergy(goCtx context.Context, request *types.QueryDestinationRequest) (*types.QueryRoutedEnergyResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if request.Destination == "" {
		return nil, status.Errorf(codes.InvalidArgument, "destination address cannot be empty")
	}

	addr, err := sdk.AccAddressFromBech32(request.Destination)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	routedEnergy := k.GetRoutedToEnergy(ctx, addr)

	return &types.QueryRoutedEnergyResponse{Value: routedEnergy}, nil
}

func (k Keeper) SourceRoutedEnergy(goCtx context.Context, request *types.QuerySourceRequest) (*types.QueryRoutedEnergyResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if request.Source == "" {
		return nil, status.Errorf(codes.InvalidArgument, "source address cannot be empty")
	}

	addr, err := sdk.AccAddressFromBech32(request.Source)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	routedEnergy := k.GetRoutedFromEnergy(ctx, addr)

	return &types.QueryRoutedEnergyResponse{Value: routedEnergy}, nil
}

func (k Keeper) Route(goCtx context.Context, request *types.QueryRouteRequest) (*types.QueryRouteResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if request.Source == "" {
		return nil, status.Errorf(codes.InvalidArgument, "source address cannot be empty")
	}
	if request.Destination == "" {
		return nil, status.Errorf(codes.InvalidArgument, "destination address cannot be empty")
	}

	src, err := sdk.AccAddressFromBech32(request.Source)
	if err != nil {
		return nil, err
	}
	dst, err := sdk.AccAddressFromBech32(request.Destination)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	route, found := k.GetRoute(ctx, src, dst)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"route with source %s and destination %s not found",
			request.Source, request.Destination)
	}

	return &types.QueryRouteResponse{Route: route}, nil
}

func (k Keeper) Routes(goCtx context.Context, request *types.QueryRoutesRequest) (*types.QueryRoutesResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if request.Pagination == nil {
		return nil, status.Errorf(codes.InvalidArgument, "pagination cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var routes types.Routes

	store := ctx.KVStore(k.storeKey)
	routesStore := prefix.NewStore(store, types.RouteKey)
	pageRes, err := query.Paginate(routesStore, request.Pagination, func(key []byte, value []byte) error {
		route, err := types.UnmarshalRoute(k.cdc, value)
		if err != nil {
			return err
		}
		routes = append(routes, route)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRoutesResponse{Routes: routes, Pagination: pageRes}, nil
}
