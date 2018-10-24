package storage

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/pkg/errors"
	cmn "github.com/tendermint/tendermint/libs/common"
	"sort"
)

type RankedCidNumber struct {
	cidNumber CidNumber
	rank      float64
}

type RankedCid struct {
	Cid  Cid
	Rank float64 `amino:"unsafe"`
}

type CidRankedLinks []RankedCidNumber

// Sort Interface functions
func (links CidRankedLinks) Len() int { return len(links) }

func (links CidRankedLinks) Less(i, j int) bool { return links[i].rank < links[j].rank }

func (links CidRankedLinks) Swap(i, j int) { links[i], links[j] = links[j], links[i] }

type InMemoryStorage struct {
	inLinks  map[CidNumber]CidLinks
	outLinks map[CidNumber]CidLinks

	cidsCount          uint64
	cidsNumbersIndexes map[Cid]CidNumber
	cidsByNumberIndex  map[CidNumber]Cid
	cidRank            []float64 // array index is cid number

	cidRankedLinksIndex []CidRankedLinks

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
	cidsByNumber := make(map[CidNumber]Cid)
	for cid, cidNumber := range cidsIndexes {
		cidsByNumber[cidNumber] = cid
	}

	s.inLinks = inLinks
	s.outLinks = outLinks
	s.cidsNumbersIndexes = cidsIndexes
	s.cidsByNumberIndex = cidsByNumber
	s.cidsCount = uint64(len(cidsIndexes))
	s.userStake = GetAllAccountsStakes(ctx, am)
}

// Also returns bool flag, whether index exists
func (s *InMemoryStorage) GetCidIndex(cid Cid) (CidNumber, bool) {

	if index, ok := s.cidsNumbersIndexes[cid]; ok {
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
	s.cidsNumbersIndexes[cid] = number
	s.cidsByNumberIndex[number] = cid
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

func (s *InMemoryStorage) GetCidRankedLinks(cid Cid, page, perPage int) ([]RankedCid, int, error) {

	cidNumber, exists := s.cidsNumbersIndexes[cid]
	if !exists {
		return nil, 0, errors.New("Searched cid not found")
	}

	cidRankedLinks := s.cidRankedLinksIndex[cidNumber]
	if len(cidRankedLinks) == 0 {
		return nil, 0, errors.New("No links for this cid")
	}

	totalSize := len(cidRankedLinks)
	startIndex := page * perPage
	if startIndex >= totalSize {
		return nil, totalSize, errors.New("Page not found")
	}

	endIndex := startIndex + perPage
	if endIndex > totalSize {
		endIndex = startIndex + (totalSize % perPage)
	}

	resultSet := cidRankedLinks[startIndex:endIndex]
	response := make([]RankedCid, 0, len(resultSet))
	for _, result := range resultSet {
		response = append(response, RankedCid{s.cidsByNumberIndex[result.cidNumber], result.rank})
	}

	return response, totalSize, nil
}

func (s *InMemoryStorage) UpdateRank(newCidRank []float64) {
	s.cidRank = newCidRank
	s.buildCidRankedLinksIndex()
}

func (s *InMemoryStorage) buildCidRankedLinksIndex() {

	newIndex := make([]CidRankedLinks, len(s.cidsNumbersIndexes))

	for _, cidNumber := range s.cidsNumbersIndexes {
		cidOutLinks := s.outLinks[cidNumber]
		cidRankedLinks := getLinksSortedByRank(cidOutLinks, s.cidRank)
		newIndex[cidNumber] = cidRankedLinks
	}

	s.cidRankedLinksIndex = newIndex
}

func getLinksSortedByRank(cidOutLinks CidLinks, cidRank []float64) CidRankedLinks {
	cidRankedLinks := make(CidRankedLinks, 0, len(cidOutLinks))
	for linkedCidNumber := range cidOutLinks {
		rankedCid := RankedCidNumber{cidNumber: linkedCidNumber, rank: cidRank[linkedCidNumber]}
		cidRankedLinks = append(cidRankedLinks, rankedCid)
	}
	sort.Stable(sort.Reverse(cidRankedLinks))
	return cidRankedLinks
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
	return len(s.cidsNumbersIndexes)
}

func (s *InMemoryStorage) PrintCidsByNumberIndex() {
	var keys []int
	for k := range s.cidsByNumberIndex {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Println("Key:", k, "Value:", s.cidsByNumberIndex[CidNumber(k)])
	}
}
