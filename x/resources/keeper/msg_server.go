package keeper

import (
	"context"
	"strconv"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	ctypes "github.com/cybercongress/go-cyber/types"

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

	neuron, err := sdk.AccAddressFromBech32(msg.Neuron)
	if err != nil {
		return nil, err
	}

	switch msg.Resource {
	case ctypes.VOLT:
		if msg.Amount.Denom != k.BaseInvestmintAmountVolt(ctx).Denom {
			return nil, sdkerrors.Wrap(types.ErrInvalidBaseResource, msg.Amount.Denom)
		}
	case ctypes.AMPERE:
		if msg.Amount.Denom != k.BaseInvestmintAmountAmpere(ctx).Denom {
			return nil, sdkerrors.Wrap(types.ErrInvalidBaseResource, msg.Amount.Denom)
		}
	}

	err, minted := k.ConvertResource(ctx, neuron, msg.Amount, msg.Resource, msg.Length)
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
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Neuron),
		),
		sdk.NewEvent(
			types.EventTypeInvestmint,
			sdk.NewAttribute(types.AttributeKeyNeuron, msg.Neuron),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyResource, msg.Resource),
			sdk.NewAttribute(types.AttributeKeyLength, strconv.FormatUint(msg.Length, 10)),
			sdk.NewAttribute(types.AttributeKeyMinted, minted.Amount.String()),
		),
	})

	return &types.MsgInvestmintResponse{}, nil
}
