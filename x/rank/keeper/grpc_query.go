package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	querytypes "github.com/cybercongress/go-cyber/types/query"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
	"github.com/cybercongress/go-cyber/x/rank/types"
)

var _ types.QueryServer = &StateKeeper{}

func (bk StateKeeper) Params(goCtx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := bk.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (bk StateKeeper) Rank(goCtx context.Context, req *types.QueryRankRequest) (*types.QueryRankResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNum, exist := bk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.Cid)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}

	rankValue := bk.index.GetRankValue(cidNum)
	return &types.QueryRankResponse{Rank: rankValue}, nil
}

func (bk *StateKeeper) Search(goCtx context.Context, req *types.QuerySearchRequest) (*types.QuerySearchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNum, exist := bk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.Cid)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}

	page, limit := uint32(0), uint32(10)
	if req.Pagination != nil {
		page, limit = req.Pagination.Page, req.Pagination.PerPage
	}
	rankedCidNumbers, totalSize, err := bk.index.Search(cidNum, page, limit)
	if err != nil {
		panic(err)
	}

	result := make([]types.RankedCid, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, types.RankedCid{Cid: string(bk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	return &types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}}, nil
}

func (bk *StateKeeper) Backlinks(goCtx context.Context, req *types.QuerySearchRequest) (*types.QuerySearchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNum, exist := bk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.Cid)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}

	page, limit := uint32(0), uint32(10)
	if req.Pagination != nil {
		page, limit = req.Pagination.Page, req.Pagination.PerPage
	}
	rankedCidNumbers, totalSize, err := bk.index.Backlinks(cidNum, page, limit)
	if err != nil {
		panic(err)
	}

	result := make([]types.RankedCid, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, types.RankedCid{Cid: string(bk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	return &types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}}, nil
}

func (bk *StateKeeper) Top(goCtx context.Context, req *querytypes.PageRequest) (*types.QuerySearchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)


	page, limit := uint32(0), uint32(100)
	if req != nil {
		page, limit = req.Page, req.PerPage
	}
	topRankedCidNumbers, totalSize, err := bk.index.Top(page, limit)
	if err != nil {
		panic(err)
	}

	result := make([]types.RankedCid, 0, len(topRankedCidNumbers))
	for _, c := range topRankedCidNumbers {
		result = append(result, types.RankedCid{Cid: string(bk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	return &types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}}, nil
}

func (bk StateKeeper) IsLinkExist(goCtx context.Context, req *types.QueryIsLinkExistRequest) (*types.QueryLinkExistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNumFrom, exist := bk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.From)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}

	cidNumTo, exist := bk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.To)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}
	addr, err := sdk.AccAddressFromBech32(req.Address); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	accountNum := bk.accountKeeper.GetAccount(ctx, addr).GetAccountNumber()

	resp := uint32(0)
	exists := bk.graphIndexedKeeper.IsLinkExist(graphtypes.CompactLink{
		uint64(cidNumFrom),
		uint64(cidNumTo),
		accountNum,
	})
	if exists {
		resp = uint32(1)
	}

	return &types.QueryLinkExistResponse{Exist: resp}, nil
}

func (bk StateKeeper) IsAnyLinkExist(goCtx context.Context, req *types.QueryIsAnyLinkExistRequest) (*types.QueryLinkExistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNumFrom, exist := bk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.From)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}

	cidNumTo, exist := bk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.To)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, "")
	}

	resp := uint32(0)
	exists := bk.graphIndexedKeeper.IsAnyLinkExist(cidNumFrom, cidNumTo)
	if exists {
		resp = uint32(1)
	}

	return &types.QueryLinkExistResponse{Exist: resp}, nil
}

func (s *StateKeeper) Karmas(goCtx context.Context, request *types.QueryKarmasRequest) (*types.QueryKarmasResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	karmas := s.GetKarmas(ctx)

	return &types.QueryKarmasResponse{Karmas: karmas}, nil
}
