package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/types/coin"
)

// Inflate every block
func BeginBlocker(ctx sdk.Context, m Minter) {
	mintedCoins := sdk.NewInt64Coin(coin.CYB, m.GetParams(ctx).TokensPerBlock)
	m.fck.AddCollectedFees(ctx, sdk.Coins{mintedCoins})
	m.stakeKeeper.InflateSupply(ctx, mintedCoins.Amount)
}
