package client

import (
	cbdlink "github.com/cybercongress/cyberd/x/link/types"
)

type CidsFilter map[cbdlink.Cid]map[cbdlink.Cid]struct{}

func (cf CidsFilter) Put(from cbdlink.Cid, to cbdlink.Cid) {

	cidLinks := cf[from]
	if cidLinks == nil {
		cidLinks = make(map[cbdlink.Cid]struct{})
	}
	cidLinks[to] = struct{}{}
	cf[from] = cidLinks
}

func (cf CidsFilter) Contains(from cbdlink.Cid, to cbdlink.Cid) bool {

	cidLinks := cf[from]
	if cidLinks == nil {
		return false
	}
	_, contains := cidLinks[to]
	return contains
}
