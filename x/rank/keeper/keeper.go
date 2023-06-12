package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cybercongress/go-cyber/merkle"
	graphkeeper "github.com/cybercongress/go-cyber/x/graph/keeper"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
	"github.com/cybercongress/go-cyber/x/rank/types"
	"github.com/tendermint/tendermint/libs/log"
)

type StateKeeper struct {
	networkCidRank types.Rank
	nextCidRank    types.Rank

	rankCalculationFinished bool
	cidCount                int64

	hasNewLinksForPeriod bool

	rankCalcChan chan types.Rank
	rankErrChan  chan error
	allowSearch  bool
	computeUnit  types.ComputeUnit

	stakeKeeper        types.StakeKeeper
	graphKeeper        types.GraphKeeper
	graphIndexedKeeper *graphkeeper.IndexKeeper
	accountKeeper      keeper.AccountKeeper

	storeKey   sdk.StoreKey
	paramSpace paramstypes.Subspace

	index         types.SearchIndex
	getIndexError types.GetError
}

func NewKeeper(
	key sdk.StoreKey,
	paramSpace paramstypes.Subspace,
	allowSearch bool,
	stakeIndex types.StakeKeeper,
	graphIndexedKeeper *graphkeeper.IndexKeeper,
	graphKeeper types.GraphKeeper,
	accountKeeper keeper.AccountKeeper,
	unit types.ComputeUnit,
) *StateKeeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &StateKeeper{
		storeKey:                key,
		paramSpace:              paramSpace,
		allowSearch:             allowSearch,
		rankCalcChan:            make(chan types.Rank, 1),
		rankErrChan:             make(chan error),
		rankCalculationFinished: true,
		stakeKeeper:             stakeIndex,
		graphIndexedKeeper:      graphIndexedKeeper,
		graphKeeper:             graphKeeper,
		accountKeeper:           accountKeeper,
		computeUnit:             unit,
		hasNewLinksForPeriod:    true,
	}
}

func (s *StateKeeper) LoadState(ctx sdk.Context) {
	s.networkCidRank = types.NewFromMerkle(s.graphKeeper.GetCidsCount(ctx), s.GetLatestMerkleTree(ctx))
	s.nextCidRank = types.NewFromMerkle(s.GetNextRankCidCount(ctx), s.GetNextMerkleTree(ctx))
	s.cidCount = int64(s.graphKeeper.GetCidsCount(ctx))

	s.index = s.BuildSearchIndex(s.Logger(ctx))
	s.index.Load(s.graphIndexedKeeper.GetOutLinks())
	s.getIndexError = s.index.Run()
}

func (s *StateKeeper) StartRankCalculation(ctx sdk.Context) {
	params := s.GetParams(ctx)

	dampingFactor, err := strconv.ParseFloat(params.DampingFactor.String(), 64)
	if err != nil {
		panic(err)
	}

	tolerance, err := strconv.ParseFloat(params.Tolerance.String(), 64)
	if err != nil {
		panic(err)
	}

	s.startRankCalculation(ctx, dampingFactor, tolerance)
	s.rankCalculationFinished = false
}

func (s *StateKeeper) BuildSearchIndex(logger log.Logger) types.SearchIndex {
	if s.allowSearch {
		return types.NewBaseSearchIndex(logger)
	} else {
		return types.NoopSearchIndex{}
	}
}

func (s *StateKeeper) EndBlocker(ctx sdk.Context) {
	s.StoreLatestBlockNumber(ctx, uint64(ctx.BlockHeight()))
	currentCidsCount := s.graphKeeper.GetCidsCount(ctx)

	s.index.PutNewLinks(s.graphIndexedKeeper.GetCurrentBlockNewLinks(ctx))
	// TODO MergeContextLinks need to be in graph's end blocker
	// but need to get new links here before nextRankLinks will be updated
	s.graphIndexedKeeper.MergeContextLinks(ctx)

	blockHasNewLinks := s.graphIndexedKeeper.HasNewLinks(ctx)
	s.hasNewLinksForPeriod = s.hasNewLinksForPeriod || blockHasNewLinks

	params := s.GetParams(ctx)

	if ctx.BlockHeight()%params.CalculationPeriod == 0 || ctx.BlockHeight() == 1 {
		dampingFactor, err := strconv.ParseFloat(params.DampingFactor.String(), 64)
		if err != nil {
			panic(err)
		}

		tolerance, err := strconv.ParseFloat(params.Tolerance.String(), 64)
		if err != nil {
			panic(err)
		}

		s.checkRankCalcFinished(ctx, true)
		s.applyNextRank(ctx)

		s.cidCount = int64(currentCidsCount)
		stakeChanged := s.stakeKeeper.DetectUsersStakeAmpereChange(ctx)

		// start new calculation
		if s.hasNewLinksForPeriod || stakeChanged {
			s.graphIndexedKeeper.UpdateRankLinks()
			s.graphKeeper.UpdateRankNeudegs()
			s.rankCalculationFinished = false
			s.hasNewLinksForPeriod = false
			s.prepareContext(ctx)
			s.startRankCalculation(ctx, dampingFactor, tolerance)
		}
	}

	s.networkCidRank.AddNewCids(currentCidsCount)
	networkMerkleTreeAsBytes := s.getNetworkMerkleTreeAsBytes()
	s.Logger(ctx).Info(
		"Latest Rank", "hash", fmt.Sprintf("%X", networkMerkleTreeAsBytes),
	)
	s.StoreLatestMerkleTree(ctx, networkMerkleTreeAsBytes)
}

