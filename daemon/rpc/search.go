package rpc

import "github.com/cybercongress/cyberd/app/storage"

type ResultSearch struct {
	Cids       []storage.RankedCid `json:"cids"`
	TotalCount int                 `json:"total"`
	Page       int                 `json:"page"`
	PerPage    int                 `json:"perPage"`
}

func Search(cid string, page, perPage int) (*ResultSearch, error) {
	if perPage == 0 {
		perPage = 100
	}
	links, totalSize, err := cyberdApp.Search(cid, page, perPage)
	return &ResultSearch{links, totalSize, page, perPage}, err
}
