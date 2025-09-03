package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AccountStakeProvider interface {
	GetAccountStakeVolt(ctx sdk.Context, addr sdk.AccAddress) int64
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}
