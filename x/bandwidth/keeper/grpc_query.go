package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/bandwidth/types"
	"context"
)

var _ types.QueryServer = &BandwidthMeter{}

func (bk BandwidthMeter) Load(goCtx context.Context, request *types.QueryLoadRequest) (*types.QueryLoadResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	load := bk.GetCurrentNetworkLoad(ctx)

	// TODO refactor to DEC
	return &types.QueryLoadResponse{Load: sdk.DecProto{Dec: load}}, nil
}

func (bk *BandwidthMeter) Price(_ context.Context, request *types.QueryPriceRequest) (*types.QueryPriceResponse, error) {
	price := bk.GetCurrentCreditPrice()

	// TODO refactor to DEC
	return &types.QueryPriceResponse{Price: sdk.DecProto{Dec: price}}, nil
}

func (bk BandwidthMeter) Account(goCtx context.Context, request *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, _ := sdk.AccAddressFromBech32(request.Address)
	ab := bk.GetCurrentAccountBandwidth(ctx, addr)

	return &types.QueryAccountResponse{AccountBandwidth: ab}, nil
}

func (bk BandwidthMeter) Params(goCtx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := bk.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}
