package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *IndexedKeeper) InitGenesis(ctx sdk.Context) {
	for _, account := range k.accountKeeper.GetAllAccounts(ctx) {
		k.InitializeStakeAmpere(
			account.GetAccountNumber(),
			uint64(k.Proxy.GetAccountTotalStakeAmper(ctx, account.GetAddress())),
		)
	}
}
