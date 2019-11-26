package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/staking"
	sdksupply "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	GetAllAccounts(ctx sdk.Context) []exported.Account
	SetAccount(ctx sdk.Context, acc exported.Account)
	IterateAccounts(ctx sdk.Context, process func(exported.Account) bool)
}

type StakingKeeper interface {
	GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []staking.Delegation
}

type SupplyKeeper interface {
	GetSupply(ctx sdk.Context) (supply sdksupply.SupplyI)
}
