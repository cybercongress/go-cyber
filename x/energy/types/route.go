package types

import (
	"fmt"
	"strings"

	//"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Route struct {
	Source 		sdk.AccAddress `json:"source" yaml:"source"`
	Destination sdk.AccAddress `json:"destination" yaml:"destination"`
	Alias		string		   `json:"alias" yaml:"alias"`
	Amount    	sdk.Int        `json:"amount" yaml:"amount"`
}

func NewRoute(src sdk.AccAddress, dst sdk.AccAddress, alias string, amount sdk.Int) Route {
	return Route{
		Source:      src,
		Destination: dst,
		Alias: 		 alias,
		Amount:      amount,
	}
}

func MustMarshalRoute(cdc *codec.Codec, route Route) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(route)
}

func MustUnmarshalRoute(cdc *codec.Codec, value []byte) Route {
	route, err := UnmarshalRoute(cdc, value)
	if err != nil {
		panic(err)
	}
	return route
}

func UnmarshalRoute(cdc *codec.Codec, value []byte) (route Route, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &route)
	return route, err
}

func (r Route) GetSource()      sdk.AccAddress { return r.Source }
func (r Route) GetDestination() sdk.AccAddress { return r.Destination }
func (r Route) GetAlias() 	    string         { return r.Alias }
func (r Route) GetAmount() 	    sdk.Int        { return r.Amount }

func (r Route) String() string {
	return fmt.Sprintf(`Delegation:
  Source: 	   %s
  Destination: %s
  Alias:       %s
  Amount:      %s`,
  r.Source, r.Destination, r.Alias, r.Amount,)
}

type Routes []Route

func (d Routes) String() (out string) {
	for _, route := range d {
		out += route.String() + "\n"
	}
	return strings.TrimSpace(out)
}