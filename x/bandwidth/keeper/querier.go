package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cybercongress/go-cyber/x/bandwidth/types"
)

func NewQuerier(bm *BandwidthMeter) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParams(ctx, *bm)
		case types.QueryLoad:
			return queryLoad(ctx, bm)
		case types.QueryPrice:
			return queryPrice(ctx, bm)
		case types.QueryAccount:
			return queryAccount(ctx, req, *bm)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, bm BandwidthMeter) ([]byte, error) {
	params := bm.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryLoad(ctx sdk.Context, bm *BandwidthMeter) ([]byte, error) {
	load := bm.GetCurrentNetworkLoad(ctx)
	res, err := codec.MarshalJSONIndent(codec.Cdc, types.NewResultLoad(load))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryPrice(ctx sdk.Context, bm *BandwidthMeter) ([]byte, error) {
	price := bm.GetCurrentCreditPrice()
	res, err := codec.MarshalJSONIndent(codec.Cdc, types.NewResultPrice(price))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryAccount(ctx sdk.Context, req abci.RequestQuery, bm BandwidthMeter) ([]byte, error) {
	var params types.QueryAccountParams

	err := codec.Cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	account := bm.GetCurrentAccountBandwidth(ctx, params.Account)

	res, err := codec.MarshalJSONIndent(codec.Cdc, account)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}


