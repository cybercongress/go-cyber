package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CYB    = "boot"
	VOLT   = "millivolt"
	AMPERE = "milliampere"
	SCYB   = "hydrogen"
)

const (
	Deca int64 = 1
	Kilo       = Deca * 1000
	Mega       = Kilo * 1000
	Giga       = Mega * 1000
	Tera       = Giga * 1000
	Peta       = Tera * 1000
)

func NewCybCoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(CYB, amount)
}

func NewSCybCoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(SCYB, amount)
}

func NewVoltCoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(VOLT, amount)
}

func NewAmpereCoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(AMPERE, amount)
}
