package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// "github.com/cosmos/cosmos-sdk/x/auth/exported"
	// "github.com/cosmos/cosmos-sdk/x/params"
	// bankexported "github.com/cosmos/cosmos-sdk/x/bank/exported"
	// authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	// "github.com/cybercongress/go-cyber/v2/x/bank/exported"
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
