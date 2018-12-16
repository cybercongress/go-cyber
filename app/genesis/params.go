package genesis

import (
	"github.com/cosmos/cosmos-sdk/x/stake/types"
	"github.com/cybercongress/cyberd/app/types/coin"
	"time"
)

const (
	// defaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	defaultUnbondingTime = 60 * 60 * 24 * 3 * time.Second
)

// DefaultParams returns a default set of parameters.
func DefaultStakeParams() types.Params {
	return types.Params{
		UnbondingTime: defaultUnbondingTime,
		MaxValidators: 146,
		BondDenom:     coin.CBD,
	}
}
