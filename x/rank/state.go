package rank

import (
	"bytes"
	"crypto/sha256"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/merkle"
	"github.com/cybercongress/cyberd/store"
	. "github.com/cybercongress/cyberd/x/link/types"
	"github.com/tendermint/tendermint/libs/log"
	"sort"
)

type RankedCidNumber struct {
	number CidNumber
	rank   float64
}

func (c RankedCidNumber) GetNumber() CidNumber { return c.number }
func (c RankedCidNumber) GetRank() float64     { return c.rank }

type RankState struct {
	cidRankedLinksIndex []cidRankedLinks
	networkCidRank      Rank // array index is cid number
	nextCidRank         Rank // array index is cid number

	lastCalculatedRankHash []byte

	allowSearch  bool
	indexErrChan chan error
}

func NewRankState(allowSearch bool) *RankState {
	return &RankState{
		allowSearch: allowSearch, indexErrChan: make(chan error),
	}
}

func (s *RankState) Load(ctx sdk.Context, mainKeeper store.MainKeeper) {
	s.lastCalculatedRankHash = mainKeeper.GetLastCalculatedRankHash(ctx)
	s.networkCidRank = Rank{Values: nil, MerkleTree: merkle.NewTree(sha256.New(), false)}
	s.nextCidRank = Rank{Values: nil, MerkleTree: merkle.NewTree(sha256.New(), false)}
	s.networkCidRank.MerkleTree.ImportSubtreesRoots(mainKeeper.GetLatestMerkleTree(ctx))
	s.nextCidRank.MerkleTree.ImportSubtreesRoots(mainKeeper.GetNextMerkleTree(ctx))
}

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

func (s *RankState) SetNextRank(newRank Rank) {
	s.nextCidRank = newRank
}

func (s *RankState) ApplyNextRank() {
	s.networkCidRank = s.nextCidRank
	s.lastCalculatedRankHash = s.nextCidRank.MerkleTree.RootHash()
	s.nextCidRank.Reset()
}

// add new cids to rank with 0 value
func (s *RankState) AddNewCids(cidCount uint64) {
	s.networkCidRank.AddNewCids(cidCount)
}

func (s *RankState) BuildCidRankedLinksIndex(cidsCount int64, outLinks Links) {
	// If search on this node is not allowed then we don't need to build index
	if !s.allowSearch || s.networkCidRank.Values == nil {
		return
	}

	newIndex := make([]cidRankedLinks, cidsCount)

	for cidNumber := CidNumber(0); cidNumber < CidNumber(cidsCount); cidNumber++ {
		cidOutLinks := outLinks[cidNumber]
		cidSortedByRankLinkedCids := s.getLinksSortedByRank(cidOutLinks)
		newIndex[cidNumber] = cidSortedByRankLinkedCids
	}

	s.cidRankedLinksIndex = newIndex
}

// Used for building index in parallel
func (s *RankState) BuildCidRankedLinksIndexInParallel(cidsCount int64, outLinks Links) {
	defer func() {
		if r := recover(); r != nil {
			s.indexErrChan <- r.(error)
		}
	}()

	s.BuildCidRankedLinksIndex(cidsCount, outLinks)
}

func (s *RankState) CheckBuildIndexError(logger log.Logger) {
	if s.allowSearch {
		select {
		case err := <-s.indexErrChan:
			// DUMB ERROR HANDLING
			logger.Error("Error during building rank index " + err.Error())
			panic(err)
		default:
		}
	}
}

func (s *RankState) getLinksSortedByRank(cidOutLinks CidLinks) cidRankedLinks {
	cidRankedLinks := make(cidRankedLinks, 0, len(cidOutLinks))
	for linkedCidNumber := range cidOutLinks {
		rankedCid := RankedCidNumber{number: linkedCidNumber, rank: s.networkCidRank.Values[linkedCidNumber]}
		cidRankedLinks = append(cidRankedLinks, rankedCid)
	}
	sort.Stable(sort.Reverse(cidRankedLinks))
	return cidRankedLinks
}

//
// GETTERS

//rank index
func (s *RankState) SearchAllowed() bool {
	return s.allowSearch
}

func (s *RankState) GetNetworkRankHash() []byte {
	return s.networkCidRank.MerkleTree.RootHash()
}

func (s *RankState) GetNetworkMerkleTreeAsBytes() []byte {
	return s.networkCidRank.MerkleTree.ExportSubtreesRoots()
}

func (s *RankState) GetNextMerkleTreeAsBytes() []byte {
	return s.nextCidRank.MerkleTree.ExportSubtreesRoots()
}

func (s *RankState) NextRankReady() bool {
	return s.lastCalculatedRankHash != nil && !bytes.Equal(s.nextCidRank.MerkleTree.RootHash(), s.lastCalculatedRankHash)
}

func (s *RankState) GetLastCidNum() CidNumber {
	return CidNumber(len(s.networkCidRank.Values) - 1)
}

//
// Local type for sorting
type cidRankedLinks []RankedCidNumber

// Sort Interface functions
func (links cidRankedLinks) Len() int { return len(links) }

func (links cidRankedLinks) Less(i, j int) bool { return links[i].rank < links[j].rank }

func (links cidRankedLinks) Swap(i, j int) { links[i], links[j] = links[j], links[i] }
