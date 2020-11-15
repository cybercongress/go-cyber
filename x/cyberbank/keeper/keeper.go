package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cybercongress/go-cyber/x/cyberbank/types"
)

type IndexedKeeper struct {
	*Proxy
	accountKeeper     types.AccountKeeper

	userTotalStake 	  map[uint64]uint64
	userNewTotalStake map[uint64]uint64
	accountToUpdate   []sdk.AccAddress
}

func NewIndexedKeeper(pbk *Proxy, ak types.AccountKeeper) *IndexedKeeper {
	indexedKeeper := &IndexedKeeper{
		Proxy: pbk,
		accountKeeper: ak,
		accountToUpdate: make([]sdk.AccAddress, 0),
	}
	hook := func(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress) {
		if from != nil {
			indexedKeeper.accountToUpdate = append(indexedKeeper.accountToUpdate, from)
		}
		if to != nil {
			indexedKeeper.accountToUpdate = append(indexedKeeper.accountToUpdate, to)
		}
	}
	pbk.AddHook(hook)
	return indexedKeeper
}

func (k IndexedKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *IndexedKeeper) Load(rankCtx sdk.Context, freshCtx sdk.Context) {
	k.userTotalStake = make(map[uint64]uint64)
	k.accountKeeper.IterateAccounts(rankCtx, k.getCollectFunc(rankCtx, k.userTotalStake))

	k.userNewTotalStake = make(map[uint64]uint64)
	k.accountKeeper.IterateAccounts(freshCtx, k.getCollectFunc(freshCtx, k.userNewTotalStake))
}

func (k *IndexedKeeper) getCollectFunc(ctx sdk.Context, userStake map[uint64]uint64) func(account exported.Account) bool {
	return func(account exported.Account) bool {
		balance := k.Proxy.GetAccountTotalStake(ctx, account.GetAddress())
		userStake[account.GetAccountNumber()] = uint64(balance)
		return false
	}
}

func  (k *IndexedKeeper) UpdateStake(account uint64, stake uint64) {
	k.userNewTotalStake[account] = stake
}

func (k *IndexedKeeper) GetTotalStakes() map[uint64]uint64 {
	return k.userTotalStake
}

func (k *IndexedKeeper) DetectUsersStakeChange(ctx sdk.Context) bool {
	stakeChanged := false
	for o, n := range k.userNewTotalStake {
		if k.userTotalStake[o] != n {
			stakeChanged = true
			k.userTotalStake[o] = n
		}
	}
	return stakeChanged
}

func (k *IndexedKeeper) InitGenesis(ctx sdk.Context) {
	for _, account := range k.accountKeeper.GetAllAccounts(ctx) {
		k.UpdateStake(
			account.GetAccountNumber(),
			uint64(k.Proxy.GetAccountTotalStake(ctx, account.GetAddress())),
		)
	}
}

func (k *IndexedKeeper) EndBlocker(ctx sdk.Context) {
	for _, addr := range k.accountToUpdate {
		stake := k.Proxy.GetAccountTotalStake(ctx, addr)
		accountNumber := k.accountKeeper.GetAccount(ctx, addr).GetAccountNumber()
		k.userNewTotalStake[accountNumber] = uint64(stake)
	}
	k.accountToUpdate = make([]sdk.AccAddress, 0)
}