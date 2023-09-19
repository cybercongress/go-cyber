package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (ik *IndexedKeeper) InitGenesis(ctx sdk.Context) {
	for _, account := range ik.accountKeeper.GetAllAccounts(ctx) {
		ik.InitializeStakeAmpere(
			account.GetAccountNumber(),
			uint64(ik.GetAccountTotalStakeAmper(ctx, account.GetAddress())),
		)
	}
}
