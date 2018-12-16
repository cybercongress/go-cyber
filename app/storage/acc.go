package storage

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cbd "github.com/cybercongress/cyberd/app/types"
	"github.com/cybercongress/cyberd/app/types/coin"
)

// returns all added cids
func GetAllAccountsStakes(ctx sdk.Context, am auth.AccountKeeper) map[cbd.AccountNumber]uint64 {

	stakes := make(map[cbd.AccountNumber]uint64)

	collect := func(acc auth.Account) bool {
		balance := acc.GetCoins().AmountOf(coin.CBD).Int64()
		stakes[cbd.AccountNumber(acc.GetAccountNumber())] = uint64(balance)
		return false
	}

	am.IterateAccounts(ctx, collect)
	return stakes
}
