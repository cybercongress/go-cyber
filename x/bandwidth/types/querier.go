package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryParameters		    = "params"
	QueryLoad			    = "load"
	QueryPrice			    = "price"
	QueryAccount		    = "account"
	QueryDesirableBandwidth = "desirable_bandwidth"
)

type QueryAccountBandwidthParams struct {
	Address   sdk.AccAddress
}

func NewQueryAccountBandwidthParams(addr sdk.AccAddress) QueryAccountBandwidthParams {
	return QueryAccountBandwidthParams{addr}
}