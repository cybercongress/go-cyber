package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	//"github.com/litvintech/cyber/x/link/internal/types"

	"github.com/cybercongress/go-cyber/x/link/types"
)

func NewQuerier(gk GraphKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryLinksAmount:
			return queryLinksAmount(ctx, gk)
		case types.QueryInLinks:
			return queryInLinks(ctx, req, gk)
		case types.QueryOutLinks:
			return queryOutLinks(ctx, req, gk)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryLinksAmount(ctx sdk.Context, gk GraphKeeper) ([]byte, error) {
	amount := gk.GetLinksCount(ctx)

	res, err := codec.MarshalJSONIndent(codec.Cdc, amount)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// TODO for development purpose
func queryInLinks(ctx sdk.Context, req abci.RequestQuery, gk GraphKeeper) ([]byte, error) {
	var params types.QueryLinksParams

	err := codec.Cdc.UnmarshalJSON(req.Data, &params);	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	inLinks, _, _ := gk.GetAllLinks(ctx)
	cidNum, exist := gk.GetCidNumber(ctx, params.Cid); if exist != true {
		return nil, sdkerrors.Wrap(types.ErrCidNotFound, "")
	}

	links := inLinks[cidNum]
	data := make([]types.Cid, 0)

	for i, _ := range links {
		data = append(data, gk.GetCid(ctx, i))
	}

	res, err := codec.MarshalJSONIndent(codec.Cdc, types.ResultLinks{Cids: data})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// TODO for development purpose
func queryOutLinks(ctx sdk.Context, req abci.RequestQuery, gk GraphKeeper) ([]byte, error) {
	var params types.QueryLinksParams

	err := codec.Cdc.UnmarshalJSON(req.Data, &params);	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	_, outLinks, _ := gk.GetAllLinks(ctx)
	cidNum, exist := gk.GetCidNumber(ctx, params.Cid); if exist != true {
		return nil, sdkerrors.Wrap(types.ErrCidNotFound, "")
	}

	links := outLinks[cidNum]
	data := make([]types.Cid, 0)

	for i, _ := range links {
		data = append(data, gk.GetCid(ctx, i))
	}

	res, err := codec.MarshalJSONIndent(codec.Cdc, types.ResultLinks{Cids: data})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

