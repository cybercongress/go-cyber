package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type AccountKeeper interface {
	IterateAccounts(ctx sdk.Context, process func(i authtypes.AccountI) (stop bool))
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetAllAccounts(ctx sdk.Context) (accounts []authtypes.AccountI)
}

type EnergyKeeper interface {
	GetRoutedToEnergy(ctx sdk.Context, delegate sdk.AccAddress) sdk.Coins
}
