package rpc

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	bdwth "github.com/cybercongress/cyberd/x/bandwidth/types"
)

type ResultAccount struct {
	Account auth.BaseAccount `json:"account"`
}

func Account(address string) (*ResultAccount, error) {
	accAddress, err := types.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	account := cyberdApp.Account(accAddress)

	if account == nil {
		return nil, errors.New("account not found")
	}

	return &ResultAccount{
		auth.BaseAccount{
			Address:       account.GetAddress(),
			Coins:         account.GetCoins(),
			PubKey:        account.GetPubKey(),
			AccountNumber: account.GetAccountNumber(),
			Sequence:      account.GetSequence()},
	}, nil
}

func AccountBandwidth(address string) (*bdwth.AcсBandwidth, error) {

	accAddress, err := types.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	accBdwth := cyberdApp.AccountBandwidth(accAddress)
	return &accBdwth, nil
}
