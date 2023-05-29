package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v2/x/grid/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(
	keeper Keeper,
) types.MsgServer {
	return &msgServer{
		keeper,
	}
}

func (k msgServer) CreateRoute(goCtx context.Context, msg *types.MsgCreateRoute) (*types.MsgCreateRouteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	src, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		return nil, err
	}
	dst, err := sdk.AccAddressFromBech32(msg.Destination)
	if err != nil {
		return nil, err
	}

	err = k.Keeper.CreateEnergyRoute(ctx, src, dst, msg.Name)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Source),
		),
		sdk.NewEvent(
			types.EventTypeCreateRoute,
			sdk.NewAttribute(types.AttributeKeySource, msg.Source),
			sdk.NewAttribute(types.AttributeKeyDestination, msg.Destination),
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
		),
	})

	return &types.MsgCreateRouteResponse{}, nil
}

func (k msgServer) EditRoute(goCtx context.Context, msg *types.MsgEditRoute) (*types.MsgEditRouteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	src, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		return nil, err
	}
	dst, err := sdk.AccAddressFromBech32(msg.Destination)
	if err != nil {
		return nil, err
	}

	err = k.Keeper.EditEnergyRoute(ctx, src, dst, msg.Value)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Source),
		),
		sdk.NewEvent(
			types.EventTypeEditRoute,
			sdk.NewAttribute(types.AttributeKeySource, msg.Source),
			sdk.NewAttribute(types.AttributeKeyDestination, msg.Destination),
			sdk.NewAttribute(types.AttributeKeyValue, msg.Value.String()),
		),
	})

	return &types.MsgEditRouteResponse{}, nil
}

func (k msgServer) DeleteRoute(goCtx context.Context, msg *types.MsgDeleteRoute) (*types.MsgDeleteRouteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	src, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		return nil, err
	}
	dst, err := sdk.AccAddressFromBech32(msg.Destination)
	if err != nil {
		return nil, err
	}

	err = k.Keeper.DeleteEnergyRoute(ctx, src, dst)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Source),
		),
		sdk.NewEvent(
			types.EventTypeDeleteRoute,
			sdk.NewAttribute(types.AttributeKeySource, msg.Source),
			sdk.NewAttribute(types.AttributeKeyDestination, msg.Destination),
		),
	})

	return &types.MsgDeleteRouteResponse{}, nil
}

func (k msgServer) EditRouteName(goCtx context.Context, msg *types.MsgEditRouteName) (*types.MsgEditRouteNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	src, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		return nil, err
	}
	dst, err := sdk.AccAddressFromBech32(msg.Destination)
	if err != nil {
		return nil, err
	}

	err = k.Keeper.EditEnergyRouteName(ctx, src, dst, msg.Name)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Source),
		),
		sdk.NewEvent(
			types.EventTypeEditRouteName,
			sdk.NewAttribute(types.AttributeKeySource, msg.Source),
			sdk.NewAttribute(types.AttributeKeyDestination, msg.Destination),
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
		),
	})

	return &types.MsgEditRouteNameResponse{}, nil
}
