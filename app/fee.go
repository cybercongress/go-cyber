package app

import sdk "github.com/cosmos/cosmos-sdk/types"

var noCoins = sdk.Coins{}

//__________________________________________________________________________________
// fee collection keeper used only for testing
type NoopFeeCollectionKeeper struct{}

func (fck NoopFeeCollectionKeeper) AddCollectedFees(sdk.Context, sdk.Coins) sdk.Coins {
	return noCoins
}

func (fck NoopFeeCollectionKeeper) GetCollectedFees(_ sdk.Context) sdk.Coins {
	return noCoins
}

func (fck NoopFeeCollectionKeeper) SetCollectedFees(_ sdk.Coins)     {}
func (fck NoopFeeCollectionKeeper) ClearCollectedFees(_ sdk.Context) {}
