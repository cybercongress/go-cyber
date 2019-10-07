package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	acc "github.com/cybercongress/cyberd/x/acc/types"
)

// Used to hodl total user stake in memory for further rank calculation.
// Updated once at block and the beginning of end block.
type IndexedKeeper struct {
	Keeper

	accKeeper acc.AccountIndexKeeper

	userTotalStake    map[acc.AccNumber]uint64
	userNewTotalStake map[acc.AccNumber]uint64

	// used to track accs with changed stake
	accsToUpdate []sdk.AccAddress
}

func NewIndexedKeeper(keeper *Keeper, accKeeper acc.AccountIndexKeeper) *IndexedKeeper {
	index := IndexedKeeper{Keeper: *keeper, accKeeper: accKeeper, accsToUpdate: make([]sdk.AccAddress, 0)}
	hook := func(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress) {
		if from != nil {
			index.accsToUpdate = append(index.accsToUpdate, from)
		}
		if to != nil {
			index.accsToUpdate = append(index.accsToUpdate, to)
		}
	}
	keeper.AddHook(hook)
	return &index
}

// todo: how to load only new stakes from last n blocks? We could iterate over whole db and compare stakes by address and amount.
func (s *IndexedKeeper) Load(rankCtx sdk.Context, freshCtx sdk.Context) {

	s.userTotalStake = make(map[acc.AccNumber]uint64)
	s.accKeeper.IterateAccounts(rankCtx, s.getCollectFunc(rankCtx, s.userTotalStake))

	s.userNewTotalStake = make(map[acc.AccNumber]uint64)
	s.accKeeper.IterateAccounts(freshCtx, s.getCollectFunc(freshCtx, s.userNewTotalStake))
}

func (s *IndexedKeeper) getCollectFunc(ctx sdk.Context, userStake map[acc.AccNumber]uint64) func(auth.Account) bool {
	return func(account auth.Account) bool {
		balance := s.Keeper.GetAccountTotalStake(ctx, account.GetAddress())
		userStake[acc.AccNumber(account.GetAccountNumber())] = uint64(balance)
		return false
	}
}

// return true if some stake changed
func (s *IndexedKeeper) FixUserStake() bool {
	stakeChanged := false
	for k, v := range s.userNewTotalStake {
		if s.userTotalStake[k] != v {
			stakeChanged = true
			s.userTotalStake[k] = v
		}
	}
	return stakeChanged
}

func (s *IndexedKeeper) UpdateStake(accNum acc.AccNumber, stake int64) {
	s.userNewTotalStake[accNum] += uint64(stake)
}

func (s *IndexedKeeper) GetTotalStakes() map[acc.AccNumber]uint64 {
	return s.userTotalStake
}

// Performs stakes updates for accountKeeper touched in current block
func (s *IndexedKeeper) EndBlocker(ctx sdk.Context) {
	for _, addr := range s.accsToUpdate {
		stake := s.Keeper.GetAccountTotalStake(ctx, addr)
		accNum := acc.AccNumber(s.accKeeper.GetAccount(ctx, addr).GetAccountNumber())
		s.userNewTotalStake[accNum] = uint64(stake)
	}
	s.accsToUpdate = make([]sdk.AccAddress, 0)
}
