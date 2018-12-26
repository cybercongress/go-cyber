package types

import . "github.com/cybercongress/cyberd/types"

// map of map, where first key is cid, second key is account.String()
// second map is used as set for fast contains check
type Links map[CidNumber]CidLinks
type CidLinks map[CidNumber]map[AccNumber]struct{}

type Cid string
type CidNumber uint64

func (links Links) Put(from CidNumber, to CidNumber, acc AccNumber) {
	cidLinks := links[from]
	if cidLinks == nil {
		cidLinks = make(CidLinks)
	}
	users := cidLinks[to]
	if users == nil {
		users = make(map[AccNumber]struct{})
	}
	users[acc] = struct{}{}
	cidLinks[to] = users
	links[from] = cidLinks
}

func (links Links) Copy() Links {

	linksCopy := make(Links, len(links))

	for from := range links {
		fromLinks := make(CidLinks, len(links[from]))
		for to := range links[from] {
			users := make(map[AccNumber]struct{}, len(links[from][to]))
			for u := range links[from][to] {
				users[u] = struct{}{}
			}
			fromLinks[to] = users
		}
		linksCopy[from] = fromLinks
	}
	return linksCopy
}
