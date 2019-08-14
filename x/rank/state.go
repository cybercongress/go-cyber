package rank

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/store"
	"github.com/cybercongress/cyberd/x/bank"
	"github.com/cybercongress/cyberd/x/link/keeper"
	. "github.com/cybercongress/cyberd/x/link/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/cybercongress/cyberd/merkle"
)

type RankState struct {
	networkCidRank Rank // array linksIndex is cid number
	nextCidRank    Rank // array linksIndex is cid number

	rankCalculationFinished bool
	cidCount                int64

	hasNewLinksForPeriod bool // indicates if new links where submitted for rank calc period

	rankCalcChan chan Rank
	rankErrChan  chan error
	allowSearch  bool
	computeUnit  ComputeUnit

	// keepers
	mainKeeper        store.MainKeeper
	stakeIndex        *bank.IndexedKeeper
	cidNumKeeper      keeper.CidNumberKeeper
	linkIndexedKeeper *keeper.LinkIndexedKeeper

	// index
	index         SearchIndex
	getIndexError GetError
}

func NewRankState(
	allowSearch bool, mainKeeper store.MainKeeper, stakeIndex *bank.IndexedKeeper,
	linkIndexedKeeper *keeper.LinkIndexedKeeper, cidNumKeeper keeper.CidNumberKeeper,
	unit ComputeUnit,
) *RankState {
	return &RankState{
		allowSearch: allowSearch, rankCalcChan: make(chan Rank, 1),
		rankErrChan: make(chan error), rankCalculationFinished: true, mainKeeper: mainKeeper,
		stakeIndex: stakeIndex, linkIndexedKeeper: linkIndexedKeeper, cidNumKeeper: cidNumKeeper,
		computeUnit: unit, hasNewLinksForPeriod: true,
	}
}

func (s *RankState) Load(ctx sdk.Context, log log.Logger) {
	s.networkCidRank = NewFromMerkle(s.mainKeeper.GetCidsCount(ctx), s.mainKeeper.GetLatestMerkleTree(ctx))
	s.nextCidRank = NewFromMerkle(s.mainKeeper.GetNextRankCidCount(ctx), s.mainKeeper.GetNextMerkleTree(ctx))
	s.cidCount = int64(s.mainKeeper.GetCidsCount(ctx))

	s.index = s.BuildSearchIndex(log)
	s.index.Load(s.linkIndexedKeeper.GetOutLinks())
	s.getIndexError = s.index.Run()

	// if we fell down and need to start new rank calculation
	if !s.mainKeeper.GetRankCalculationFinished(ctx) {
		s.startRankCalculation(ctx, log)
		s.rankCalculationFinished = false
	}
}

func (s *RankState) BuildSearchIndex(logger log.Logger) SearchIndex {
	if s.allowSearch {
		return NewBaseSearchIndex(logger)
	} else {
		return NoopSearchIndex{}
	}
}

func (s *RankState) EndBlocker(ctx sdk.Context, log log.Logger) {
	currentCidsCount := s.mainKeeper.GetCidsCount(ctx)

	s.index.PutNewLinks(s.linkIndexedKeeper.GetCurrentBlockNewLinks())
	s.checkIndexFailure(log)

	blockHasNewLinks := s.linkIndexedKeeper.EndBlocker()
	s.hasNewLinksForPeriod = s.hasNewLinksForPeriod || blockHasNewLinks
	if ctx.BlockHeight()%CalculationPeriod == 0 || ctx.BlockHeight() == 1 {

		s.checkRankCalcFinished(ctx, true, log)
		s.applyNextRank()

		s.cidCount = int64(currentCidsCount)
		stakeChanged := s.stakeIndex.FixUserStake()

		// start new calculation
		if s.hasNewLinksForPeriod || stakeChanged {
			s.linkIndexedKeeper.FixLinks()
			s.rankCalculationFinished = false
			s.hasNewLinksForPeriod = false
			s.mainKeeper.StoreRankCalculationFinished(ctx, false)
			s.startRankCalculation(ctx, log)
		}
	} else {
		s.checkRankCalcFinished(ctx, false, log)
	}
	s.networkCidRank.AddNewCids(currentCidsCount)
	s.mainKeeper.StoreLatestMerkleTree(ctx, s.getNetworkMerkleTreeAsBytes())
}

func (s *RankState) Search(cidNumber CidNumber, page, perPage int) ([]RankedCidNumber, int, error) {
	return s.index.Search(cidNumber, page, perPage)
}

func (s *RankState) GetRankValue(cidNumber CidNumber) float64 {
	return s.index.GetRankValue(cidNumber)
}

func (s *RankState) startRankCalculation(ctx sdk.Context, log log.Logger) {
	calcCtx := NewCalcContext(ctx, s.linkIndexedKeeper, s.cidNumKeeper, s.stakeIndex, s.allowSearch)
	go CalculateRankInParallel(calcCtx, s.rankCalcChan, s.rankErrChan, s.computeUnit, log)
}

func (s *RankState) checkRankCalcFinished(ctx sdk.Context, block bool, log log.Logger) {

	if !s.rankCalculationFinished {
		for {
			select {
			case newRank := <-s.rankCalcChan:
				s.handleNextRank(ctx, newRank)
				return
			case err := <-s.rankErrChan:
				// DUMB ERROR HANDLING
				log.Error("Error during rank calculation " + err.Error())
				panic(err.Error())
			default:
				if !block {
					return
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (s *RankState) handleNextRank(ctx sdk.Context, newRank Rank) {
	s.nextCidRank = newRank
	s.mainKeeper.StoreNextMerkleTree(ctx, s.getNextMerkleTreeAsBytes())
	s.mainKeeper.StoreNextRankCidCount(ctx, newRank.CidCount)
	s.rankCalculationFinished = true
	s.mainKeeper.StoreRankCalculationFinished(ctx, true)
}

func (s *RankState) applyNextRank() {

	if !s.nextCidRank.IsEmpty() {
		s.networkCidRank = s.nextCidRank
		s.index.PutNewRank(s.networkCidRank)
	}
	s.nextCidRank.Clear()
}

func (s *RankState) checkIndexFailure(logger log.Logger) {
	// Check search index for error
	err := s.getIndexError()
	if err != nil {
		logger.Error("Search index failed!", "error", err)
		panic(err)
	}
}

//
// GETTERS
//

func (s *RankState) GetNetworkRankHash() []byte {
	return s.networkCidRank.MerkleTree.RootHash()
}

func (s *RankState) getNetworkMerkleTreeAsBytes() []byte {
	return s.networkCidRank.MerkleTree.ExportSubtreesRoots()
}

func (s *RankState) getNextMerkleTreeAsBytes() []byte {
	return s.nextCidRank.MerkleTree.ExportSubtreesRoots()
}

func (s *RankState) GetLastCidNum() CidNumber {
	return CidNumber(len(s.networkCidRank.Values) - 1)
}

func (s *RankState) GetMerkleTree() *merkle.Tree {
	return s.networkCidRank.MerkleTree
}