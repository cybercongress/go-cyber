package app

import "github.com/cosmos/cosmos-sdk/x/auth"

func NewAccount() auth.Account {
	return &auth.BaseAccount{}
}
