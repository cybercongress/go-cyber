package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)


type AccountKeeper interface {
	IterateAccounts(ctx sdk.Context, process func(authexported.Account) (stop bool))
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account // only used for simulation
	GetAllAccounts(ctx sdk.Context) (accounts []authexported.Account)
}

type EnergyKeeper interface {
	GetRoutedToEnergy(ctx sdk.Context, delegate sdk.AccAddress) sdk.Int
}

type StakingKeeper interface {
	GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []stakingtypes.Delegation
}

type SupplyKeeper interface {
	GetSupply(ctx sdk.Context) (supply exported.SupplyI)
}