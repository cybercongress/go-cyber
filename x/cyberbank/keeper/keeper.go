package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cybercongress/go-cyber/x/cyberbank/types"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/tendermint/tendermint/libs/log"
)

type IndexedKeeper struct {
	*Proxy
	accountKeeper types.AccountKeeper
	cdc           codec.BinaryCodec
	authKey       sdk.StoreKey

	userTotalStakeAmpere    map[uint64]uint64
	userNewTotalStakeAmpere map[uint64]uint64
	accountToUpdate         []sdk.AccAddress
}

func NewIndexedKeeper(
	cdc codec.BinaryCodec,
	authKey sdk.StoreKey,
	pbk *Proxy,
	ak types.AccountKeeper,
) *IndexedKeeper {
	indexedKeeper := &IndexedKeeper{
		Proxy:           pbk,
		cdc:             cdc,
		authKey:         authKey,
		accountKeeper:   ak,
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

func (k *IndexedKeeper) InitializeStakeAmpere(account uint64, stake uint64) {
	k.userTotalStakeAmpere[account] = stake
	k.userNewTotalStakeAmpere[account] = stake
}

func (k *IndexedKeeper) GetTotalStakesAmpere() map[uint64]uint64 {
	return k.userTotalStakeAmpere
}

func (k *IndexedKeeper) DetectUsersStakeAmpereChange() bool {
	stakeChanged := false
	for o, n := range k.userNewTotalStakeAmpere {
		if _, ok := k.userTotalStakeAmpere[o]; ok {
			if k.userTotalStakeAmpere[o] != n {
				stakeChanged = true
				k.userTotalStakeAmpere[o] = n
			}
		} else {
			k.userTotalStakeAmpere[o] = n
		}
	}

	return stakeChanged
}

func (k *IndexedKeeper) UpdateAccountsStakeAmpere(ctx sdk.Context) {
	for _, addr := range k.accountToUpdate {
		k.Logger(ctx).Debug("account to update:", "address", addr.String())
		stake := k.GetAccountTotalStakeAmper(ctx, addr)
		if k.accountKeeper.GetAccount(ctx, addr) == nil {
			k.Logger(ctx).Info("skipped account:", "address", addr.String())
			continue
		}
		accountNumber := k.accountKeeper.GetAccount(ctx, addr).GetAccountNumber()
		k.userNewTotalStakeAmpere[accountNumber] = uint64(stake)
	}

	// trigger full account map rebuild in case of account' missing (and if new contract deployed)
	// TODO migrate logic to storage listener in sdk 46?
	// NOTE returns last not applied yet next! account number
	// equal to current length of accounts ids array, but last id is equal to next-1
	nextAccountNumber := k.GetNextAccountNumber(ctx)
	if uint64(len(k.userNewTotalStakeAmpere)) != nextAccountNumber {
		startTime := time.Now()
		for i := nextAccountNumber - 1; i > 0; i-- {
			if _, ok := k.userNewTotalStakeAmpere[i]; !ok {
				k.Logger(ctx).Info("added to stake index:", "account", i)
				// TODO update in next release
				// stake := k.GetAccountTotalStakeAmper(ctx, addr)
				k.userNewTotalStakeAmpere[i] = 0
			}
		}
		k.Logger(ctx).Info("rebuild stake index:", "duration", time.Since(startTime).String())
	}

	k.accountToUpdate = make([]sdk.AccAddress, 0)
}

func (k IndexedKeeper) GetNextAccountNumber(ctx sdk.Context) uint64 {
	var accNumber uint64
	store := ctx.KVStore(k.authKey)
	bz := store.Get([]byte("globalAccountNumber"))

	if bz == nil {
		accNumber = 0
	} else {
		val := gogotypes.UInt64Value{}

		err := k.cdc.Unmarshal(bz, &val)
		if err != nil {
			panic(err)
		}

		accNumber = val.GetValue()
	}

	return accNumber
}
