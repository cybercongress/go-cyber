package types

import (
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// GetTotalVestingPeriodLength returns the summed length of all vesting periods.
func GetTotalVestingPeriodLength(periods vestingtypes.Periods) int64 {
	length := int64(0)
	for _, period := range periods {
		length += period.Length
	}
	return length
}
