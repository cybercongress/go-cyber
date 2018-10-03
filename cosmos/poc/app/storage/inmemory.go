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

	// persistent storages
	persistentStorage CyberdPersistentStorages
	am                auth.AccountMapper
}

func NewInMemoryStorage(persistentStorage CyberdPersistentStorages, am auth.AccountMapper) *InMemoryStorage {

	return &InMemoryStorage{
		persistentStorage: persistentStorage,
		am:                am,
	}
}

// Load from underlying persistent storage
// Heavy operation
func (s *InMemoryStorage) Load(ctx sdk.Context) {

	inLinks, outLinks, err := s.persistentStorage.InLinks.GetAllLinks(ctx)
	if err != nil {
		cmn.Exit(err.Error())
	}

	cidsIndexes := s.persistentStorage.CidIndex.GetFullCidsIndex(ctx)

	s.inLinks = inLinks
	s.outLinks = outLinks
	s.cidsIndexes = cidsIndexes
	s.cidsCount = uint64(len(cidsIndexes))
	s.userStake = GetAllAccountsStakes(ctx, s.am)
	s.cidRank = s.persistentStorage.Rank.GetFullRank(ctx)
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
	//
	// out links
	cidLinks := s.outLinks[link.FromCid]
	if cidLinks == nil {
		cidLinks = make(CidLinks)
	}
	users := cidLinks[link.ToCid]
	if users == nil {
		users = make(map[AccountNumber]struct{})
	}
	users[link.Creator] = struct{}{}
	cidLinks[link.ToCid] = users
	s.outLinks[link.FromCid] = cidLinks

	//
	// in links
	cidLinks = s.inLinks[link.ToCid]
	if cidLinks == nil {
		cidLinks = make(CidLinks)
	}
	users = cidLinks[link.FromCid]
	if users == nil {
		users = make(map[AccountNumber]struct{})
	}
	users[link.Creator] = struct{}{}
	cidLinks[link.FromCid] = users
	s.inLinks[link.ToCid] = cidLinks
}
