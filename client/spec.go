package client

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cbd "github.com/cybercongress/cyberd/app/types"
	bdwth "github.com/cybercongress/cyberd/x/bandwidth/types"
)

type Link struct {
	from    cbd.Cid
	to      cbd.Cid
	address sdk.AccAddress
}

type CyberdClient interface {
	// Cyberd Client Specification

	// returns current connected node chain id
	GetChainId() (string, error)

	// get current bandwidth credits price
	// price 1 is price for situation, when all users use all their bandwidth (all blocks are filled for 100%)
	// if price < 1, that means blocks filled partially, thus allow more active users to do more transactions
	// if price > 1, that means network is under high load.
	GetCurrentBandwidthCreditPrice() (float64, error)

	// returns account for given address
	GetAccount(address sdk.AccAddress) (auth.Account, error)

	// returns account bandwidth information for given account
	GetAccountBandwidth(address sdk.AccAddress) (bdwth.AccountBandwidth, error)

	// links two cids for given user
	// this method also should check, either cids are correct cids and given user is msg signer
	// do not wait till tx will be mined, just returns results from tx mempool check
	SubmitLinkAsync(link Link) error

	// see `SubmitLinkAsync`. Links will be submitted as single tx with multiple msges.
	SubmitLinksAsync(links []Link) error
}
