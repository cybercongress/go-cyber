package rpc

import (
	"github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/link"

	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

type ResultAccountLinks struct {
	Links	   []link.Link 	   `json:"links"`
	TotalCount int             `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"perPage"`
}

func AccountLinks(ctx *rpctypes.Context, address string, page, perPage int) (*ResultAccountLinks, error) {
	if perPage == 0 {
		perPage = 100
	}

	accAddress, err := types.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	links, totalSize, err := cyberdApp.AccountLinks(accAddress, page, perPage)
	return &ResultAccountLinks{links, totalSize, page, perPage}, err
}


