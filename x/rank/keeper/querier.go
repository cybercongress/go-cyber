package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	querytypes "github.com/cybercongress/go-cyber/v2/types/query"
	graphtypes "github.com/cybercongress/go-cyber/v2/x/graph/types"
	"github.com/cybercongress/go-cyber/v2/x/rank/types"

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
		case types.QueryKarma:
			return queryKarma(ctx, req, sk, legacyQuerierCdc)
		case types.QueryEntropy:
			return queryEntropy(ctx, req, sk, legacyQuerierCdc)
		case types.QueryNegentropy:
			return queryNegentropy(ctx, req, sk, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, _ abci.RequestQuery, sk StateKeeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := sk.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryRank(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryRankParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	cidNum, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.Cid))
	if !exist {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, params.Cid)
	}

	rankValue := sk.index.GetRankValue(cidNum)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryRankResponse{Rank: rankValue})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func querySearch(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QuerySearchParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	cidNum, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.Cid))
	if !exist {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}

	rankedCidNumbers, totalSize, err := sk.index.Search(cidNum, params.Page, params.PerPage)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	result := make([]types.RankedParticle, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, types.RankedParticle{Particle: string(sk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryBacklinks(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QuerySearchParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	cidNum, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.Cid))
	if !exist {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, params.Cid)
	}

	rankedCidNumbers, totalSize, err := sk.index.Backlinks(cidNum, params.Page, params.PerPage)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	result := make([]types.RankedParticle, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, types.RankedParticle{Particle: string(sk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryTop(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryTopParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	topRankedCidNumbers, totalSize, err := sk.index.Top(params.Page, params.PerPage)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	result := make([]types.RankedParticle, 0, len(topRankedCidNumbers))
	for _, c := range topRankedCidNumbers {
		result = append(result, types.RankedParticle{Particle: string(sk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryIsLinkExist(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryIsLinkExistParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	cidNumFrom, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.From))
	if !exist {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, params.From)
	}

	cidNumTo, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.To))
	if !exist {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, params.To)
	}

	var accountNum uint64
	account := sk.accountKeeper.GetAccount(ctx, params.Address)
	if account != nil {
		accountNum = account.GetAccountNumber()
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Invalid neuron address")
	}

	exists := sk.graphIndexedKeeper.IsLinkExist(graphtypes.CompactLink{
		From:    uint64(cidNumFrom),
		To:      uint64(cidNumTo),
		Account: accountNum,
	})

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, &types.QueryLinkExistResponse{Exist: exists})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryIsAnyLinkExist(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryIsAnyLinkExistParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	cidNumFrom, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.From))
	if !exist {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, params.From)
	}

	cidNumTo, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.To))
	if !exist {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, params.To)
	}

	exists := sk.graphIndexedKeeper.IsAnyLinkExist(cidNumFrom, cidNumTo)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, &types.QueryLinkExistResponse{Exist: exists})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryEntropy(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryEntropyParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	cidNum, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(params.Cid))
	if !exist {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, params.Cid)
	}

	entropy := sk.GetEntropy(cidNum)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, &types.QueryNegentropyParticleResponse{Entropy: entropy})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryNegentropy(_ sdk.Context, _ abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	negentropy := sk.GetNegEntropy()

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, &types.QueryNegentropyResponse{Negentropy: negentropy})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryKarma(ctx sdk.Context, req abci.RequestQuery, sk *StateKeeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryKarmaParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var accountNum uint64
	account := sk.accountKeeper.GetAccount(ctx, params.Address)
	if account != nil {
		accountNum = account.GetAccountNumber()
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Invalid neuron address")
	}

	karma := sk.GetKarma(accountNum)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, &types.QueryKarmaResponse{Karma: karma})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
