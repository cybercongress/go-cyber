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

func (k msgServer) CronAddJob(goCtx context.Context, msg *types.MsgCronAddJob) (*types.MsgCronAddJobResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.AddJob(
		ctx, cr, ct,
		*msg.Trigger, *msg.Load,
		msg.Label, graph.Cid(msg.Cid),
	)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddJob,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract, msg.Contract),
			sdk.NewAttribute(types.AttributeKeyJobTrigger, msg.Trigger.String()),
			sdk.NewAttribute(types.AttributeKeyJobLoad, msg.Load.String()),
			sdk.NewAttribute(types.AttributeKeyJobLabel, msg.Label),
			sdk.NewAttribute(types.AttributeKeyJobCID, string(msg.Cid)),
		),
	)

	return &types.MsgCronAddJobResponse{}, nil
}

func (k msgServer) CronRemoveJob(goCtx context.Context, msg *types.MsgCronRemoveJob) (*types.MsgCronRemoveJobResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.RemoveJob(ctx, cr, ct, msg.Label)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveJob,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract, msg.Contract),
		),
	)

	return &types.MsgCronRemoveJobResponse{}, nil
}

func (k msgServer) CronChangeJobCID(goCtx context.Context, msg *types.MsgCronChangeJobCID) (*types.MsgCronChangeJobCIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.ChangeJobCID(ctx, cr, ct, msg.Label, graph.Cid(msg.Cid))
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeJobCID,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract,msg.Contract),
			sdk.NewAttribute(types.AttributeKeyJobCID, msg.Cid),
		),
	)

	return &types.MsgCronChangeJobCIDResponse{}, nil
}

func (k msgServer) CronChangeJobLabel(goCtx context.Context, msg *types.MsgCronChangeJobLabel) (*types.MsgCronChangeJobLabelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.ChangeJobLabel(ctx, cr, ct, msg.Label, msg.NewLabel)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeJobLabel,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract, msg.Contract),
			sdk.NewAttribute(types.AttributeKeyJobLabel, msg.Label),
		),
	)

	return &types.MsgCronChangeJobLabelResponse{}, nil
}

func (k msgServer) CronChangeJobCallData(goCtx context.Context, msg *types.MsgCronChangeJobCallData) (*types.MsgCronChangeJobCallDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.ChangeJobCallData(ctx, cr, ct, msg.Label, msg.CallData)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeJobCallData,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract, msg.Contract),
			sdk.NewAttribute(types.AttributeKeyJobCallData, msg.CallData),
		),
	)

	return &types.MsgCronChangeJobCallDataResponse{}, nil
}

func (k msgServer) CronChangeJobGasPrice(goCtx context.Context, msg *types.MsgCronChangeJobGasPrice) (*types.MsgCronChangeJobGasPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.ChangeJobGasPrice(ctx, cr, ct, msg.Label, msg.GasPrice)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeJobGasPrice,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract, msg.Contract),
			sdk.NewAttribute(types.AttributeKeyJobGasPrice, msg.GasPrice.String()),
		),
	)

	return &types.MsgCronChangeJobGasPriceResponse{}, nil
}

func (k msgServer) CronChangeJobPeriod(goCtx context.Context, msg *types.MsgCronChangeJobPeriod) (*types.MsgCronChangeJobPeriodResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.ChangeJobPeriod(ctx, cr, ct, msg.Label, msg.Period)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeJobPeriod,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract, msg.Contract),
			sdk.NewAttribute(types.AttributeKeyJobPeriod, string(msg.Period)),
		),
	)

	return &types.MsgCronChangeJobPeriodResponse{}, nil
}

func (k msgServer) CronChangeJobBlock(goCtx context.Context, msg *types.MsgCronChangeJobBlock) (*types.MsgCronChangeJobBlockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.ChangeJobBlock(ctx, cr, ct, msg.Label, msg.Block)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeJobBlock,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract, msg.Contract),
			sdk.NewAttribute(types.AttributeKeyJobBlock, string(msg.Block)),
		),
	)

	return &types.MsgCronChangeJobBlockResponse{}, nil
}