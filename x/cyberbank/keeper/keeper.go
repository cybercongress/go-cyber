package keeper

import (
	"fmt"
	"math"
	"time"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	ctypes "github.com/cybercongress/go-cyber/types"

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
	energyKeeper  types.EnergyKeeper
	authKey       sdk.StoreKey
	cdc           codec.BinaryCodec

	userTotalStakeAmpere    map[uint64]uint64
	userNewTotalStakeAmpere map[uint64]uint64
	accountToUpdate         []sdk.AccAddress
}

func NewIndexedKeeper(
	cdc codec.BinaryCodec,
	pbk *BankProxyKeeper,
	ak types.AccountKeeper,
	authKey sdk.StoreKey,
) *IndexedKeeper {
	indexedKeeper := &IndexedKeeper{
		BankProxyKeeper: pbk,
		accountKeeper:   ak,
		authKey:         authKey,
		cdc:             cdc,
		accountToUpdate: make([]sdk.AccAddress, 0),
	}
	hook := func(ctx sdk.Context, accounts []sdk.AccAddress) {
		indexedKeeper.accountToUpdate = append(indexedKeeper.accountToUpdate, accounts...)
	}
	pbk.AddBalanceListener(hook)

	return indexedKeeper
}

func (p *IndexedKeeper) SetGridKeeper(ek types.EnergyKeeper) {
	p.energyKeeper = ek
}

func (p *IndexedKeeper) SetAccountKeeper(ak authkeeper.AccountKeeper) {
	p.accountKeeper = ak
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

func (p IndexedKeeper) GetTotalSupplyVolt(ctx sdk.Context) int64 {
	return p.BankProxyKeeper.GetSupply(ctx, ctypes.VOLT).Amount.Int64()
}

func (p IndexedKeeper) GetTotalSupplyAmper(ctx sdk.Context) int64 {
	return p.BankProxyKeeper.GetSupply(ctx, ctypes.AMPERE).Amount.Int64()
}

func (p IndexedKeeper) GetAccountStakePercentageVolt(ctx sdk.Context, addr sdk.AccAddress) float64 {
	a := p.GetAccountTotalStakeVolt(ctx, addr)
	aFloat := float64(a)

	b := p.GetTotalSupplyVolt(ctx)
	bFloat := float64(b)

	c := aFloat / bFloat

	if math.IsNaN(c) {
		return 0
	}
	return c
}

func (p IndexedKeeper) GetAccountTotalStakeVolt(ctx sdk.Context, addr sdk.AccAddress) int64 {
	return p.BankProxyKeeper.GetBalance(ctx, addr, ctypes.VOLT).Amount.Int64() + p.GetRoutedTo(ctx, addr).AmountOf(ctypes.VOLT).Int64()
}

func (p IndexedKeeper) GetAccountTotalStakeAmper(ctx sdk.Context, addr sdk.AccAddress) int64 {
	return p.BankProxyKeeper.GetBalance(ctx, addr, ctypes.AMPERE).Amount.Int64() + p.GetRoutedTo(ctx, addr).AmountOf(ctypes.AMPERE).Int64()
}

func (p IndexedKeeper) GetRoutedTo(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return p.energyKeeper.GetRoutedToEnergy(ctx, addr)
}

func (k *IndexedKeeper) getCollectFunc(ctx sdk.Context, userStake map[uint64]uint64) func(account authtypes.AccountI) bool {
	return func(account authtypes.AccountI) bool {
		balance := k.GetAccountTotalStakeAmper(ctx, account.GetAddress())
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

func (k *IndexedKeeper) DetectUsersStakeAmpereChange(ctx sdk.Context) bool {
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
