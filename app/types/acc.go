package types

import "github.com/cosmos/cosmos-sdk/x/auth"

func NewCyberdAccount() auth.Account {
	return &auth.BaseAccount{}
}
