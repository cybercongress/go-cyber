package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v2/x/dmn/types"
	graph "github.com/cybercongress/go-cyber/v2/x/graph/types"
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

func (k msgServer) CreateThought(goCtx context.Context, msg *types.MsgCreateThought) (*types.MsgCreateThoughtResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.SaveThought(
		ctx, program,
		msg.Trigger, msg.Load,
		msg.Name, graph.Cid(msg.Particle),
	)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Program),
		),
		sdk.NewEvent(
			types.EventTypeCreateThought,
			sdk.NewAttribute(types.AttributeKeyThoughtProgram, msg.Program),
			sdk.NewAttribute(types.AttributeKeyThoughtTrigger, msg.Trigger.String()),
			sdk.NewAttribute(types.AttributeKeyThoughtLoad, msg.Load.String()),
			sdk.NewAttribute(types.AttributeKeyThoughtName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyThoughtParticle, msg.Particle),
		),
	})

	return &types.MsgCreateThoughtResponse{}, nil
}

func (k msgServer) ForgetThought(goCtx context.Context, msg *types.MsgForgetThought) (*types.MsgForgetThoughtResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.RemoveThoughtFull(ctx, program, msg.Name)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Program),
		),
		sdk.NewEvent(
			types.EventTypeForgetThought,
			sdk.NewAttribute(types.AttributeKeyThoughtProgram, msg.Program),
		),
	})

	return &types.MsgForgetThoughtResponse{}, nil
}

func (k msgServer) ChangeThoughtParticle(goCtx context.Context, msg *types.MsgChangeThoughtParticle) (*types.MsgChangeThoughtParticleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.UpdateThoughtParticle(ctx, program, msg.Name, graph.Cid(msg.Particle))
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Program),
		),
		sdk.NewEvent(
			types.EventTypeChangeThoughtParticle,
			sdk.NewAttribute(types.AttributeKeyThoughtProgram, msg.Program),
			sdk.NewAttribute(types.AttributeKeyThoughtParticle, msg.Particle),
		),
	})

	return &types.MsgChangeThoughtParticleResponse{}, nil
}

func (k msgServer) ChangeThoughtName(goCtx context.Context, msg *types.MsgChangeThoughtName) (*types.MsgChangeThoughtNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.UpdateThoughtName(ctx, program, msg.Name, msg.NewName)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Program),
		),
		sdk.NewEvent(
			types.EventTypeChangeThoughtName,
			sdk.NewAttribute(types.AttributeKeyThoughtProgram, msg.Program),
			sdk.NewAttribute(types.AttributeKeyThoughtName, msg.Name),
		),
	})

	return &types.MsgChangeThoughtNameResponse{}, nil
}

func (k msgServer) ChangeThoughtInput(goCtx context.Context, msg *types.MsgChangeThoughtInput) (*types.MsgChangeThoughtInputResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.UpdateThoughtCallData(ctx, program, msg.Name, msg.Input)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Program),
		),
		sdk.NewEvent(
			types.EventTypeChangeThoughtInput,
			sdk.NewAttribute(types.AttributeKeyThoughtProgram, msg.Program),
			sdk.NewAttribute(types.AttributeKeyThoughtInput, msg.Input),
		),
	})

	return &types.MsgChangeThoughtInputResponse{}, nil
}

func (k msgServer) ChangeThoughtGasPrice(goCtx context.Context, msg *types.MsgChangeThoughtGasPrice) (*types.MsgChangeThoughtGasPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.UpdateThoughtGasPrice(ctx, program, msg.Name, msg.GasPrice)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Program),
		),
		sdk.NewEvent(
			types.EventTypeChangeThoughtGasPrice,
			sdk.NewAttribute(types.AttributeKeyThoughtProgram, msg.Program),
			sdk.NewAttribute(types.AttributeKeyThoughtGasPrice, msg.GasPrice.String()),
		),
	})

	return &types.MsgChangeThoughtGasPriceResponse{}, nil
}

func (k msgServer) ChangeThoughtPeriod(goCtx context.Context, msg *types.MsgChangeThoughtPeriod) (*types.MsgChangeThoughtPeriodResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.UpdateThoughtPeriod(ctx, program, msg.Name, msg.Period)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Program),
		),
		sdk.NewEvent(
			types.EventTypeChangeThoughtPeriod,
			sdk.NewAttribute(types.AttributeKeyThoughtProgram, msg.Program),
			sdk.NewAttribute(types.AttributeKeyThoughtPeriod, string(msg.Period)),
		),
	})

	return &types.MsgChangeThoughtPeriodResponse{}, nil
}

func (k msgServer) ChangeThoughtBlock(goCtx context.Context, msg *types.MsgChangeThoughtBlock) (*types.MsgChangeThoughtBlockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.UpdateThoughtBlock(ctx, program, msg.Name, msg.Block)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Program),
		),
		sdk.NewEvent(
			types.EventTypeChangeThoughtBlock,
			sdk.NewAttribute(types.AttributeKeyThoughtProgram, msg.Program),
			sdk.NewAttribute(types.AttributeKeyThoughtBlock, string(msg.Block)),
		),
	})

	return &types.MsgChangeThoughtBlockResponse{}, nil
}
