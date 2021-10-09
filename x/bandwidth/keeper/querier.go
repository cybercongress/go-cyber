package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cybercongress/go-cyber/x/bandwidth/types"
)

func NewQuerier(bm *BandwidthMeter, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParams(ctx, req, *bm, legacyQuerierCdc)
		case types.QueryLoad:
			return queryLoad(ctx, req, bm, legacyQuerierCdc)
		case types.QueryPrice:
			return queryPrice(ctx, req, bm, legacyQuerierCdc)
		case types.QueryDesirableBandwidth:
			return queryTotalBandwidth(ctx, req, bm, legacyQuerierCdc)
		case types.QueryAccount:
			return queryNeuronBandwidth(ctx, req, *bm, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, _ abci.RequestQuery, bm BandwidthMeter, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := bm.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryLoad(ctx sdk.Context, _ abci.RequestQuery, bm *BandwidthMeter, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	load := bm.GetCurrentNetworkLoad(ctx)
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryLoadResponse{Load: sdk.DecProto{Dec: load}})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryPrice(_ sdk.Context, _ abci.RequestQuery, bm *BandwidthMeter, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	price := bm.GetCurrentCreditPrice()
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryPriceResponse{Price: sdk.DecProto{Dec: price}})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryTotalBandwidth(ctx sdk.Context, _ abci.RequestQuery, bm *BandwidthMeter, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	totalBandwidth := bm.GetDesirableBandwidth(ctx)
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryTotalBandwidthResponse{TotalBandwidth: totalBandwidth})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryNeuronBandwidth(ctx sdk.Context, req abci.RequestQuery, bm BandwidthMeter, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	var params types.QueryAccountBandwidthParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	neuronBandwidth := bm.GetCurrentAccountBandwidth(ctx, params.Address)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryAccountResponse{NeuronBandwidth: neuronBandwidth})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}


