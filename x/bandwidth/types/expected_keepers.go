package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AccountStakeProvider interface {
	GetAccountStakePercentage(ctx sdk.Context, address sdk.AccAddress) float64
}
