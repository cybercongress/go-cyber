package types

import "github.com/cosmos/cosmos-sdk/x/auth"

type AccNumber uint64

func NewCyberdAccount() auth.Account {
	return &auth.BaseAccount{}
}
