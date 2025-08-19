package keeper

import (
	"context"
	"strconv"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"cosmossdk.io/errors"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"

	ctypes "github.com/cybercongress/go-cyber/v6/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v6/x/resources/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper}
}

var _ types.MsgServer = msgServer{}

func (server msgServer) Investmint(goCtx context.Context, msg *types.MsgInvestmint) (*types.MsgInvestmintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	neuron, err := sdk.AccAddressFromBech32(msg.Neuron)
	if err != nil {
		return nil, err
	}

	switch msg.Resource {
	case ctypes.VOLT:
		if msg.Amount.Denom != server.GetParams(ctx).BaseInvestmintAmountVolt.Denom {
			return nil, errors.Wrap(types.ErrInvalidBaseResource, msg.Amount.Denom)
		}
	case ctypes.AMPERE:
		if msg.Amount.Denom != server.GetParams(ctx).BaseInvestmintAmountAmpere.Denom {
			return nil, errors.Wrap(types.ErrInvalidBaseResource, msg.Amount.Denom)
		}
	}

	minted, err := server.ConvertResource(ctx, neuron, msg.Amount, msg.Resource, msg.Length)
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

func (server msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if server.authority != req.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", server.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := server.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
