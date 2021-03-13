package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/energy/types"
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

func (k msgServer) CreateEnergyRoute(goCtx context.Context, msg *types.MsgCreateEnergyRoute) (*types.MsgCreateEnergyRouteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	src, _ := sdk.AccAddressFromBech32(msg.Source)
	dst, _ := sdk.AccAddressFromBech32(msg.Destination)
	err := k.Keeper.CreateEnergyRoute(ctx, src, dst, msg.Alias)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateEnergyRoute,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySource, msg.Source),
			sdk.NewAttribute(types.AttributeKeyDestination, msg.Destination),
			sdk.NewAttribute(types.AttributeKeyAlias, msg.Alias),
		),
	)

	return &types.MsgCreateEnergyRouteResponse{}, nil
}

func (k msgServer) EditEnergyRoute(goCtx context.Context, msg *types.MsgEditEnergyRoute) (*types.MsgEditEnergyRouteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	src, _ := sdk.AccAddressFromBech32(msg.Source)
	dst, _ := sdk.AccAddressFromBech32(msg.Destination)

	err := k.Keeper.EditEnergyRoute(ctx, src, dst, msg.Value)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEditEnergyRoute,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySource, msg.Source),
			sdk.NewAttribute(types.AttributeKeyDestination, msg.Destination),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Value.String()),
		),
	)

	return &types.MsgEditEnergyRouteResponse{}, nil
}

func (k msgServer) DeleteEnergyRoute(goCtx context.Context, msg *types.MsgDeleteEnergyRoute) (*types.MsgDeleteEnergyRouteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	src, _ := sdk.AccAddressFromBech32(msg.Source)
	dst, _ := sdk.AccAddressFromBech32(msg.Destination)

	err := k.Keeper.DeleteEnergyRoute(ctx, src, dst)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDeleteEnergyRoute,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySource, msg.Source),
			sdk.NewAttribute(types.AttributeKeyDestination, msg.Destination),
		),
	)

	return &types.MsgDeleteEnergyRouteResponse{}, nil
}

func (k msgServer) EditEnergyRouteAlias(goCtx context.Context, msg *types.MsgEditEnergyRouteAlias) (*types.MsgEditEnergyRouteAliasResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	src, _ := sdk.AccAddressFromBech32(msg.Source)
	dst, _ := sdk.AccAddressFromBech32(msg.Destination)

	err := k.Keeper.EditEnergyRouteAlias(ctx, src, dst, msg.Alias)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEditEnergyRouteAlias,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySource, msg.Source),
			sdk.NewAttribute(types.AttributeKeyDestination, msg.Destination),
			sdk.NewAttribute(types.AttributeKeyAlias, msg.Alias),
		),
	)

	return &types.MsgEditEnergyRouteAliasResponse{}, nil
}
