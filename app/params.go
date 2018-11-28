package app

import (
	"github.com/cosmos/cosmos-sdk/x/stake/types"
	"time"
)

const (
	// defaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	defaultUnbondingTime time.Duration = 60 * 60 * 24 * 3 * time.Second
)

// DefaultParams returns a default set of parameters.
func DefaultStakeParams() types.Params {
	return types.Params{
		UnbondingTime: defaultUnbondingTime,
		MaxValidators: 100,
		BondDenom:     "CBD",
	}
}
