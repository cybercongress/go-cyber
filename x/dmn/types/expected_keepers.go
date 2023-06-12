package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	//"github.com/cosmos/cosmos-sdk/x/params"
)

type BankKeeper interface {
	OnCoinsTransfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress)
	SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) error
}

type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
}