package rank

import (
	"crypto/sha256"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/merkle"
	"github.com/cybercongress/cyberd/store"
	"github.com/cybercongress/cyberd/x/bank"
	"github.com/cybercongress/cyberd/x/link/keeper"
	. "github.com/cybercongress/cyberd/x/link/types"
	"github.com/tendermint/tendermint/libs/log"
	"time"
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
	index *SearchIndex
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
	s.networkCidRank = Rank{
		Values:     nil,
		MerkleTree: merkle.NewTree(sha256.New(), false),
		CidCount:   s.mainKeeper.GetCidsCount(ctx),
	}
	s.nextCidRank = Rank{
		Values:     nil,
		MerkleTree: merkle.NewTree(sha256.New(), false),
		CidCount:   s.mainKeeper.GetNextRankCidCount(ctx),
	}
	s.networkCidRank.MerkleTree.ImportSubtreesRoots(s.mainKeeper.GetLatestMerkleTree(ctx))
	s.nextCidRank.MerkleTree.ImportSubtreesRoots(s.mainKeeper.GetNextMerkleTree(ctx))

	s.cidCount = int64(s.mainKeeper.GetCidsCount(ctx))

	if s.allowSearch {
		s.index = NewSearchIndex(log)
		s.index.Load(s.linkIndexedKeeper)

		go s.index.startListenNewLinks()
		go s.index.startListenNewRank()
	}

	// if we fell down and need to start new rank calculation
	if !s.mainKeeper.GetRankCalculationFinished(ctx) {
		s.startRankCalculation(ctx, log)
		s.rankCalculationFinished = false
	}
}

func (s *RankState) EndBlocker(ctx sdk.Context, log log.Logger) {
	currentCidsCount := s.mainKeeper.GetCidsCount(ctx)

	if s.allowSearch {
		s.submitNewLinksToIndex()
		s.checkIndexFailure(log)
	}

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
	// todo: hardcode. remove after next relaunch
	if ctx.BlockHeight() >= CalculationPeriod  {
		s.addNewCids(currentCidsCount)
	}
	s.mainKeeper.StoreLatestMerkleTree(ctx, s.getNetworkMerkleTreeAsBytes())
}

func (s *RankState) Search(cidNumber CidNumber, page, perPage int) ([]RankedCidNumber, int, error) {
	if s.allowSearch {
		return s.index.Search(cidNumber, page, perPage)
	}
	return nil, 0, errors.New("search is not allowed on this node")
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
		if s.allowSearch {
			s.index.RankChan <- s.networkCidRank.CopyWithoutTree()
		}
	}
	s.nextCidRank.Clear()
}

// add new cids to rank with 0 value
func (s *RankState) addNewCids(cidCount uint64) {
	s.networkCidRank.AddNewCids(cidCount)
}

//
// SEARCH INDEX METHODS
//

func (s *RankState) submitNewLinksToIndex() {
	for _, link := range s.linkIndexedKeeper.GetCurrentBlockLinks() {
		if !s.linkIndexedKeeper.IsAnyLinkExist(link.From(), link.To()) {
			s.index.LinksChan <- link
		}
	}
}

func (s *RankState) checkIndexFailure(logger log.Logger) {
	// Check search index for error
	err := s.index.checkIndexError()
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
