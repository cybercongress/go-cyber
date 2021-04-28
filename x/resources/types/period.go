package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

func NewPeriod(amount sdk.Coins, length int64) vestingtypes.Period {
	return vestingtypes.Period{Amount: amount, Length: length}
}
