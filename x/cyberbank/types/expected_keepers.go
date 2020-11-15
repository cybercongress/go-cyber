package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)


type AccountKeeper interface {
	IterateAccounts(ctx sdk.Context, process func(authexported.Account) (stop bool))
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account // only used for simulation
	GetAllAccounts(ctx sdk.Context) (accounts []authexported.Account)
}

//type BankKeeper interface {
//	ValidateBalance(ctx sdk.Context, addr sdk.AccAddress) error
//	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
//	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
//	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
//	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
//	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
//	IterateAccountBalances(ctx sdk.Context, addr sdk.AccAddress, cb func(coin sdk.Coin) (stop bool))
//	IterateAllBalances(ctx sdk.Context, cb func(address sdk.AccAddress, coin sdk.Coin) (stop bool))
//	InputOutputCoins(ctx sdk.Context, inputs []bank.Input, outputs []bank.Output) error
//	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
//	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
//	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
//	SetBalance(ctx sdk.Context, addr sdk.AccAddress, balance sdk.Coin) error
//	SetBalances(ctx sdk.Context, addr sdk.AccAddress, balances sdk.Coins) error
//	GetSendEnabled(ctx sdk.Context) bool
//	SetSendEnabled(ctx sdk.Context, enabled bool)
//	BlacklistedAddr(addr sdk.AccAddress) bool
//	GetSupply(ctx sdk.Context) bankexported.SupplyI
//	SetSupply(ctx sdk.Context, supply bankexported.SupplyI)
//	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
//	SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt sdk.Coins) error
//	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
//	DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
//	UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
//	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
//	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
//	DelegateCoins(ctx sdk.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) error
//	UndelegateCoins(ctx sdk.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) error
//
//	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
//	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool
//	SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error
//}

type StakingKeeper interface {
	GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []stakingtypes.Delegation
}