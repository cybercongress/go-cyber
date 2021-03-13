package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cybercongress/go-cyber/x/graph/types"
)

func NewQuerier(gk GraphKeeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryLinksAmount:
			return queryLinksAmount(ctx, req, gk, legacyQuerierCdc)
		case types.QueryCidsAmount:
			return queryCidsAmount(ctx, req, gk, legacyQuerierCdc)
		case types.QueryGraphStats:
			return queryGraphStats(ctx, req, gk, legacyQuerierCdc)
		case types.QueryInLinks:
			return queryInLinks(ctx, req, gk, legacyQuerierCdc)
		case types.QueryOutLinks:
			return queryOutLinks(ctx, req, gk, legacyQuerierCdc)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryLinksAmount(ctx sdk.Context, _ abci.RequestQuery, gk GraphKeeper, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	amount := gk.GetLinksCount(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryLinksAmountResponse{Amount: amount})

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryCidsAmount(ctx sdk.Context, _ abci.RequestQuery, gk GraphKeeper, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	amount := gk.GetCidsCount(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryCidsAmountResponse{Amount: amount})

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryGraphStats(ctx sdk.Context, _ abci.RequestQuery, gk GraphKeeper, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	links := gk.GetLinksCount(ctx)
	cids := gk.GetCidsCount(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryGraphStatsResponse{links, cids})

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// SANDBOX ZONE
func queryInLinks(ctx sdk.Context, req abci.RequestQuery, gk GraphKeeper, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	var params types.QueryLinksParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	inLinks, _, _ := gk.GetAllLinks(ctx)
	cidNum, exist := gk.GetCidNumber(ctx, types.Cid(params.Cid)); if exist != true {
		return nil, sdkerrors.Wrap(types.ErrCidNotFound, "")
	}

	links := inLinks[cidNum]
	data := make([]string, 0)

	for i, _ := range links {
		data = append(data, string(gk.GetCid(ctx, i)))
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryLinksResponse{Cids: data})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// SANDBOX ZONE
func queryOutLinks(ctx sdk.Context, req abci.RequestQuery, gk GraphKeeper, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	var params types.QueryLinksParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	_, outLinks, _ := gk.GetAllLinks(ctx)
	cidNum, exist := gk.GetCidNumber(ctx, types.Cid(params.Cid)); if exist != true {
		return nil, sdkerrors.Wrap(types.ErrCidNotFound, "")
	}

	links := outLinks[cidNum]
	data := make([]string, 0)

	for i, _ := range links {
		data = append(data, string(gk.GetCid(ctx, i)))
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryLinksResponse{Cids: data})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

