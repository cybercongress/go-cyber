package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/app/types/coin"
)

// Inflate every block
func BeginBlocker(ctx sdk.Context, m Minter) {
	mintedCoin := sdk.NewInt64Coin(coin.CBD, int64(m.blockReward))
	m.fck.AddCollectedFees(ctx, sdk.Coins{mintedCoin})
	m.stakeKeeper.InflateSupply(ctx, sdk.NewDecFromInt(mintedCoin.Amount))
}
