package energy

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cybercongress/go-cyber/x/energy/keeper"
	"github.com/cybercongress/go-cyber/x/energy/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreateEnergyRoute:
			return handleMsgCreateEnergyRoute(ctx, k, msg)
		case types.MsgEditEnergyRoute:
			return handleMsgEditEnergyRoute(ctx, k, msg)
		case types.MsgDeleteEnergyRoute:
			return handleMsgDeleteEnergyRoute(ctx, k, msg)
		case types.MsgEditEnergyRouteAlias:
			return handleMsgEditEnergyRouteAlias(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName,  msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreateEnergyRoute(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateEnergyRoute) (*sdk.Result, error) {
	err := k.CreateEnergyRoute(ctx, msg.Source, msg.Destination, msg.Alias)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateEnergyRoute,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySource, msg.Source.String()),
			sdk.NewAttribute(types.AttributeKeyDestination, msg.Destination.String()),
			sdk.NewAttribute(types.AttributeKeyAlias, msg.Alias),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgEditEnergyRoute(ctx sdk.Context, k keeper.Keeper, msg types.MsgEditEnergyRoute) (*sdk.Result, error) {
	err := k.EditEnergyRoute(ctx, msg.Source, msg.Destination, msg.Value)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEditEnergyRoute,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySource, msg.Source.String()),
			sdk.NewAttribute(types.AttributeKeyDestination, msg.Destination.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Value.Amount.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDeleteEnergyRoute(ctx sdk.Context, k keeper.Keeper, msg types.MsgDeleteEnergyRoute) (*sdk.Result, error) {
	err := k.DeleteEnergyRoute(ctx, msg.Source, msg.Destination)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDeleteEnergyRoute,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySource, msg.Source.String()),
			sdk.NewAttribute(types.AttributeKeyDestination, msg.Destination.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgEditEnergyRouteAlias(ctx sdk.Context, k keeper.Keeper, msg types.MsgEditEnergyRouteAlias) (*sdk.Result, error) {
	err := k.EditEnergyRouteAlias(ctx, msg.Source, msg.Destination, msg.Alias)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEditEnergyRouteAlias,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySource, msg.Source.String()),
			sdk.NewAttribute(types.AttributeKeyDestination, msg.Destination.String()),
			sdk.NewAttribute(types.AttributeKeyAlias, msg.Alias),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
