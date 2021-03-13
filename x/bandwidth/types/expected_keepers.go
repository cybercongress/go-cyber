package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AccountStakeProvider interface {
	GetAccountStakePercentageVolt(ctx sdk.Context, address sdk.AccAddress) float64
}
