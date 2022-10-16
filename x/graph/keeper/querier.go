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
		case types.QueryGraphStats:
			return queryGraphStats(ctx, req, gk, legacyQuerierCdc)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryGraphStats(ctx sdk.Context, _ abci.RequestQuery, gk GraphKeeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	links := gk.GetLinksCount(ctx)
	cids := gk.GetCidsCount(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryGraphStatsResponse{links, cids})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
