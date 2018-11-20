package types

import (
	"github.com/pkg/errors"
	"sort"
)

type CidWithRank struct {
	Cid  Cid
	Rank float64
}

// internal representation
type CidNumberWithRank struct {
	cidNumber CidNumber
	rank      float64
}

func (c CidNumberWithRank) CidNumber() CidNumber { return c.cidNumber }
func (c CidNumberWithRank) Rank() float64        { return c.rank }

/* Represents calculated rank with various indices */
type Rank struct {
	cidsCount uint64

	// array index - cid number
	rank []float64

	// array index - cid number
	// for all cids create index of sorted by rank connected cids.
	sortedByRankLinkedCids []cidsWithRanks
}

func NewRank(rank []float64, outLinks Links, cidsCount uint64) Rank {

	index := make([]cidsWithRanks, cidsCount)
	for cidNumber := CidNumber(0); cidNumber < CidNumber(cidsCount); cidsCount++ {
		cidOutLinks := outLinks[cidNumber]
		cidSortedByRankLinkedCids := getSortedByRankLinkedCids(cidOutLinks, rank)
		index[cidNumber] = cidSortedByRankLinkedCids
	}
	return Rank{
		cidsCount: cidsCount, rank: rank,
		sortedByRankLinkedCids: index,
	}
}

// sort given cids (represent by out links) by rank value
func getSortedByRankLinkedCids(cidOutLinks CidLinks, cidRank []float64) cidsWithRanks {
	cids := make(cidsWithRanks, 0, len(cidOutLinks))
	for linkedCidNumber := range cidOutLinks {
		rankedCid := CidNumberWithRank{cidNumber: linkedCidNumber, rank: cidRank[linkedCidNumber]}
		cids = append(cids, rankedCid)
	}
	sort.Stable(sort.Reverse(cids))
	return cids
}

func (s *Rank) GetSortedByRankLinkedCids(cid CidNumber, page, perPage int) ([]CidNumberWithRank, int, error) {

	if cid >= CidNumber(s.cidsCount) {
		return nil, 0, nil
	}

	sortedCids := s.sortedByRankLinkedCids[cid]
	if len(sortedCids) == 0 {
		return nil, 0, nil
	}

	totalSize := len(sortedCids)
	startIndex := page * perPage
	if startIndex >= totalSize {
		return nil, totalSize, errors.New("Page not found")
	}

	endIndex := startIndex + perPage
	if endIndex > totalSize {
		endIndex = startIndex + (totalSize % perPage)
	}

	resultSet := sortedCids[startIndex:endIndex]
	return resultSet, totalSize, nil
}

type cidsWithRanks []CidNumberWithRank

// Sort Interface functions
func (links cidsWithRanks) Len() int { return len(links) }

func (links cidsWithRanks) Less(i, j int) bool { return links[i].rank < links[j].rank }

func (links cidsWithRanks) Swap(i, j int) { links[i], links[j] = links[j], links[i] }
