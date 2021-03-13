package exported

import sdk "github.com/cosmos/cosmos-sdk/types"

type EnergyKeeper interface {
	GetRoutedToEnergy(ctx sdk.Context, delegate sdk.AccAddress) sdk.Coins
}
