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
)

type RankState struct {
	// todo: remove rank Values, store only merkle tree (or even root hash). Move values to link_index
	networkCidRank Rank // array linksIndex is cid number
	nextCidRank    Rank // array linksIndex is cid number

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

	if s.allowSearch {
		s.index = NewSearchIndex(log)
		s.index.Load(s.linkIndexedKeeper)

		go s.index.startListenNewLinks()
		go s.index.startListenNewRank()
	}

	// if we have fallen and need to start new rank calculation
	// todo: what if rank didn't changed in new calculation???
	if latestBlockHeight != 0 && !s.nextRankReady() {
		s.startRankCalculation(ctx, log)
		s.rankCalculationFinished = false
	}
}

func (s *RankState) EndBlocker(ctx sdk.Context, log log.Logger) {
	currentCidsCount := s.mainKeeper.GetCidsCount(ctx)

	if s.allowSearch {
		for _, link := range s.linkIndexedKeeper.GetCurrentBlockLinks() {
			outLinks := s.linkIndexedKeeper.GetNextOutLinks()
			toCids, exists := outLinks[link.From()]
			if exists {
				_, exists := toCids[link.To()]
				if exists {
					continue
				}
			}
			s.index.LinksChan <- link
		}
	}

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
	s.addNewCids(currentCidsCount)
	s.mainKeeper.StoreLatestMerkleTree(ctx, s.getNetworkMerkleTreeAsBytes())

	if s.allowSearch {
		// Check search index for error
		err := s.index.checkIndexError()
		if err != nil {
			log.Error("Search index failed!", "error", err)
			panic(err)
		}
	}
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
	if s.allowSearch {
		s.index.RankChan <- s.networkCidRank
	}
}

// add new cids to rank with 0 value
func (s *RankState) addNewCids(cidCount uint64) {
	s.networkCidRank.AddNewCids(cidCount)
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
