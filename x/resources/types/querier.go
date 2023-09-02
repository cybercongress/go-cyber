package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryParams     = "params"
	QueryInvestmint = "investmint"
)

type QueryInvestmintParams struct {
	Amount   sdk.Coin
	Resource string
	Length   uint64
}

func NewQueryInvestmintParams(amount sdk.Coin, resource string, length uint64) QueryInvestmintParams {
	return QueryInvestmintParams{amount, resource, length}
}