func (s *StateKeeper) startRankCalculation(ctx sdk.Context, dampingFactor float64, tolerance float64) {
	calcCtx := types.NewCalcContext(
		s.graphIndexedKeeper, s.graphKeeper, s.stakeKeeper,
		s.allowSearch, dampingFactor, tolerance,
		s.GetContextCidCount(ctx),
		s.GetContextLinkCount(ctx),
		s.GetAccountCount(ctx),
	)

	go CalculateRankInParallel(calcCtx, s.rankCalcChan, s.rankErrChan, s.computeUnit, s.Logger(ctx))
}

func (s *StateKeeper) checkRankCalcFinished(ctx sdk.Context, block bool) {
	if !s.rankCalculationFinished {
		for {
			select {
			case newRank := <-s.rankCalcChan:
				s.handleNextRank(ctx, newRank)
				return
			case err := <-s.rankErrChan:
				s.Logger(ctx).Error("Error during cyber~Rank calculation, call cyber_devs! " + err.Error())
				panic(err.Error())
			default:
				if !block {
					return
				}
			}
			s.Logger(ctx).Info("Waiting for cyber~Rank calculation to finish")

			time.Sleep(2000 * time.Millisecond)
		}
	}
}

func (s *StateKeeper) handleNextRank(ctx sdk.Context, newRank types.Rank) {
	s.nextCidRank = newRank
	nextMerkleTreeAsBytes := s.getNextMerkleTreeAsBytes()
	s.Logger(ctx).Info(
		"Next Rank", "hash", fmt.Sprintf("%X", nextMerkleTreeAsBytes),
	)
	s.StoreNextMerkleTree(ctx, nextMerkleTreeAsBytes)
	s.StoreNextRankCidCount(ctx, newRank.CidCount)
	s.rankCalculationFinished = true
}

func (s *StateKeeper) applyNextRank(ctx sdk.Context) {
	if !s.nextCidRank.IsEmpty() {
		s.networkCidRank = s.nextCidRank
		s.index.PutNewRank(s.networkCidRank)
	}
	s.nextCidRank.Clear()
}

func (s *StateKeeper) GetRankValueByNumber(number uint64) uint64 {
	if number >= uint64(len(s.networkCidRank.RankValues)) {
		return 0
	}
	return s.networkCidRank.RankValues[number]
}

func (s *StateKeeper) GetRankValueByParticle(ctx sdk.Context, particle string) (uint64, error) {
	number, exist := s.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(particle))
	if exist != true {
		return 0, sdkerrors.ErrInvalidRequest
	}
	return s.networkCidRank.RankValues[number], nil
}

func (s *StateKeeper) GetNextNetworkRankHash() []byte {
	return s.nextCidRank.MerkleTree.RootHash()
}

func (s *StateKeeper) GetNetworkRankHash() []byte {
	return s.networkCidRank.MerkleTree.RootHash()
}

func (s *StateKeeper) getNetworkMerkleTreeAsBytes() []byte {
	return s.networkCidRank.MerkleTree.ExportSubtreesRoots()
}

func (s *StateKeeper) getNextMerkleTreeAsBytes() []byte {
	return s.nextCidRank.MerkleTree.ExportSubtreesRoots()
}

func (s *StateKeeper) GetLastCidNum() graphtypes.CidNumber {
	return graphtypes.CidNumber(len(s.networkCidRank.RankValues) - 1)
}

func (s *StateKeeper) GetMerkleTree() *merkle.Tree {
	return s.networkCidRank.MerkleTree
}

func (s *StateKeeper) GetKarma(accountNumber uint64) (karma uint64) {
	return s.networkCidRank.KarmaValues[accountNumber]
}

