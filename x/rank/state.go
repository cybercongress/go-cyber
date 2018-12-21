package rank

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/store"
	"github.com/cybercongress/cyberd/x/link/keeper"
	. "github.com/cybercongress/cyberd/x/link/types"
	"github.com/tendermint/tendermint/libs/log"
	"math"
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
	networkCidRank      []float64 // array index is cid number
	nextCidRank         []float64 // array index is cid number

	networkRankHash []byte
	nextRankHash    []byte

	allowSearch  bool
	indexErrChan chan error
}

func NewRankState(linkIndex *keeper.LinkIndexedKeeper, allowSearch bool) *RankState {

	return &RankState{
		linkIndex: linkIndex, allowSearch: allowSearch, indexErrChan: make(chan error),
	}
}

func (s *RankState) Load(ctx sdk.Context, mainKeeper store.MainKeeper) {
	s.networkRankHash = mainKeeper.GetAppHash(ctx)
	s.nextRankHash = mainKeeper.GetNextAppHash(ctx)
	if s.networkRankHash == nil {
		s.networkRankHash = s.calculateHash(s.networkCidRank)
	}
	if s.nextRankHash == nil {
		s.nextRankHash = s.calculateHash(s.nextCidRank)
	}
}

// TODO: add flag to not build index at all
func (s *RankState) GetCidRankedLinks(cidNumber CidNumber, page, perPage int) ([]RankedCidNumber, int, error) {

	if !s.allowSearch {
		return nil, 0, errors.New("search on this node is not allowed")
	}

	if s.cidRankedLinksIndex == nil {
		return nil, 0, errors.New("index currently unavailable after node restart")
	}

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

func (s *RankState) SetNextRank(newRank []float64) {
	s.nextCidRank = newRank
	s.nextRankHash = s.calculateHash(newRank)
}

func (s *RankState) calculateHash(rank []float64) []byte {
	rankAsBytes := make([]byte, 8*len(rank))
	for i, f64 := range rank {
		binary.LittleEndian.PutUint64(rankAsBytes[i*8:i*8+8], math.Float64bits(f64))
	}

	hash := sha256.Sum256(rankAsBytes)
	return hash[:]
}

func (s *RankState) ApplyNextRank(cidsCount int64) {
	s.networkCidRank = s.nextCidRank
	s.networkRankHash = s.nextRankHash
	s.nextCidRank = nil
}

func (s *RankState) BuildCidRankedLinksIndex(cidsCount int64) {
	// If search on this node is not allowed then we don't need to build index
	if !s.allowSearch {
		return
	}

	newIndex := make([]cidRankedLinks, cidsCount)

	for cidNumber := CidNumber(0); cidNumber < CidNumber(cidsCount); cidNumber++ {
		cidOutLinks := s.linkIndex.GetOutLinks()[cidNumber]
		cidSortedByRankLinkedCids := s.getLinksSortedByRank(cidOutLinks)
		newIndex[cidNumber] = cidSortedByRankLinkedCids
	}

	s.cidRankedLinksIndex = newIndex
}

// Used for building index in parallel
func (s *RankState) BuildCidRankedLinksIndexInParallel(cidsCount int64) {
	defer func() {
		if r := recover(); r != nil {
			s.indexErrChan <- r.(error)
		}
	}()

	s.BuildCidRankedLinksIndex(cidsCount)
}

func (s *RankState) CheckBuildIndexError(logger log.Logger) {
	if s.allowSearch {
		select {
		case err := <- s.indexErrChan:
			// DUMB ERROR HANDLING
			logger.Error("Error during building rank index " + err.Error())
			panic(err.Error())
		default:
		}
	}
}

func (s *RankState) getLinksSortedByRank(cidOutLinks CidLinks) cidRankedLinks {
	cidRankedLinks := make(cidRankedLinks, 0, len(cidOutLinks))
	for linkedCidNumber := range cidOutLinks {
		rankedCid := RankedCidNumber{number: linkedCidNumber, rank: s.networkCidRank[linkedCidNumber]}
		cidRankedLinks = append(cidRankedLinks, rankedCid)
	}
	sort.Stable(sort.Reverse(cidRankedLinks))
	return cidRankedLinks
}

//
// GETTERS

//rank index
func (s *RankState) GetNetworkRank() []float64 {
	return s.networkCidRank
}

func (s *RankState) SearchAllowed() bool {
	return s.allowSearch
}

func (s *RankState) GetNextRankHash() []byte {
	return s.nextRankHash
}

func (s *RankState) GetCurrentRankHash() []byte {
	return s.networkRankHash
}

func (s *RankState) NextRankReady() bool {
	return !bytes.Equal(s.nextRankHash, s.networkRankHash)
}

//
// Local type for sorting
type cidRankedLinks []RankedCidNumber

// Sort Interface functions
func (links cidRankedLinks) Len() int { return len(links) }

func (links cidRankedLinks) Less(i, j int) bool { return links[i].rank < links[j].rank }

func (links cidRankedLinks) Swap(i, j int) { links[i], links[j] = links[j], links[i] }
