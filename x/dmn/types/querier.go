package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryParams   	   = "params"
	QueryThought       = "thought"
	QueryThoughtStats  = "thought_stats"
	QueryThoughts      = "thoughts"
	QueryThoughtsStats = "thoughts_stats"
)

type QueryThoughtParams struct {
	Program sdk.AccAddress
	Name    string
}

func NewQueryThoughtParams(program sdk.AccAddress, name string) QueryThoughtParams {
	return QueryThoughtParams{program, name}
}