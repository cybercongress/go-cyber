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
	cidRank     []float64 // array index is cid number

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

func (s *InMemoryStorage) AddCid(cid Cid, number CidNumber) {
	s.cidsIndexes[cid] = number
}

//
func (s *InMemoryStorage) GetOverallLinkStake(from CidNumber, to CidNumber) int64 {

	stake := int64(0)
	users := s.outLinks[from][to]
	for user := range users {
		stake += s.userStake[user]
	}
	return stake
}

func (s *InMemoryStorage) GetOverallOutLinksStake(from CidNumber) int64 {

	stake := int64(0)
	for to := range s.outLinks[from] {
		stake += s.GetOverallLinkStake(from, to)
	}
	return stake
}

func (s *InMemoryStorage) UpdateRank(newCidRank []float64) {
	s.cidRank = newCidRank
}

//
// GETTERS
func (s *InMemoryStorage) GetRank() []float64 {
	return s.cidRank
}

func (s *InMemoryStorage) GetInLinks() map[CidNumber]CidLinks {
	return s.inLinks
}

func (s *InMemoryStorage) GetOutLinks() map[CidNumber]CidLinks {
	return s.outLinks
}

func (s *InMemoryStorage) GetStakes() map[AccountNumber]int64 {
	return s.userStake
}

func (s *InMemoryStorage) GetCidsCount() int {
	return len(s.cidsIndexes)
}
