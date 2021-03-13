package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankexported "github.com/cosmos/cosmos-sdk/x/bank/exported"
)


type AccountKeeper interface {
	IterateAccounts(ctx sdk.Context, process func(i authtypes.AccountI) (stop bool))
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetAllAccounts(ctx sdk.Context) (accounts []authtypes.AccountI)
}

type EnergyKeeper interface {
	GetRoutedToEnergy(ctx sdk.Context, delegate sdk.AccAddress) sdk.Coins
}

//type StakingKeeper interface {
//	GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []stakingtypes.Delegation
//}

type BankKeeper interface {
	GetSupply(ctx sdk.Context) bankexported.SupplyI
}