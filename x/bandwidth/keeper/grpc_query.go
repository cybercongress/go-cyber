package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"context"

	"github.com/cybercongress/go-cyber/x/bandwidth/types"
)

var _ types.QueryServer = &BandwidthMeter{}

func (bm BandwidthMeter) Params(goCtx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := bm.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (bm *BandwidthMeter) Load(goCtx context.Context, request *types.QueryLoadRequest) (*types.QueryLoadResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	load := bm.GetCurrentNetworkLoad(ctx)

	return &types.QueryLoadResponse{Load: sdk.DecProto{Dec: load}}, nil
}

func (bm *BandwidthMeter) Price(_ context.Context, _ *types.QueryPriceRequest) (*types.QueryPriceResponse, error) {
	price := bm.GetCurrentCreditPrice()

	return &types.QueryPriceResponse{Price: sdk.DecProto{Dec: price}}, nil
}

func (bm *BandwidthMeter) TotalBandwidth(goCtx context.Context, _ *types.QueryTotalBandwidthRequest) (*types.QueryTotalBandwidthResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	totalBandwidth := bm.GetDesirableBandwidth(ctx)

	return &types.QueryTotalBandwidthResponse{TotalBandwidth: totalBandwidth}, nil
}

func (bm *BandwidthMeter) Account(goCtx context.Context, request *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
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

	neuronBandwidth := bm.GetCurrentAccountBandwidth(ctx, addr)

	return &types.QueryAccountResponse{NeuronBandwidth: neuronBandwidth}, nil
}
