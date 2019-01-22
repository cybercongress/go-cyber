package rank

import (
	"bytes"
	"crypto/sha256"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/merkle"
	"github.com/cybercongress/cyberd/store"
	"github.com/cybercongress/cyberd/x/bank"
	"github.com/cybercongress/cyberd/x/link/keeper"
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

	rankCalculationFinished bool
	cidCount                int64

	rankCalcChan chan Rank
	rankErrChan  chan error
	allowSearch  bool
	indexErrChan chan error
	computeUnit  ComputeUnit

	// keepers
	mainKeeper        store.MainKeeper
	stakeIndex        *bank.IndexedKeeper
	cidNumKeeper      keeper.CidNumberKeeper
	linkIndexedKeeper keeper.LinkIndexedKeeper
}

func NewRankState(
	allowSearch bool, mainKeeper store.MainKeeper, stakeIndex *bank.IndexedKeeper,
	linkIndexedKeeper keeper.LinkIndexedKeeper, cidNumKeeper keeper.CidNumberKeeper,
	unit ComputeUnit,
) *RankState {
	return &RankState{
		allowSearch: allowSearch, indexErrChan: make(chan error), rankCalcChan: make(chan Rank, 1),
		rankErrChan: make(chan error), rankCalculationFinished: true, mainKeeper: mainKeeper,
		stakeIndex: stakeIndex, linkIndexedKeeper: linkIndexedKeeper, cidNumKeeper: cidNumKeeper,
		computeUnit: unit,
	}
}

func (s *RankState) Load(ctx sdk.Context, latestBlockHeight int64, log log.Logger) {
	s.lastCalculatedRankHash = s.mainKeeper.GetLastCalculatedRankHash(ctx)
	s.networkCidRank = Rank{Values: nil, MerkleTree: merkle.NewTree(sha256.New(), false)}
	s.nextCidRank = Rank{Values: nil, MerkleTree: merkle.NewTree(sha256.New(), false)}
	s.networkCidRank.MerkleTree.ImportSubtreesRoots(s.mainKeeper.GetLatestMerkleTree(ctx))
	s.nextCidRank.MerkleTree.ImportSubtreesRoots(s.mainKeeper.GetNextMerkleTree(ctx))

	s.cidCount = int64(s.mainKeeper.GetCidsCount(ctx))

	// if we have fallen and need to start new rank calculation
	// todo: what if rank didn't changed in new calculation???
	if latestBlockHeight != 0 && !s.nextRankReady() {
		s.startRankCalculation(ctx, log)
		s.rankCalculationFinished = false
	}
}

func (s *RankState) EndBlocker(ctx sdk.Context, log log.Logger) {
	currentCidsCount := s.mainKeeper.GetCidsCount(ctx)
	s.linkIndexedKeeper.EndBlocker()
	if ctx.BlockHeight()%CalculationPeriod == 0 || ctx.BlockHeight() == 1 {

		if !s.rankCalculationFinished {
			select {
			case newRank := <-s.rankCalcChan:
				s.handleNewRank(ctx, newRank)
			case err := <-s.rankErrChan:
				// DUMB ERROR HANDLING
				log.Error("Error during rank calculation " + err.Error())
				panic(err.Error())
			}
		}

		s.applyNextRank()
		// Recalculate index
		// todo state copied
		outLinksCopy := Links(s.linkIndexedKeeper.GetOutLinks()).Copy()
		go s.buildCidRankedLinksIndexInParallel(s.cidCount, outLinksCopy)

		s.mainKeeper.StoreLastCalculatedRankHash(ctx, s.GetNetworkRankHash())

		// start new calculation
		s.rankCalculationFinished = false
		s.cidCount = int64(currentCidsCount)
		s.linkIndexedKeeper.FixLinks()
		s.stakeIndex.FixUserStake()
		s.startRankCalculation(ctx, log)

	} else if !s.rankCalculationFinished {

		select {
		case newRank := <-s.rankCalcChan:
			s.handleNewRank(ctx, newRank)
		case err := <-s.rankErrChan:
			// DUMB ERROR HANDLING
			log.Error("Error during rank calculation " + err.Error())
			panic(err.Error())
		default:
		}

	}
	//CHECK INDEX BUILDING FOR ERROR:
	s.addNewCids(currentCidsCount)
	s.mainKeeper.StoreLatestMerkleTree(ctx, s.getNetworkMerkleTreeAsBytes())
	s.checkBuildIndexError(log)
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

func (s *RankState) startRankCalculation(ctx sdk.Context, log log.Logger) {
	calcCtx := NewCalcContext(ctx, s.linkIndexedKeeper, s.cidNumKeeper, s.stakeIndex, s.allowSearch)
	go CalculateRankInParallel(calcCtx, s.rankCalcChan, s.rankErrChan, s.computeUnit, log)
}

func (s *RankState) handleNewRank(ctx sdk.Context, newRank Rank) {
	s.setNextRank(newRank)
	s.mainKeeper.StoreNextMerkleTree(ctx, s.getNextMerkleTreeAsBytes())
	s.rankCalculationFinished = true
}

func (s *RankState) setNextRank(newRank Rank) {
	s.nextCidRank = newRank
}

func (s *RankState) applyNextRank() {
	s.networkCidRank = s.nextCidRank
	s.lastCalculatedRankHash = s.nextCidRank.MerkleTree.RootHash()
	s.nextCidRank.Reset()
}

// add new cids to rank with 0 value
func (s *RankState) addNewCids(cidCount uint64) {
	s.networkCidRank.AddNewCids(cidCount)
}

func (s *RankState) buildCidRankedLinksIndex(cidsCount int64, outLinks Links) {
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
func (s *RankState) buildCidRankedLinksIndexInParallel(cidsCount int64, outLinks Links) {
	defer func() {
		if r := recover(); r != nil {
			s.indexErrChan <- r.(error)
		}
	}()

	s.buildCidRankedLinksIndex(cidsCount, outLinks)
}

func (s *RankState) checkBuildIndexError(logger log.Logger) {
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

func (s *RankState) GetNetworkRankHash() []byte {
	return s.networkCidRank.MerkleTree.RootHash()
}

func (s *RankState) getNetworkMerkleTreeAsBytes() []byte {
	return s.networkCidRank.MerkleTree.ExportSubtreesRoots()
}

func (s *RankState) getNextMerkleTreeAsBytes() []byte {
	return s.nextCidRank.MerkleTree.ExportSubtreesRoots()
}

func (s *RankState) nextRankReady() bool {
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
