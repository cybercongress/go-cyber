package types

type CidsFilter map[Cid]map[Cid]struct{}

func (cf CidsFilter) Put(from Cid, to Cid) {

	cidLinks := cf[from]
	if cidLinks == nil {
		cidLinks = make(map[Cid]struct{})
	}
	cidLinks[to] = struct{}{}
	cf[from] = cidLinks
}

func (cf CidsFilter) Contains(from Cid, to Cid) bool {

	cidLinks := cf[from]
	if cidLinks == nil {
		return false
	}
	_, contains := cidLinks[to]
	return contains
}
