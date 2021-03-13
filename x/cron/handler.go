package cron

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cybercongress/go-cyber/x/cron/keeper"
	"github.com/cybercongress/go-cyber/x/cron/types"
)

func NewHandler(
	k keeper.Keeper,
) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCronAddJob:
			res, err := msgServer.CronAddJob(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCronRemoveJob:
			res, err := msgServer.CronRemoveJob(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCronChangeJobCID:
			res, err := msgServer.CronChangeJobCID(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCronChangeJobLabel:
			res, err := msgServer.CronChangeJobLabel(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCronChangeJobCallData:
			res, err := msgServer.CronChangeJobCallData(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCronChangeJobGasPrice:
			res, err := msgServer.CronChangeJobGasPrice(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCronChangeJobPeriod:
			res, err := msgServer.CronChangeJobPeriod(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCronChangeJobBlock:
			res, err := msgServer.CronChangeJobBlock(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest,
				"unrecognized %s message type: %v", types.ModuleName, msg.Type())
		}
	}
}

