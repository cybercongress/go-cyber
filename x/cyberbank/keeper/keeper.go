package keeper

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cybercongress/go-cyber/x/cyberbank/types"
)


type IndexedKeeper struct {
	*Proxy
	accountKeeper     		types.AccountKeeper

	userTotalStakeAmpere    map[uint64]uint64
	userNewTotalStakeAmpere map[uint64]uint64
	accountToUpdate         []sdk.AccAddress
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

func (k *IndexedKeeper) LoadState(rankCtx sdk.Context, freshCtx sdk.Context) {
	k.userTotalStakeAmpere = make(map[uint64]uint64)
	k.accountKeeper.IterateAccounts(rankCtx, k.getCollectFunc(rankCtx, k.userTotalStakeAmpere))

	k.userNewTotalStakeAmpere = make(map[uint64]uint64)
	k.accountKeeper.IterateAccounts(freshCtx, k.getCollectFunc(freshCtx, k.userNewTotalStakeAmpere))
}

func (k *IndexedKeeper) getCollectFunc(ctx sdk.Context, userStake map[uint64]uint64) func(account authtypes.AccountI) bool {
	return func(account authtypes.AccountI) bool {
		balance := k.Proxy.GetAccountTotalStakeAmper(ctx, account.GetAddress())
		userStake[account.GetAccountNumber()] = uint64(balance)
		return false
	}
}

func  (k *IndexedKeeper) InitializeStakeAmpere(account uint64, stake uint64) {
	k.userTotalStakeAmpere[account] = stake
}

func (k *IndexedKeeper) GetTotalStakesAmpere() map[uint64]uint64 {
	return k.userTotalStakeAmpere
}

func (k *IndexedKeeper) DetectUsersStakeAmpereChange(ctx sdk.Context) bool {
	stakeChanged := false
	for o, n := range k.userNewTotalStakeAmpere {
		if k.userTotalStakeAmpere[o] != n {
			stakeChanged = true
			k.userTotalStakeAmpere[o] = n
		}
	}
	return stakeChanged
}

func (k *IndexedKeeper) UpdateAccountsStakeAmpere(ctx sdk.Context) {
	for _, addr := range k.accountToUpdate {
		stake := k.GetAccountTotalStakeAmper(ctx, addr)
		fmt.Printf("[%s] %s \n", addr.String(), strconv.FormatUint(uint64(stake), 10))
		if k.accountKeeper.GetAccount(ctx, addr) == nil { continue }
		accountNumber := k.accountKeeper.GetAccount(ctx, addr).GetAccountNumber()
		k.userNewTotalStakeAmpere[accountNumber] = uint64(stake)
	}

	k.accountToUpdate = make([]sdk.AccAddress, 0)
}