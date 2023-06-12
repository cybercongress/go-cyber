package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type CyberbankKeeper interface {
	OnCoinsTransfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress)
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
	GetModuleAddress(name string) sdk.AccAddress
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
}
