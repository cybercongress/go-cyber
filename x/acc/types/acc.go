package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

type AccNumber uint64

func NewCyberdAccount() auth.Account {
	return &auth.BaseAccount{}
}

// cbd1qqqqqqqqqqqqqqqqqqqqqqqqqqqqph4dvtfakz    0.1.0
// cyber1qqqqqqqqqqqqqqqqqqqqqqqqqqqqph4djv52fh  0.1.1
func GetBurnAddress() (addr sdk.AccAddress) {
	hex := "000000000000000000000000000000000000dead"
	addr, _ = sdk.AccAddressFromHex(hex)
	return
}
