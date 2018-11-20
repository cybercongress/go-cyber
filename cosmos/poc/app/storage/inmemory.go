package storage

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	. "github.com/cybercongress/cyberd/cosmos/poc/app/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type InMemoryStorage struct {
	inLinks  map[CidNumber]CidLinks
	outLinks map[CidNumber]CidLinks

	cidsCount          uint64
	cidsNumbersIndexes map[Cid]CidNumber
	cidsByNumberIndex  map[CidNumber]Cid
	userStake          map[AccountNumber]uint64
}

func (s *InMemoryStorage) Empty() {
	s.inLinks = make(map[CidNumber]CidLinks)
	s.outLinks = make(map[CidNumber]CidLinks)
	s.cidsNumbersIndexes = make(map[Cid]CidNumber)
	s.cidsByNumberIndex = make(map[CidNumber]Cid)
	s.userStake = make(map[AccountNumber]uint64)
}

// Load from underlying persistent storage
// Heavy operation
func (s *InMemoryStorage) Load(ctx sdk.Context, ps CyberdPersistentStorages, am auth.AccountKeeper) {

	inLinks, outLinks, err := ps.Links.GetAllLinks(ctx)
	if err != nil {
		cmn.Exit(err.Error())
	}

	cidsIndexes := ps.CidIndex.GetFullCidsIndex(ctx)
	cidsByNumber := make(map[CidNumber]Cid)
	for cid, cidNumber := range cidsIndexes {
		cidsByNumber[cidNumber] = cid
	}

	s.inLinks = inLinks
	s.outLinks = outLinks
	s.cidsNumbersIndexes = cidsIndexes
	s.cidsByNumberIndex = cidsByNumber
	s.cidsCount = uint64(len(cidsIndexes))
	//s.userStake = GetAllAccountsStakes(ctx, am)
}

func (s *InMemoryStorage) UpdateStakeByNumber(acc AccountNumber, stake int64) {
	s.userStake[acc] += uint64(stake)
}

func (s *InMemoryStorage) AddLink(link Link) {
	Links(s.outLinks).Put(link.From(), link.To(), link.Acc())
	Links(s.inLinks).Put(link.To(), link.From(), link.Acc())
}

func (s *InMemoryStorage) AddCid(cid Cid, number CidNumber) {
	s.cidsNumbersIndexes[cid] = number
	s.cidsByNumberIndex[number] = cid
}

//
func (s *InMemoryStorage) GetOverallLinkStake(from CidNumber, to CidNumber) uint64 {

	stake := uint64(0)
	users := s.outLinks[from][to]
	for user := range users {
		stake += s.userStake[user]
	}
	return stake
}

func (s *InMemoryStorage) GetOverallOutLinksStake(from CidNumber) uint64 {

	stake := uint64(0)
	for to := range s.outLinks[from] {
		stake += s.GetOverallLinkStake(from, to)
	}
	return stake
}

//
// GETTERS
func (s *InMemoryStorage) GetInLinks() map[CidNumber]CidLinks {
	return s.inLinks
}

func (s *InMemoryStorage) GetOutLinks() map[CidNumber]CidLinks {
	return s.outLinks
}

func (s *InMemoryStorage) GetStakes() map[AccountNumber]uint64 {
	return s.userStake
}

func (s *InMemoryStorage) GetCidsCount() uint64 {
	return uint64(len(s.cidsNumbersIndexes))
}
