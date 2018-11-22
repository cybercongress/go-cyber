package types

// map of map, where first key is cid, second key is account.String()
// second map is used as set for fast contains check
type Links map[CidNumber]CidLinks
type CidLinks map[CidNumber]map[AccountNumber]struct{}

type Cid string
type CidNumber uint64

func (links Links) Put(from CidNumber, to CidNumber, acc AccountNumber) {
	cidLinks := links[from]
	if cidLinks == nil {
		cidLinks = make(CidLinks)
	}
	users := cidLinks[to]
	if users == nil {
		users = make(map[AccountNumber]struct{})
	}
	users[acc] = struct{}{}
	cidLinks[to] = users
	links[from] = cidLinks
}
