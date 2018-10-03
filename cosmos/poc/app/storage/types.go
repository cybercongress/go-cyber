package storage

// map of map, where first key is cid, second key is account.String()
// second map is used as set for fast contains check
type CidLinks map[CidNumber]map[AccountNumber]struct{}

//todo should be uint
type AccountNumber string

type Cid string

type CidNumber uint64

type LinkedCids struct {
	FromCid CidNumber
	ToCid   CidNumber
	Creator AccountNumber
}

// All addition storages introduced by cyberd
type CyberdPersistentStorages struct {
	CidIndex CidIndexStorage
	Links    LinksStorage
	Rank     RankStorage
}

type CidsLinks map[CidNumber]CidLinks

func (mp CidsLinks) Put(c1 CidNumber, c2 CidNumber, acc AccountNumber) {
	cidLinks := mp[c1]
	if cidLinks == nil {
		cidLinks = make(CidLinks)
	}
	users := cidLinks[c2]
	if users == nil {
		users = make(map[AccountNumber]struct{})
	}
	users[acc] = struct{}{}
	cidLinks[c2] = users
	mp[c1] = cidLinks
}
