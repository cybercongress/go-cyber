package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ctypes "github.com/cybercongress/go-cyber/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cybercongress/go-cyber/x/resources/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(goCtx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) Investmint(goCtx context.Context, request *types.QueryInvestmintRequest) (*types.QueryInvestmintResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if request.Amount.Denom != ctypes.SCYB {
		return nil, sdkerrors.Wrap(types.ErrInvalidBaseResource, request.Amount.String())
	}

	if request.Resource != ctypes.VOLT && request.Resource != ctypes.AMPERE {
		return nil, sdkerrors.Wrap(types.ErrResourceNotExist, request.Resource)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	amount := k.CalculateInvestmint(ctx, request.Amount, request.Resource, request.Length)

	return &types.QueryInvestmintResponse{Amount: amount}, nil
}