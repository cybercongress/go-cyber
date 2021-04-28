package keeper

import (
	"context"
	//"sort"

	//"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	//"github.com/cosmos/cosmos-sdk/types/query"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/status"
	"github.com/cybercongress/go-cyber/x/graph/types"
)

var _ types.QueryServer = GraphKeeper{}

// TODO developers endpoint, remove before release
func (gk GraphKeeper) OutLinks(goCtx context.Context, req *types.QueryLinksRequest) (*types.QueryLinksResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNum, exist := gk.GetCidNumber(ctx, types.Cid(req.Cid)); if exist != true {
		return nil, sdkerrors.Wrap(types.ErrCidNotFound, req.Cid)
	}
	// TODO add pagination
	_, outLinks, _ := gk.GetAllLinks(ctx)

	links := outLinks[cidNum]
	data := make([]string, 0)

	for i, _ := range links {
		data = append(data, string(gk.GetCid(ctx, i)))
	}

	return &types.QueryLinksResponse{Cids: data}, nil
}

// TODO developers endpoint, remove before release
func (gk GraphKeeper) InLinks(goCtx context.Context, req *types.QueryLinksRequest) (*types.QueryLinksResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	cidNum, exist := gk.GetCidNumber(ctx, types.Cid(req.Cid)); if exist != true {
		return nil, sdkerrors.Wrap(types.ErrCidNotFound, req.Cid)
	}
	// TODO add pagination
	inLinks, _, _ := gk.GetAllLinks(ctx)

	links := inLinks[cidNum]
	data := make([]string, 0)

	for i, _ := range links {
		data = append(data, string(gk.GetCid(ctx, i)))
	}

	return &types.QueryLinksResponse{Cids: data}, nil
}

func (gk GraphKeeper) LinksAmount(goCtx context.Context, _ *types.QueryLinksAmountRequest) (*types.QueryLinksAmountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	amount := gk.GetLinksCount(ctx)
	return &types.QueryLinksAmountResponse{Amount: amount}, nil
}

func (gk GraphKeeper) CidsAmount(goCtx context.Context, _ *types.QueryCidsAmountRequest) (*types.QueryCidsAmountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	amount := gk.GetCidsCount(ctx)
	return &types.QueryCidsAmountResponse{Amount: amount}, nil
}

func (gk GraphKeeper) GraphStats(goCtx context.Context, _ *types.QueryGraphStatsRequest) (*types.QueryGraphStatsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	links := gk.GetLinksCount(ctx)
	cids := gk.GetCidsCount(ctx)
	return &types.QueryGraphStatsResponse{links, cids}, nil
}