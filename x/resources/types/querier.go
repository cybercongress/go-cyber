package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryParams 	      = "params"
	QueryInvestmintAmount = "investmint_amount"
)

type QueryInvestmintAmountParams struct {
	Amount sdk.Coin
	Resource string
	Length uint64
}

func NewQueryInvestmintAmountParams(amount sdk.Coin, resource string, length uint64) QueryInvestmintAmountParams {
	return QueryInvestmintAmountParams{amount, resource, length}
}