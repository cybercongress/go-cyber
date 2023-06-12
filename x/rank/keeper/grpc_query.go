package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	querytypes "github.com/cybercongress/go-cyber/types/query"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
	"github.com/cybercongress/go-cyber/x/rank/types"
)

var _ types.QueryServer = &StateKeeper{}

func (bk StateKeeper) Params(goCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := bk.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (bk StateKeeper) Rank(goCtx context.Context, req *types.QueryRankRequest) (*types.QueryRankResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNum, exist := bk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.Particle)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, req.Particle)
	}

	rankValue := bk.index.GetRankValue(cidNum)
	return &types.QueryRankResponse{Rank: rankValue}, nil
}

func (bk *StateKeeper) Search(goCtx context.Context, req *types.QuerySearchRequest) (*types.QuerySearchResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNum, exist := bk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.Particle)); if exist != true {
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

	result := make([]types.RankedParticle, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, types.RankedParticle{Particle: string(bk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	return &types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}}, nil
}

func (bk *StateKeeper) Backlinks(goCtx context.Context, req *types.QuerySearchRequest) (*types.QuerySearchResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNum, exist := bk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.Particle)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, req.Particle)
	}

	page, limit := uint32(0), uint32(10)
	if req.Pagination != nil {
		page, limit = req.Pagination.Page, req.Pagination.PerPage
	}
	rankedCidNumbers, totalSize, err := bk.index.Backlinks(cidNum, page, limit)
	if err != nil {
		panic(err)
	}

	result := make([]types.RankedParticle, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, types.RankedParticle{Particle: string(bk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	return &types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}}, nil
}

func (bk *StateKeeper) Top(goCtx context.Context, req *querytypes.PageRequest) (*types.QuerySearchResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO check pagination
	page, limit := uint32(0), uint32(100)
	page, limit = req.Page, req.PerPage
	topRankedCidNumbers, totalSize, err := bk.index.Top(page, limit)
	if err != nil {
		panic(err)
	}

	result := make([]types.RankedParticle, 0, len(topRankedCidNumbers))
	for _, c := range topRankedCidNumbers {
		result = append(result, types.RankedParticle{Particle: string(bk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	return &types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}}, nil
}

func (bk StateKeeper) IsLinkExist(goCtx context.Context, req *types.QueryIsLinkExistRequest) (*types.QueryLinkExistResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	addr, err := sdk.AccAddressFromBech32(req.Address); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNumFrom, exist := bk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.From)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, req.From)
	}

	cidNumTo, exist := bk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.To)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, req.To)
	}

	var accountNum uint64
	account := bk.accountKeeper.GetAccount(ctx, addr)
	if account != nil {
		accountNum = account.GetAccountNumber()
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Invalid neuron address")
	}

	exists := bk.graphIndexedKeeper.IsLinkExist(graphtypes.CompactLink{
		uint64(cidNumFrom),
		uint64(cidNumTo),
		accountNum,
	})

	return &types.QueryLinkExistResponse{Exist: exists}, nil
}

func (bk StateKeeper) IsAnyLinkExist(goCtx context.Context, req *types.QueryIsAnyLinkExistRequest) (*types.QueryLinkExistResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNumFrom, exist := bk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.From)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, req.From)
	}

	cidNumTo, exist := bk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.To)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, req.To)
	}

	exists := bk.graphIndexedKeeper.IsAnyLinkExist(cidNumFrom, cidNumTo)

	return &types.QueryLinkExistResponse{Exist: exists}, nil
}

func (s *StateKeeper) ParticleNegentropy(goCtx context.Context, request *types.QueryNegentropyPartilceRequest) (*types.QueryNegentropyParticleResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNum, exist := s.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(request.Particle)); if exist != true {
		return nil, sdkerrors.Wrap(graphtypes.ErrCidNotFound, request.Particle)
	}

	entropyValue := s.GetEntropy(cidNum)
	return &types.QueryNegentropyParticleResponse{Entropy: entropyValue}, nil
}

func (s *StateKeeper) Negentropy(_ context.Context, _ *types.QueryNegentropyRequest) (*types.QueryNegentropyResponse, error) {
	negentropy := s.GetNegEntropy()
	return &types.QueryNegentropyResponse{Negentropy: negentropy}, nil
}

func (s *StateKeeper) Karma(goCtx context.Context, request *types.QueryKarmaRequest) (*types.QueryKarmaResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	addr, err := sdk.AccAddressFromBech32(request.Neuron); if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var accountNum uint64
	account := s.accountKeeper.GetAccount(ctx, addr)
	if account != nil {
		accountNum = account.GetAccountNumber()
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Invalid neuron address")
	}

	karma := s.GetKarma(accountNum)

	return &types.QueryKarmaResponse{Karma: karma}, nil
}
