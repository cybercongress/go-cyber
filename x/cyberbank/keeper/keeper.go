package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cybercongress/go-cyber/x/cyberbank/types"
)

type IndexedKeeper struct {
	*BankProxyKeeper
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
	pbk *BankProxyKeeper,
	ak types.AccountKeeper,
) *IndexedKeeper {
	indexedKeeper := &IndexedKeeper{
		BankProxyKeeper: pbk,
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
	pbk.AddBalanceListener(hook)

	return indexedKeeper
}

func (ik IndexedKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (ik *IndexedKeeper) LoadState(rankCtx sdk.Context, freshCtx sdk.Context) {
	ik.userTotalStakeAmpere = make(map[uint64]uint64)
	ik.accountKeeper.IterateAccounts(rankCtx, ik.getCollectFunc(rankCtx, ik.userTotalStakeAmpere))

	ik.userNewTotalStakeAmpere = make(map[uint64]uint64)
	ik.accountKeeper.IterateAccounts(freshCtx, ik.getCollectFunc(freshCtx, ik.userNewTotalStakeAmpere))
}

func (ik *IndexedKeeper) getCollectFunc(ctx sdk.Context, userStake map[uint64]uint64) func(account authtypes.AccountI) bool {
	return func(account authtypes.AccountI) bool {
		balance := ik.BankProxyKeeper.GetAccountTotalStakeAmper(ctx, account.GetAddress())
		userStake[account.GetAccountNumber()] = uint64(balance)
		return false
	}
}

func (ik *IndexedKeeper) InitializeStakeAmpere(account uint64, stake uint64) {
	ik.userTotalStakeAmpere[account] = stake
	ik.userNewTotalStakeAmpere[account] = stake
}

func (ik *IndexedKeeper) GetTotalStakesAmpere() map[uint64]uint64 {
	return ik.userTotalStakeAmpere
}

func (ik *IndexedKeeper) DetectUsersStakeAmpereChange(ctx sdk.Context) bool {
	stakeChanged := false
	for o, n := range ik.userNewTotalStakeAmpere {
		if _, ok := ik.userTotalStakeAmpere[o]; ok {
			if ik.userTotalStakeAmpere[o] != n {
				stakeChanged = true
				ik.userTotalStakeAmpere[o] = n
			}
		} else {
			ik.userTotalStakeAmpere[o] = n
		}
	}

	return stakeChanged
}

func (ik *IndexedKeeper) UpdateAccountsStakeAmpere(ctx sdk.Context) {
	for _, addr := range ik.accountToUpdate {
		ik.Logger(ctx).Debug("account to update:", "address", addr.String())
		stake := ik.GetAccountTotalStakeAmper(ctx, addr)
		if ik.accountKeeper.GetAccount(ctx, addr) == nil {
			ik.Logger(ctx).Info("skipped account:", "address", addr.String())
			continue
		}
		accountNumber := ik.accountKeeper.GetAccount(ctx, addr).GetAccountNumber()
		ik.userNewTotalStakeAmpere[accountNumber] = uint64(stake)
	}

	// trigger full account map rebuild in case of account' missing (and if new contract deployed)
	// TODO migrate logic to storage listener in sdk 46?
	// NOTE returns last not applied yet next! account number
	// equal to current length of accounts ids array, but last id is equal to next-1
	nextAccountNumber := ik.GetNextAccountNumber(ctx)
	if uint64(len(ik.userNewTotalStakeAmpere)) != nextAccountNumber {
		startTime := time.Now()
		for i := nextAccountNumber - 1; i > 0; i-- {
			if _, ok := ik.userNewTotalStakeAmpere[i]; !ok {
				ik.Logger(ctx).Info("added to stake index:", "account", i)
				// TODO update in next release
				// stake := ik.GetAccountTotalStakeAmper(ctx, addr)
				ik.userNewTotalStakeAmpere[i] = 0
			}
		}
		ik.Logger(ctx).Info("rebuild stake index:", "duration", time.Since(startTime).String())
	}

	ik.accountToUpdate = make([]sdk.AccAddress, 0)
}

func (ik IndexedKeeper) GetNextAccountNumber(ctx sdk.Context) uint64 {
	var accNumber uint64
	store := ctx.KVStore(ik.authKey)
	bz := store.Get([]byte("globalAccountNumber"))

	if bz == nil {
		accNumber = 0
	} else {
		val := gogotypes.UInt64Value{}

		err := ik.cdc.Unmarshal(bz, &val)
		if err != nil {
			panic(err)
		}

		accNumber = val.GetValue()
	}

	return accNumber
}
