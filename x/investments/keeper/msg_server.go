package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/investments/types"
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

func (k msgServer) Invest(goCtx context.Context, msg *types.MsgInvest) (*types.MsgInvestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	investor, err := sdk.AccAddressFromBech32(msg.Investor)
	if err != nil {
		return nil, err
	}

	err = k.PutInvestment(ctx, investor, msg.Amount, msg.Resource, msg.EndTime) // TODO return minted amount and pass to event
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeInvest,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyInvestor, msg.Investor),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyResource, msg.Resource),
			sdk.NewAttribute(types.AttributeKeyEndTime, strconv.FormatInt(msg.EndTime, 10)),
		),
	)

	return &types.MsgInvestResponse{}, nil
}

func (k msgServer) CreateResource(goCtx context.Context, msg *types.MsgCreateResource) (*types.MsgCreateResourceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, err
	}

	err = k.PutCreateResource(ctx, msg.Resource, sender, receiver)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
			sdk.NewAttribute(types.AttributeKeyResource, msg.Resource.String()),
		),
	)

	return &types.MsgCreateResourceResponse{}, nil
}

func (k msgServer) RedeemResource(goCtx context.Context, msg *types.MsgRedeemResource) (*types.MsgRedeemResourceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = k.PutRedeemResource(ctx, msg.Resource, sender)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyResource, msg.Resource.String()),
		),
	)

	return &types.MsgRedeemResourceResponse{}, nil
}

