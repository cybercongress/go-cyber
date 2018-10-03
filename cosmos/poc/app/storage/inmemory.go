package storage

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type InMemoryStorage struct {
	inLinks  map[CidNumber]CidLinks
	outLinks map[CidNumber]CidLinks

	cidsCount   uint64
	cidsIndexes map[Cid]CidNumber
	cidRank     []float64

	userStake map[AccountNumber]int64
}

// Load from underlying persistent storage
// Heavy operation
func (s *InMemoryStorage) Load(ctx sdk.Context, ps CyberdPersistentStorages, am auth.AccountMapper) {

	inLinks, outLinks, err := ps.Links.GetAllLinks(ctx)
	if err != nil {
		cmn.Exit(err.Error())
	}

	cidsIndexes := ps.CidIndex.GetFullCidsIndex(ctx)

	s.inLinks = inLinks
	s.outLinks = outLinks
	s.cidsIndexes = cidsIndexes
	s.cidsCount = uint64(len(cidsIndexes))
	s.userStake = GetAllAccountsStakes(ctx, am)
	s.cidRank = ps.Rank.GetFullRank(ctx)
}

// Also returns bool flag, whether index exists
func (s *InMemoryStorage) GetCidIndex(cid Cid) (CidNumber, bool) {

	if index, ok := s.cidsIndexes[cid]; ok {
		return index, true
	}
	return 0, false
}

func (s *InMemoryStorage) UpdateStake(acc sdk.AccAddress, stake int64) {
	s.userStake[AccountNumber(acc.String())] += stake
}

func (s *InMemoryStorage) AddLink(link LinkedCids) {

	CidsLinks(s.outLinks).Put(link.FromCid, link.ToCid, link.Creator)
	CidsLinks(s.inLinks).Put(link.ToCid, link.FromCid, link.Creator)
}
