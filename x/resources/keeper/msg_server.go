package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/x/resources/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper}
}

func (k msgServer) Investmint(goCtx context.Context, msg *types.MsgInvestmint) (*types.MsgInvestmintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	agent, err := sdk.AccAddressFromBech32(msg.Agent)
	if err != nil {
		return nil, err
	}

	err = k.ConvertResource(ctx, agent, msg.Amount, msg.Resource, msg.Length)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeConvert,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyAgent, msg.Agent),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyResource, msg.Resource),
			sdk.NewAttribute(types.AttributeKeyEndTime, strconv.FormatUint(msg.Length, 10)),
		),
	)

	return &types.MsgInvestmintResponse{}, nil
}