package rank

import (
	"errors"
	"github.com/cybercongress/cyberd/x/link/keeper"
	. "github.com/cybercongress/cyberd/x/link/types"
	"sort"
)

type RankedCidNumber struct {
	number CidNumber
	rank   float64
}

func (c RankedCidNumber) GetNumber() CidNumber { return c.number }
func (c RankedCidNumber) GetRank() float64     { return c.rank }

type RankState struct {
	linkIndex *keeper.LinkIndexedKeeper

	cidRankedLinksIndex []cidRankedLinks
	cidRank             []float64 // array index is cid number
}

func NewRankState(linkIndex *keeper.LinkIndexedKeeper) *RankState {

	return &RankState{
		linkIndex: linkIndex,
	}
}

func (s *RankState) GetCidRankedLinks(cidNumber CidNumber, page, perPage int) ([]RankedCidNumber, int, error) {

	cidRankedLinks := s.cidRankedLinksIndex[cidNumber]
	if len(cidRankedLinks) == 0 {
		return nil, 0, errors.New("no links for this cid")
	}

	totalSize := len(cidRankedLinks)
	startIndex := page * perPage
	if startIndex >= totalSize {
		return nil, totalSize, errors.New("page not found")
	}

	endIndex := startIndex + perPage
	if endIndex > totalSize {
		endIndex = startIndex + (totalSize % perPage)
	}

	resultSet := cidRankedLinks[startIndex:endIndex]

	return resultSet, totalSize, nil
}

func (s *RankState) UpdateRank(newCidRank []float64, cidsCount int64) {
	s.cidRank = newCidRank
	s.buildCidRankedLinksIndex(cidsCount)
}

func (s *RankState) buildCidRankedLinksIndex(cidsCount int64) {

	newIndex := make([]cidRankedLinks, cidsCount)

	for cidNumber := CidNumber(0); cidNumber < CidNumber(cidsCount); cidNumber++ {
		cidOutLinks := s.linkIndex.GetOutLinks()[cidNumber]
		cidSortedByRankLinkedCids := s.getLinksSortedByRank(cidOutLinks)
		newIndex[cidNumber] = cidSortedByRankLinkedCids
	}

	s.cidRankedLinksIndex = newIndex
}

func (s *RankState) getLinksSortedByRank(cidOutLinks CidLinks) cidRankedLinks {
	cidRankedLinks := make(cidRankedLinks, 0, len(cidOutLinks))
	for linkedCidNumber := range cidOutLinks {
		rankedCid := RankedCidNumber{number: linkedCidNumber, rank: s.cidRank[linkedCidNumber]}
		cidRankedLinks = append(cidRankedLinks, rankedCid)
	}
	sort.Stable(sort.Reverse(cidRankedLinks))
	return cidRankedLinks
}

//
// GETTERS

//rank index
func (s *RankState) GetRank() []float64 {
	return s.cidRank
}

//
// Local type for sorting
type cidRankedLinks []RankedCidNumber

// Sort Interface functions
func (links cidRankedLinks) Len() int { return len(links) }

func (links cidRankedLinks) Less(i, j int) bool { return links[i].rank < links[j].rank }

func (links cidRankedLinks) Swap(i, j int) { links[i], links[j] = links[j], links[i] }
