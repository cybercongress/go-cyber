package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewRoute(src sdk.AccAddress, dst sdk.AccAddress, alias string, value sdk.Coins) Route {
	return Route{
		Source:      src.String(),
		Destination: dst.String(),
		Alias: 		 alias,
		Value:       value,
	}
}

type Routes []Route

func NewValue(coins sdk.Coins) Value {
	return Value{coins}
}