func (s *StateKeeper) GetEntropy(cidNum graphtypes.CidNumber) uint64 {
	return s.networkCidRank.EntropyValues[cidNum]
}

func (s *StateKeeper) GetNegEntropy() uint64 {
	return s.networkCidRank.NegEntropy
}

func (s *StateKeeper) GetIndexError() error {
	return s.getIndexError()
}

func (s *StateKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	s.paramSpace.GetParamSet(ctx, &params)
	return params
}

func (s *StateKeeper) SetParams(ctx sdk.Context, params types.Params) {
	s.paramSpace.SetParamSet(ctx, &params)
}

func (k *StateKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (sk StateKeeper) GetLatestBlockNumber(ctx sdk.Context) uint64 {
	store := ctx.KVStore(sk.storeKey)
	numberAsBytes := store.Get(types.LatestBlockNumber)
	if numberAsBytes == nil {
		return 0
	}
	return binary.LittleEndian.Uint64(numberAsBytes)
}

func (sk StateKeeper) StoreLatestBlockNumber(ctx sdk.Context, number uint64) {
	store := ctx.KVStore(sk.storeKey)
	numberAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(numberAsBytes, number)
	store.Set(types.LatestBlockNumber, numberAsBytes)
}

func (sk StateKeeper) GetLatestMerkleTree(ctx sdk.Context) []byte {
	store := ctx.KVStore(sk.storeKey)
	return store.Get(types.LatestMerkleTree)
}

func (sk StateKeeper) StoreLatestMerkleTree(ctx sdk.Context, treeAsBytes []byte) {
	store := ctx.KVStore(sk.storeKey)
	if bytes.Compare(sk.GetLatestMerkleTree(ctx), treeAsBytes) != 0 {
		store.Set(types.LatestMerkleTree, treeAsBytes)
	}
}

func (sk StateKeeper) GetNextMerkleTree(ctx sdk.Context) []byte {
	store := ctx.KVStore(sk.storeKey)
	return store.Get(types.NextMerkleTree)
}

func (sk StateKeeper) StoreNextMerkleTree(ctx sdk.Context, treeAsBytes []byte) {
	store := ctx.KVStore(sk.storeKey)
	if bytes.Compare(sk.GetNextMerkleTree(ctx), treeAsBytes) != 0 {
		store.Set(types.NextMerkleTree, treeAsBytes)
	}
}

func (sk StateKeeper) GetNextRankCidCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(sk.storeKey)
	numberAsBytes := store.Get(types.NextRankCidCount)
	if numberAsBytes == nil {
		return sk.graphKeeper.GetCidsCount(ctx)
	}
	return binary.LittleEndian.Uint64(numberAsBytes)
}

func (sk StateKeeper) StoreNextRankCidCount(ctx sdk.Context, number uint64) {
	if sk.GetNextRankCidCount(ctx) != number {
		store := ctx.KVStore(sk.storeKey)
		numberAsBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(numberAsBytes, number)
		store.Set(types.NextRankCidCount, numberAsBytes)
	}
}

func (sk StateKeeper) GetContextCidCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(sk.storeKey)
	numberAsBytes := store.Get(types.ContextCidCount)
	if numberAsBytes == nil {
		return 0
	}
	return binary.LittleEndian.Uint64(numberAsBytes)
}

func (sk StateKeeper) StoreContextCidCount(ctx sdk.Context, number uint64) {
	if sk.GetContextCidCount(ctx) != number {
		store := ctx.KVStore(sk.storeKey)
		numberAsBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(numberAsBytes, number)
		store.Set(types.ContextCidCount, numberAsBytes)
	}
}

func (sk StateKeeper) GetContextLinkCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(sk.storeKey)
	numberAsBytes := store.Get(types.ContextLinkCount)
	if numberAsBytes == nil {
		return 0
	}
	return binary.LittleEndian.Uint64(numberAsBytes)
}

func (sk StateKeeper) StoreContextLinkCount(ctx sdk.Context, number uint64) {
	if sk.GetContextLinkCount(ctx) != number {
		store := ctx.KVStore(sk.storeKey)
		numberAsBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(numberAsBytes, number)
		store.Set(types.ContextLinkCount, numberAsBytes)
	}
}

func (s *StateKeeper) prepareContext(ctx sdk.Context) {
	s.StoreContextCidCount(ctx, s.graphKeeper.GetCidsCount(ctx))
	s.StoreContextLinkCount(ctx, s.graphIndexedKeeper.GetLinksCount(ctx))
}

func (s *StateKeeper) GetAccountCount(ctx sdk.Context) uint64 {
	return s.stakeKeeper.GetNextAccountNumber(ctx)
}
