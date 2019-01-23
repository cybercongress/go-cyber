package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"log"
)

var noCoins = sdk.Coins{}

//__________________________________________________________________________________
// fee collection keeper used only for testing
type NoopFeeCollectionKeeper struct{}

func (fck NoopFeeCollectionKeeper) AddCollectedFees(_ sdk.Context, c sdk.Coins) sdk.Coins {
	log.Println(c.String())
	return noCoins
}

func (fck NoopFeeCollectionKeeper) GetCollectedFees(_ sdk.Context) sdk.Coins {
	return noCoins
}

func (fck NoopFeeCollectionKeeper) SetCollectedFees(_ sdk.Coins)     {}
func (fck NoopFeeCollectionKeeper) ClearCollectedFees(_ sdk.Context) {}
