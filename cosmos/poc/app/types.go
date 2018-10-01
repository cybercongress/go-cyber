package app

import "github.com/cosmos/cosmos-sdk/x/auth"

// map of map, where first key is cid, second key is account.String()
// second map is used as set for fast contains check
type CidLinks map[string]map[string]struct{}

func NewAccount() auth.Account {
	return &auth.BaseAccount{}
}
