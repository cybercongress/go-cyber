package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/cron/types"
	graph "github.com/cybercongress/go-cyber/x/graph/types"
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

func (k msgServer) AddJob(goCtx context.Context, msg *types.MsgAddJob) (*types.MsgAddJobResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.SaveJob(
		ctx, program,
		msg.Trigger, msg.Load,
		msg.Label, graph.Cid(msg.Cid),
	)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeAddJob,
			sdk.NewAttribute(types.AttributeKeyJobProgram, msg.Program),
			sdk.NewAttribute(types.AttributeKeyJobTrigger, msg.Trigger.String()),
			sdk.NewAttribute(types.AttributeKeyJobLoad, msg.Load.String()),
			sdk.NewAttribute(types.AttributeKeyJobLabel, msg.Label),
			sdk.NewAttribute(types.AttributeKeyJobCID, msg.Cid),
		),
	})

	return &types.MsgAddJobResponse{}, nil
}

func (k msgServer) RemoveJob(goCtx context.Context, msg *types.MsgRemoveJob) (*types.MsgRemoveJobResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.RemoveJobFull(ctx, program, msg.Label)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeRemoveJob,
			sdk.NewAttribute(types.AttributeKeyJobProgram, msg.Program),
		),
	})

	return &types.MsgRemoveJobResponse{}, nil
}

func (k msgServer) ChangeJobCID(goCtx context.Context, msg *types.MsgChangeJobCID) (*types.MsgChangeJobCIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.UpdateJobCID(ctx, program, msg.Label, graph.Cid(msg.Cid))
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeChangeJobCID,
			sdk.NewAttribute(types.AttributeKeyJobProgram, msg.Program),
			sdk.NewAttribute(types.AttributeKeyJobCID, msg.Cid),
		),
	})

	return &types.MsgChangeJobCIDResponse{}, nil
}

func (k msgServer) ChangeJobLabel(goCtx context.Context, msg *types.MsgChangeJobLabel) (*types.MsgChangeJobLabelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.UpdateJobLabel(ctx, program, msg.Label, msg.NewLabel)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeChangeJobLabel,
			sdk.NewAttribute(types.AttributeKeyJobProgram, msg.Program),
			sdk.NewAttribute(types.AttributeKeyJobLabel, msg.Label),
		),
	})

	return &types.MsgChangeJobLabelResponse{}, nil
}

func (k msgServer) ChangeJobCallData(goCtx context.Context, msg *types.MsgChangeJobCallData) (*types.MsgChangeJobCallDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.UpdateJobCallData(ctx, program, msg.Label, msg.CallData)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeChangeJobCallData,
			sdk.NewAttribute(types.AttributeKeyJobProgram, msg.Program),
			sdk.NewAttribute(types.AttributeKeyJobCallData, msg.CallData),
		),
	})

	return &types.MsgChangeJobCallDataResponse{}, nil
}

func (k msgServer) ChangeJobGasPrice(goCtx context.Context, msg *types.MsgChangeJobGasPrice) (*types.MsgChangeJobGasPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.UpdateJobGasPrice(ctx, program, msg.Label, msg.GasPrice)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeChangeJobGasPrice,
			sdk.NewAttribute(types.AttributeKeyJobProgram, msg.Program),
			sdk.NewAttribute(types.AttributeKeyJobGasPrice, msg.GasPrice.String()),
		),
	})

	return &types.MsgChangeJobGasPriceResponse{}, nil
}

func (k msgServer) ChangeJobPeriod(goCtx context.Context, msg *types.MsgChangeJobPeriod) (*types.MsgChangeJobPeriodResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.UpdateJobPeriod(ctx, program, msg.Label, msg.Period)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeChangeJobPeriod,
			sdk.NewAttribute(types.AttributeKeyJobProgram, msg.Program),
			sdk.NewAttribute(types.AttributeKeyJobPeriod, string(msg.Period)),
		),
	})

	return &types.MsgChangeJobPeriodResponse{}, nil
}

func (k msgServer) ChangeJobBlock(goCtx context.Context, msg *types.MsgChangeJobBlock) (*types.MsgChangeJobBlockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, _ := sdk.AccAddressFromBech32(msg.Program)

	err := k.UpdateJobBlock(ctx, program, msg.Label, msg.Block)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeChangeJobBlock,
			sdk.NewAttribute(types.AttributeKeyJobProgram, msg.Program),
			sdk.NewAttribute(types.AttributeKeyJobBlock, string(msg.Block)),
		),
	})

	return &types.MsgChangeJobBlockResponse{}, nil
}