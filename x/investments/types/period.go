package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// NewPeriod returns a new vesting period
func NewPeriod(amount sdk.Coins, length int64) vestingtypes.Period {
	return vestingtypes.Period{Amount: amount, Length: length}
}
