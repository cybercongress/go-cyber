package rpc

import "github.com/cybercongress/cyberd/cosmos/poc/app/storage"

type ResultSearch struct {
	Cids       []storage.RankedCid `json:"cids"`
	TotalCount int                 `json:"totalCount"`
}

func Search(cid string, page, perPage int) (*ResultSearch, error) {
	links, totalSize, err := cyberdApp.Search(cid, page, perPage)
	return &ResultSearch{links, totalSize}, err
}
