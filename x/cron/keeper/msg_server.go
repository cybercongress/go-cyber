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

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.SaveJob(
		ctx, cr, ct,
		msg.Trigger, msg.Load,
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

	return &types.MsgAddJobResponse{}, nil
}

func (k msgServer) RemoveJob(goCtx context.Context, msg *types.MsgRemoveJob) (*types.MsgRemoveJobResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.RemoveJobFull(ctx, cr, ct, msg.Label)
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

	return &types.MsgRemoveJobResponse{}, nil
}

func (k msgServer) ChangeJobCID(goCtx context.Context, msg *types.MsgChangeJobCID) (*types.MsgChangeJobCIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.UpdateJobCID(ctx, cr, ct, msg.Label, graph.Cid(msg.Cid))
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

	return &types.MsgChangeJobCIDResponse{}, nil
}

func (k msgServer) ChangeJobLabel(goCtx context.Context, msg *types.MsgChangeJobLabel) (*types.MsgChangeJobLabelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.UpdateJobLabel(ctx, cr, ct, msg.Label, msg.NewLabel)
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

	return &types.MsgChangeJobLabelResponse{}, nil
}

func (k msgServer) ChangeJobCallData(goCtx context.Context, msg *types.MsgChangeJobCallData) (*types.MsgChangeJobCallDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.UpdateJobCallData(ctx, cr, ct, msg.Label, msg.CallData)
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

	return &types.MsgChangeJobCallDataResponse{}, nil
}

func (k msgServer) ChangeJobGasPrice(goCtx context.Context, msg *types.MsgChangeJobGasPrice) (*types.MsgChangeJobGasPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.UpdateJobGasPrice(ctx, cr, ct, msg.Label, msg.GasPrice)
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

	return &types.MsgChangeJobGasPriceResponse{}, nil
}

func (k msgServer) ChangeJobPeriod(goCtx context.Context, msg *types.MsgChangeJobPeriod) (*types.MsgChangeJobPeriodResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.UpdateJobPeriod(ctx, cr, ct, msg.Label, msg.Period)
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

	return &types.MsgChangeJobPeriodResponse{}, nil
}

func (k msgServer) ChangeJobBlock(goCtx context.Context, msg *types.MsgChangeJobBlock) (*types.MsgChangeJobBlockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ct, _ := sdk.AccAddressFromBech32(msg.Contract)
	cr, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.UpdateJobBlock(ctx, cr, ct, msg.Label, msg.Block)
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

	return &types.MsgChangeJobBlockResponse{}, nil
}