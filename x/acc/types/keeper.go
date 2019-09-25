package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto"
)

type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) auth.Account
	NewAccount(ctx sdk.Context, account auth.Account) auth.Account
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.Account
	GetAllAccounts(ctx sdk.Context) []auth.Account
	SetAccount(ctx sdk.Context, account auth.Account)
	RemoveAccount(ctx sdk.Context, account auth.Account)
	IterateAccounts(ctx sdk.Context, process func(auth.Account) (stop bool))
	GetPubKey(ctx sdk.Context, addr sdk.AccAddress) (crypto.PubKey, sdk.Error)
	GetSequence(ctx sdk.Context, addr sdk.AccAddress) (uint64, sdk.Error)
	GetNextAccountNumber(ctx sdk.Context) uint64

	SetParams(ctx sdk.Context, params auth.Params)
	GetParams(ctx sdk.Context) auth.Params
}

type AccountIndexKeeper interface {
	AccountKeeper

	GetAccountKeeper() auth.AccountKeeper // TODO: Remove when update to the new version SDK

	AddToIndex(account auth.Account)
	RefreshIndex(ctx sdk.Context)
	GetAccountAddress(number AccNumber) (addr sdk.AccAddress, ok bool)
	GetAccountAddresses(numbers []AccNumber) []sdk.AccAddress
}
