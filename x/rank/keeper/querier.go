package keeper

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	querytypes "github.com/cybercongress/go-cyber/types/query"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
	"github.com/cybercongress/go-cyber/x/rank/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewQuerier(sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParams(ctx, req, *sk, legacyQuerierCdc)
		case types.QueryRank:
			return queryRank(ctx, req, sk, legacyQuerierCdc)
		case types.QuerySearch:
			return querySearch(ctx, req, sk, legacyQuerierCdc)
		case types.QueryBacklinks:
			return queryBacklinks(ctx, req, sk, legacyQuerierCdc)
		case types.QueryTop:
			return queryTop(ctx, req, sk, legacyQuerierCdc)
		case types.QueryIsLinkExist:
			return queryIsLinkExist(ctx, req, sk, legacyQuerierCdc)
		case types.QueryIsAnyLinkExist:
			return queryIsAnyLinkExist(ctx, req, sk, legacyQuerierCdc)
		case types.QueryKarmas:
			return queryKarmas(ctx, req, sk, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, _ abci.RequestQuery, sk StateKeeper, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	params := sk.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryRank(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	var params types.QueryRankRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	cidNum, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.Cid)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}

	rankValue := sk.index.GetRankValue(cidNum)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryRankResponse{Rank: rankValue})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func querySearch(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	var params types.QuerySearchRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	cidNum, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.Cid)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}

	rankedCidNumbers, totalSize, err := sk.index.Search(cidNum, params.Pagination.Page, params.Pagination.PerPage)
	if err != nil {
		panic(err)
	}

	result := make([]types.RankedCid, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, types.RankedCid{Cid: string(sk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryBacklinks(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	var params types.QuerySearchRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	cidNum, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.Cid)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}

	rankedCidNumbers, totalSize, err := sk.index.Backlinks(cidNum, params.Pagination.Page, params.Pagination.PerPage)
	if err != nil {
		panic(err)
	}

	result := make([]types.RankedCid, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, types.RankedCid{Cid: string(sk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryTop(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	var params querytypes.PageRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	topRankedCidNumbers, totalSize, err := sk.index.Top(params.Page, params.PerPage)
	if err != nil {
		panic(err)
	}

	result := make([]types.RankedCid, 0, len(topRankedCidNumbers))
	for _, c := range topRankedCidNumbers {
		result = append(result, types.RankedCid{Cid: string(sk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryIsLinkExist(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	var params types.QueryIsLinkExistRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	cidNumFrom, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.From)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}

	cidNumTo, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.To)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}
	addr, err := sdk.AccAddressFromBech32(params.Address); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	accountNum := sk.accountKeeper.GetAccount(ctx, addr).GetAccountNumber()

	resp := uint32(0)
	exists := sk.graphIndexedKeeper.IsLinkExist(graphtypes.CompactLink{
		uint64(cidNumFrom),
		uint64(cidNumTo),
		accountNum,
	})
	if exists {
		resp = uint32(1)
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, &types.QueryLinkExistResponse{Exist: resp})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryIsAnyLinkExist(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	var params types.QueryIsAnyLinkExistRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	cidNumFrom, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.From)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}

	cidNumTo, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.To)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}

	resp := uint32(0)
	exists := sk.graphIndexedKeeper.IsAnyLinkExist(cidNumFrom, cidNumTo)
	if exists {
		resp = uint32(1)
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, &types.QueryLinkExistResponse{Exist: resp})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryKarmas(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino,) ([]byte, error) {
	var params types.QueryKarmasRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	karmas := sk.GetKarmas(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, &types.QueryKarmasResponse{Karmas: karmas})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

