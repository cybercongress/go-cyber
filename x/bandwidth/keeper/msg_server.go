package keeper

import (
	"context"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v5/x/bandwidth/types"
)

type msgServer struct {
	BandwidthMeter
}

func NewMsgServerImpl(meter BandwidthMeter) types.MsgServer {
	return &msgServer{BandwidthMeter: meter}
}

var _ types.MsgServer = msgServer{}

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
