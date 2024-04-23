package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/ipfs/go-cid"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	graphtypes "github.com/cybercongress/go-cyber/v4/x/graph/types"
	"github.com/cybercongress/go-cyber/v4/x/rank/types"
	querytypes "github.com/cybercongress/go-cyber/v4/x/rank/types"
)

var _ types.QueryServer = &StateKeeper{}

func (sk StateKeeper) Params(goCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := sk.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (sk StateKeeper) Rank(goCtx context.Context, req *types.QueryRankRequest) (*types.QueryRankResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	particle, err := cid.Decode(req.Particle)
	if err != nil {
		return nil, graphtypes.ErrInvalidParticle
	}

	if particle.Version() != 0 {
		return nil, graphtypes.ErrCidVersion
	}

	cidNum, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.Particle))
	if !exist {
		return nil, errorsmod.Wrap(graphtypes.ErrCidNotFound, req.Particle)
	}

	// rankValue := sk.index.GetRankValue(cidNum) // TODO it was the bug, test bindings
	rankValue := sk.networkCidRank.RankValues[cidNum]
	return &types.QueryRankResponse{Rank: rankValue}, nil
}

func (sk *StateKeeper) Search(goCtx context.Context, req *types.QuerySearchRequest) (*types.QuerySearchResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNum, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.Particle))
	if !exist {
		return nil, errorsmod.Wrap(graphtypes.ErrCidNotFound, "")
	}

	page, limit := uint32(0), uint32(10)
	if req.Pagination != nil {
		page, limit = req.Pagination.Page, req.Pagination.PerPage
	}
	rankedCidNumbers, totalSize, err := sk.index.Search(cidNum, page, limit)
	if err != nil {
		panic(err)
	}

	result := make([]types.RankedParticle, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, types.RankedParticle{Particle: string(sk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	return &types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}}, nil
}

func (sk *StateKeeper) Backlinks(goCtx context.Context, req *types.QuerySearchRequest) (*types.QuerySearchResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNum, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.Particle))
	if !exist {
		return nil, errorsmod.Wrap(graphtypes.ErrCidNotFound, req.Particle)
	}

	page, limit := uint32(0), uint32(10)
	if req.Pagination != nil {
		page, limit = req.Pagination.Page, req.Pagination.PerPage
	}
	rankedCidNumbers, totalSize, err := sk.index.Backlinks(cidNum, page, limit)
	if err != nil {
		panic(err)
	}

	result := make([]types.RankedParticle, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, types.RankedParticle{Particle: string(sk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	return &types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}}, nil
}

func (sk *StateKeeper) Top(goCtx context.Context, req *querytypes.PageRequest) (*types.QuerySearchResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if req.PerPage > uint32(1000) {
		return nil, sdkerrors.ErrInvalidRequest
	}
	page, limit := req.Page, req.PerPage
	topRankedCidNumbers, totalSize, err := sk.index.Top(page, limit)
	if err != nil {
		panic(err)
	}

	result := make([]types.RankedParticle, 0, len(topRankedCidNumbers))
	for _, c := range topRankedCidNumbers {
		result = append(result, types.RankedParticle{Particle: string(sk.graphKeeper.GetCid(ctx, c.GetNumber())), Rank: c.GetRank()})
	}

	return &types.QuerySearchResponse{Result: result, Pagination: &querytypes.PageResponse{Total: totalSize}}, nil
}

func (sk StateKeeper) IsLinkExist(goCtx context.Context, req *types.QueryIsLinkExistRequest) (*types.QueryLinkExistResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNumFrom, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.From))
	if !exist {
		return nil, errorsmod.Wrap(graphtypes.ErrCidNotFound, req.From)
	}

	cidNumTo, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.To))
	if !exist {
		return nil, errorsmod.Wrap(graphtypes.ErrCidNotFound, req.To)
	}

	var accountNum uint64
	account := sk.accountKeeper.GetAccount(ctx, addr)
	if account != nil {
		accountNum = account.GetAccountNumber()
	} else {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "Invalid neuron address")
	}

	exists := sk.graphIndexedKeeper.IsLinkExist(graphtypes.CompactLink{
		From:    uint64(cidNumFrom),
		To:      uint64(cidNumTo),
		Account: accountNum,
	})

	return &types.QueryLinkExistResponse{Exist: exists}, nil
}

func (sk StateKeeper) IsAnyLinkExist(goCtx context.Context, req *types.QueryIsAnyLinkExistRequest) (*types.QueryLinkExistResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNumFrom, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.From))
	if !exist {
		return nil, errorsmod.Wrap(graphtypes.ErrCidNotFound, req.From)
	}

	cidNumTo, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(req.To))
	if !exist {
		return nil, errorsmod.Wrap(graphtypes.ErrCidNotFound, req.To)
	}

	exists := sk.graphIndexedKeeper.IsAnyLinkExist(cidNumFrom, cidNumTo)

	return &types.QueryLinkExistResponse{Exist: exists}, nil
}

func (sk *StateKeeper) ParticleNegentropy(goCtx context.Context, request *types.QueryNegentropyPartilceRequest) (*types.QueryNegentropyParticleResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNum, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(request.Particle))
	if !exist {
		return nil, errorsmod.Wrap(graphtypes.ErrCidNotFound, request.Particle)
	}

	entropyValue := sk.GetEntropy(cidNum)
	return &types.QueryNegentropyParticleResponse{Entropy: entropyValue}, nil
}

func (sk *StateKeeper) Negentropy(_ context.Context, _ *types.QueryNegentropyRequest) (*types.QueryNegentropyResponse, error) {
	negentropy := sk.GetNegEntropy()
	return &types.QueryNegentropyResponse{Negentropy: negentropy}, nil
}

func (sk *StateKeeper) Karma(goCtx context.Context, request *types.QueryKarmaRequest) (*types.QueryKarmaResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	addr, err := sdk.AccAddressFromBech32(request.Neuron)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var accountNum uint64
	account := sk.accountKeeper.GetAccount(ctx, addr)
	if account != nil {
		accountNum = account.GetAccountNumber()
	} else {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "Invalid neuron address")
	}

	karma := sk.GetKarma(accountNum)

	return &types.QueryKarmaResponse{Karma: karma}, nil
}
