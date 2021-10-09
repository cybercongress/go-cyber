package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewRoute creates a new route object
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

// MustMarshalRoute returns the route bytes. Panics if fails
func MustMarshalRoute(cdc codec.BinaryCodec, route Route) []byte {
	return cdc.MustMarshal(&route)
}

// MustUnmarshalRoute return the unmarshaled route from bytes. Panics if fails.
func MustUnmarshalRoute(cdc codec.BinaryCodec, value []byte) Route {
	route, err := UnmarshalRoute(cdc, value)
	if err != nil {
		panic(err)
	}

	return route
}

func UnmarshalRoute(cdc codec.BinaryCodec, value []byte) (route Route, err error) {
	err = cdc.Unmarshal(value, &route)
	return route, err
}