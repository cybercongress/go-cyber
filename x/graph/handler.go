package graph

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	cyberbankkeeper "github.com/cybercongress/go-cyber/x/cyberbank/keeper"
	//сtypes "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/graph/types"
	"github.com/cybercongress/go-cyber/x/graph/keeper"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
)

func NewHandler(
	gk keeper.GraphKeeper,
	ik *keeper.IndexKeeper,
	ak authkeeper.AccountKeeper,
	bk *cyberbankkeeper.IndexedKeeper,
) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(gk, ik, ak, bk)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCyberlink:
			res, err :=  msgServer.Cyberlink(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest,
				"unrecognized %s message type: %v", types.ModuleName, msg.Type())
		}
	}
}
