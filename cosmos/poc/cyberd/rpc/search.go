package rpc

import (
	cbd "github.com/cybercongress/cyberd/cosmos/poc/app/types"
)

type ResultSearch struct {
	Cids       []cbd.CidWithRank `json:"cids"`
	TotalCount int               `json:"total"`
	Page       int               `json:"page"`
	PerPage    int               `json:"perPage"`
}

func Search(cid string, page, perPage int) (*ResultSearch, error) {
	if perPage == 0 {
		perPage = 100
	}
	links, totalSize, err := cyberdApp.Search(cbd.Cid(cid), page, perPage)
	return &ResultSearch{links, totalSize, page, perPage}, err
}
