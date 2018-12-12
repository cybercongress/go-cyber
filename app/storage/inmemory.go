package storage

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	. "github.com/cybercongress/cyberd/app/types"
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

	userStake map[AccountNumber]uint64
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
	s.userStake = GetAllAccountsStakes(ctx, am)
}

// Also returns bool flag, whether index exists
func (s *InMemoryStorage) GetCidIndex(cid Cid) (CidNumber, bool) {

	if index, ok := s.cidsNumbersIndexes[cid]; ok {
		return index, true
	}
	return 0, false
}

func (s *InMemoryStorage) UpdateStake(acc AccountNumber, stake int64) {
	s.userStake[acc] += uint64(stake)
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

// sorted by cid
func (s *InMemoryStorage) GetSortedInLinks(cid CidNumber) (CidLinks, []CidNumber, bool) {
	links := s.inLinks[cid]

	if len(links) == 0 {
		return nil, nil, false
	}

	numbers := make([]CidNumber, 0, len(links))
	for num := range links {
		numbers = append(numbers, num)
	}

	sort.Slice(numbers, func(i, j int) bool { return numbers[i] < numbers[j] })

	return links, numbers, true
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
