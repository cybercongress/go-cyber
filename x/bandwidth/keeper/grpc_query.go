package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"context"

	"github.com/cybercongress/go-cyber/x/bandwidth/types"
)

var _ types.QueryServer = &BandwidthMeter{}

func (bk BandwidthMeter) Params(goCtx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := bk.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (bk *BandwidthMeter) Load(goCtx context.Context, request *types.QueryLoadRequest) (*types.QueryLoadResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	load := bk.GetCurrentNetworkLoad(ctx)

	return &types.QueryLoadResponse{Load: sdk.DecProto{Dec: load}}, nil
}

func (bk *BandwidthMeter) Price(_ context.Context, _ *types.QueryPriceRequest) (*types.QueryPriceResponse, error) {
	price := bk.GetCurrentCreditPrice()

	return &types.QueryPriceResponse{Price: sdk.DecProto{Dec: price}}, nil
}

func (bk *BandwidthMeter) DesirableBandwidth(goCtx context.Context, _ *types.QueryDesirableBandwidthRequest) (*types.QueryDesirableBandwidthResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	desirableBandwidth := bk.GetDesirableBandwidth(ctx)

	return &types.QueryDesirableBandwidthResponse{DesirableBandwidth: desirableBandwidth}, nil
}

func (bk *BandwidthMeter) Account(goCtx context.Context, request *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if request.Address == "" {
		return nil, status.Errorf(codes.InvalidArgument, "source address cannot be empty")
	}

	addr, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	accountBandwidth := bk.GetCurrentAccountBandwidth(ctx, addr)

	return &types.QueryAccountResponse{AccountBandwidth: accountBandwidth}, nil
}
