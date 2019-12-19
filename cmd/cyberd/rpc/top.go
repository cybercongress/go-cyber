package rpc

import (
	"github.com/tendermint/tendermint/rpc/lib/types"
)

func Top(ctx *rpctypes.Context, page, perPage int) (*ResultSearch, error) {
	if perPage == 0 {
		perPage = 100
	}
	cids, totalSize, err := cyberdApp.Top(page, perPage)
	return &ResultSearch{cids, totalSize, page, perPage}, err
}
