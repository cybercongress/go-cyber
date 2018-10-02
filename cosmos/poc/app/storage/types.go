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
