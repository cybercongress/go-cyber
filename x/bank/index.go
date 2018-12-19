package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cbd "github.com/cybercongress/cyberd/types"
)

type IndexedKeeper struct {
	Keeper

	accKeeper auth.AccountKeeper

	userStake map[cbd.AccNumber]uint64
}

func NewIndexedKeeper(keeper Keeper, accKeeper auth.AccountKeeper) IndexedKeeper {
	return IndexedKeeper{Keeper: keeper, accKeeper: accKeeper}
}

func (s *IndexedKeeper) Empty() {
	s.userStake = make(map[cbd.AccNumber]uint64)
}

func (s *IndexedKeeper) Load(ctx sdk.Context) {

	s.userStake = make(map[cbd.AccNumber]uint64)

	collect := func(acc auth.Account) bool {
		balance := s.Keeper.GetAccountTotalStake(ctx, acc.GetAddress())
		s.userStake[cbd.AccNumber(acc.GetAccountNumber())] = uint64(balance)
		return false
	}

	s.accKeeper.IterateAccounts(ctx, collect)
}

func (s *IndexedKeeper) UpdateStake(acc cbd.AccNumber, stake int64) {
	s.userStake[acc] += uint64(stake)
}

func (s *IndexedKeeper) GetStakes() map[cbd.AccNumber]uint64 {
	return s.userStake
}
