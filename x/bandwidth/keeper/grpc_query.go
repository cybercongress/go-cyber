package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cybercongress/go-cyber/v6/x/bandwidth/types"
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

	return &types.QueryLoadResponse{Load: load}, nil
}

func (bm *BandwidthMeter) Price(_ context.Context, _ *types.QueryPriceRequest) (*types.QueryPriceResponse, error) {
	price := bm.GetCurrentCreditPrice()

	return &types.QueryPriceResponse{Price: price}, nil
}

func (bm *BandwidthMeter) TotalBandwidth(goCtx context.Context, _ *types.QueryTotalBandwidthRequest) (*types.QueryTotalBandwidthResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	totalBandwidth := bm.GetDesirableBandwidth(ctx)

	return &types.QueryTotalBandwidthResponse{TotalBandwidth: totalBandwidth}, nil
}

func (bm *BandwidthMeter) NeuronBandwidth(goCtx context.Context, request *types.QueryNeuronBandwidthRequest) (*types.QueryNeuronBandwidthResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if request.Neuron == "" {
		return nil, status.Errorf(codes.InvalidArgument, "source address cannot be empty")
	}

	addr, err := sdk.AccAddressFromBech32(request.Neuron)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	neuronBandwidth := bm.GetCurrentVoltsAccountBandwidth(ctx, addr)

	return &types.QueryNeuronBandwidthResponse{NeuronBandwidth: neuronBandwidth}, nil
}
