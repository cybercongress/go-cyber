package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/cyberd/app/coin"
	"github.com/cybercongress/cyberd/x/bandwidth/types"
)

func InitGenesis(
	ctx sdk.Context, bwKeeper AccountBandwidthKeeper, accKeeper auth.AccountKeeper, maxBandwidth types.MaxAccBandwidth,
) {

	accKeeper.IterateAccounts(ctx, func(account auth.Account) (stop bool) {

		addressStake := account.GetCoins().AmountOf(coin.CBD)
		bwKeeper.SetAccountBandwidth(
			ctx,
			types.NewGenesisAccBandwidth(account.GetAddress(), maxBandwidth(ctx, addressStake.Int64())),
		)

		return false
	})

}
