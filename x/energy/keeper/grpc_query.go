package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/energy/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(goCtx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) SourceRoutes(goCtx context.Context, request *types.QuerySourceRequest) (*types.QueryRoutesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, _ := sdk.AccAddressFromBech32(request.Source)
	routes := k.GetSourceRoutes(ctx, addr, 10)

	return &types.QueryRoutesResponse{Routes: routes}, nil
}

func (k Keeper) DestinationRoutes(goCtx context.Context, request *types.QueryDestinationRequest) (*types.QueryRoutesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, _ := sdk.AccAddressFromBech32(request.Destination)
	routes := k.GetDestinationRoutes(ctx, addr)

	return &types.QueryRoutesResponse{Routes: routes}, nil
}

func (k Keeper) DestinationRoutedEnergy(goCtx context.Context, request *types.QueryDestinationRequest) (*types.QueryRoutedEnergyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, _ := sdk.AccAddressFromBech32(request.Destination)

	routedEnergy := k.GetRoutedToEnergy(ctx, addr)

	return &types.QueryRoutedEnergyResponse{Value: routedEnergy}, nil
}

func (k Keeper) SourceRoutedEnergy(goCtx context.Context, request *types.QuerySourceRequest) (*types.QueryRoutedEnergyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, _ := sdk.AccAddressFromBech32(request.Source)

	routedEnergy := k.GetRoutedFromEnergy(ctx, addr)

	return &types.QueryRoutedEnergyResponse{Value: routedEnergy}, nil
}

func (k Keeper) Route(goCtx context.Context, request *types.QueryRouteRequest) (*types.QueryRouteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	src, _ := sdk.AccAddressFromBech32(request.Source)
	dst, _ := sdk.AccAddressFromBech32(request.Destination)

	route, _ := k.GetRoute(ctx, src, dst)

	return &types.QueryRouteResponse{Route: &route}, nil
}

func (k Keeper) Routes(goCtx context.Context, request *types.QueryRoutesRequest) (*types.QueryRoutesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	routes := k.GetAllRoutes(ctx)

	return &types.QueryRoutesResponse{Routes: routes}, nil
}