package keeper

import (
	"context"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
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

	err, minted := k.ConvertResource(ctx, agent, msg.Amount, msg.Resource, msg.Length)
	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounterWithLabels(
		[]string{types.ModuleName, "investmint"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("resource", msg.Resource),
		},
	)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeInvestmint,
			sdk.NewAttribute(types.AttributeKeyAgent, msg.Agent),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyResource, msg.Resource),
			sdk.NewAttribute(types.AttributeKeyLength, strconv.FormatUint(msg.Length, 10)),
			sdk.NewAttribute(types.AttributeKeyMinted, minted.Amount.String()),
		),
	})

	return &types.MsgInvestmintResponse{}, nil
}