package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/x/auth/keeper"

	"github.com/cybercongress/go-cyber/v5/merkle"
	graphkeeper "github.com/cybercongress/go-cyber/v5/x/graph/keeper"
	graphtypes "github.com/cybercongress/go-cyber/v5/x/graph/types"
	"github.com/cybercongress/go-cyber/v5/x/rank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cometbft/cometbft/libs/log"
)

type StateKeeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	networkCidRank types.Rank
	nextCidRank    types.Rank

	rankCalculationFinished bool
	//cidCount                int64

	hasNewLinksForPeriod bool

	rankCalcChan chan types.Rank
	rankErrChan  chan error
	allowSearch  bool
	computeUnit  types.ComputeUnit

	stakeKeeper        types.StakeKeeper
	graphKeeper        types.GraphKeeper
	graphIndexedKeeper *graphkeeper.IndexKeeper
	accountKeeper      keeper.AccountKeeper

	index         types.SearchIndex
	getIndexError types.GetError

	authority string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	allowSearch bool,
	stakeIndex types.StakeKeeper,
	graphIndexedKeeper *graphkeeper.IndexKeeper,
	graphKeeper types.GraphKeeper,
	accountKeeper keeper.AccountKeeper,
	unit types.ComputeUnit,
	authority string,
) *StateKeeper {
	return &StateKeeper{
		storeKey:                key,
		cdc:                     cdc,
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
		authority:               authority,
	}
}

func (sk *StateKeeper) GetAuthority() string { return sk.authority }

func (sk *StateKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k *StateKeeper) SetParams(ctx sdk.Context, p types.Params) error {
	if err := p.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&p)
	store.Set(types.ParamsKey, bz)

	return nil
}

func (k *StateKeeper) GetParams(ctx sdk.Context) (p types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return p
	}

	k.cdc.MustUnmarshal(bz, &p)
	return p
}

func (sk *StateKeeper) LoadState(ctx sdk.Context) {
	sk.networkCidRank = types.NewFromMerkle(sk.graphKeeper.GetCidsCount(ctx), sk.GetLatestMerkleTree(ctx))
	sk.nextCidRank = types.NewFromMerkle(sk.GetNextRankCidCount(ctx), sk.GetNextMerkleTree(ctx))

	sk.index = sk.BuildSearchIndex(sk.Logger(ctx))
	sk.index.Load(sk.graphIndexedKeeper.GetOutLinks())
	sk.getIndexError = sk.index.Run()
}

func (sk *StateKeeper) StartRankCalculation(ctx sdk.Context) {
	params := sk.GetParams(ctx)

	// TODO remove this after upgrade to v4 because on network upgrade block cannot access rank params
	if params.CalculationPeriod == 0 {
		params = types.Params{
			CalculationPeriod: 5,
			DampingFactor:     sdk.NewDecWithPrec(85, 2),
			Tolerance:         sdk.NewDecWithPrec(1, 3),
		}
	}

	dampingFactor, err := strconv.ParseFloat(params.DampingFactor.String(), 64)
	if err != nil {
		panic(err)
	}

	tolerance, err := strconv.ParseFloat(params.Tolerance.String(), 64)
	if err != nil {
		panic(err)
	}

	sk.startRankCalculation(ctx, dampingFactor, tolerance)
	sk.rankCalculationFinished = false
}

func (sk *StateKeeper) BuildSearchIndex(logger log.Logger) types.SearchIndex {
	if sk.allowSearch {
		return types.NewBaseSearchIndex(logger)
	}
	return types.NoopSearchIndex{}
}

func (sk *StateKeeper) EndBlocker(ctx sdk.Context) {
	sk.StoreLatestBlockNumber(ctx, uint64(ctx.BlockHeight()))
	currentCidsCount := sk.graphKeeper.GetCidsCount(ctx)

	sk.index.PutNewLinks(sk.graphIndexedKeeper.GetCurrentBlockNewLinks(ctx))
	// TODO MergeContextLinks need to be in graph'sk end blocker
	// but need to get new links here before nextRankLinks will be updated
	sk.graphIndexedKeeper.MergeContextLinks(ctx)

	blockHasNewLinks := sk.graphIndexedKeeper.HasNewLinks(ctx)
	sk.hasNewLinksForPeriod = sk.hasNewLinksForPeriod || blockHasNewLinks

	params := sk.GetParams(ctx)

	if ctx.BlockHeight()%params.CalculationPeriod == 0 || ctx.BlockHeight() == 1 {
		dampingFactor, err := strconv.ParseFloat(params.DampingFactor.String(), 64)
		if err != nil {
			panic(err)
		}

		tolerance, err := strconv.ParseFloat(params.Tolerance.String(), 64)
		if err != nil {
			panic(err)
		}

		sk.checkRankCalcFinished(ctx, true)
		sk.applyNextRank()
		stakeChanged := sk.stakeKeeper.DetectUsersStakeAmpereChange(ctx)

		// start new calculation
		if sk.hasNewLinksForPeriod || stakeChanged {
			sk.graphIndexedKeeper.UpdateRankLinks()
			sk.graphKeeper.UpdateRankNeudegs()
			sk.rankCalculationFinished = false
			sk.hasNewLinksForPeriod = false
			sk.prepareContext(ctx)
			sk.startRankCalculation(ctx, dampingFactor, tolerance)
		}
	}

	sk.networkCidRank.AddNewCids(currentCidsCount)
	networkMerkleTreeAsBytes := sk.getNetworkMerkleTreeAsBytes()
	sk.Logger(ctx).Info(
		"Latest Rank", "hash", fmt.Sprintf("%X", networkMerkleTreeAsBytes),
	)

	sk.StoreLatestMerkleTree(ctx, networkMerkleTreeAsBytes)
}

func (sk *StateKeeper) startRankCalculation(ctx sdk.Context, dampingFactor float64, tolerance float64) {
	calcCtx := types.NewCalcContext(
		sk.graphIndexedKeeper, sk.graphKeeper, sk.stakeKeeper,
		sk.allowSearch, dampingFactor, tolerance,
		sk.GetContextCidCount(ctx),
		sk.GetContextLinkCount(ctx),
		sk.GetAccountCount(ctx),
	)

	go CalculateRankInParallel(calcCtx, sk.rankCalcChan, sk.rankErrChan, sk.computeUnit, sk.Logger(ctx))
}

func (sk *StateKeeper) checkRankCalcFinished(ctx sdk.Context, block bool) {
	if !sk.rankCalculationFinished {
		for {
			select {
			case newRank := <-sk.rankCalcChan:
				sk.handleNextRank(ctx, newRank)
				return
			case err := <-sk.rankErrChan:
				sk.Logger(ctx).Error("Error during cyber~Rank calculation, call cyber_devs! " + err.Error())
				panic(err.Error())
			default:
				if !block {
					return
				}
			}
			sk.Logger(ctx).Info("Waiting for cyber~Rank calculation to finish")

			time.Sleep(2000 * time.Millisecond)
		}
	}
}

func (sk *StateKeeper) handleNextRank(ctx sdk.Context, newRank types.Rank) {
	sk.nextCidRank = newRank
	nextMerkleTreeAsBytes := sk.getNextMerkleTreeAsBytes()
	sk.Logger(ctx).Info(
		"Next Rank", "hash", fmt.Sprintf("%X", nextMerkleTreeAsBytes),
	)
	sk.StoreNextMerkleTree(ctx, nextMerkleTreeAsBytes)
	sk.StoreNextRankCidCount(ctx, newRank.CidCount)
	sk.rankCalculationFinished = true
}

func (sk *StateKeeper) applyNextRank() {
	if !sk.nextCidRank.IsEmpty() {
		sk.networkCidRank = sk.nextCidRank
		sk.index.PutNewRank(sk.networkCidRank)
	}
	sk.nextCidRank.Clear()
}

func (sk *StateKeeper) GetRankValueByNumber(number uint64) uint64 {
	if number >= uint64(len(sk.networkCidRank.RankValues)) {
		return 0
	}
	return sk.networkCidRank.RankValues[number]
}

func (sk *StateKeeper) GetRankValueByParticle(ctx sdk.Context, particle string) (uint64, error) {
	number, exist := sk.graphKeeper.GetCidNumber(ctx, graphtypes.Cid(particle))
	if !exist {
		return 0, sdkerrors.ErrInvalidRequest
	}
	return sk.networkCidRank.RankValues[number], nil
}

func (sk *StateKeeper) GetNextNetworkRankHash() []byte {
	return sk.nextCidRank.MerkleTree.RootHash()
}

func (sk *StateKeeper) GetNetworkRankHash() []byte {
	return sk.networkCidRank.MerkleTree.RootHash()
}

func (sk *StateKeeper) getNetworkMerkleTreeAsBytes() []byte {
	return sk.networkCidRank.MerkleTree.ExportSubtreesRoots()
}

func (sk *StateKeeper) getNextMerkleTreeAsBytes() []byte {
	return sk.nextCidRank.MerkleTree.ExportSubtreesRoots()
}

func (sk *StateKeeper) GetLastCidNum() graphtypes.CidNumber {
	return graphtypes.CidNumber(len(sk.networkCidRank.RankValues) - 1)
}

func (sk *StateKeeper) GetMerkleTree() *merkle.Tree {
	return sk.networkCidRank.MerkleTree
}

func (sk *StateKeeper) GetKarma(accountNumber uint64) (karma uint64) {
	return sk.networkCidRank.KarmaValues[accountNumber]
}

func (sk *StateKeeper) GetEntropy(cidNum graphtypes.CidNumber) uint64 {
	return sk.networkCidRank.EntropyValues[cidNum]
}

func (sk *StateKeeper) GetNegEntropy() uint64 {
	return sk.networkCidRank.NegEntropy
}

func (sk *StateKeeper) GetIndexError() error {
	return sk.getIndexError()
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
	if bytes.Compare(sk.GetLatestMerkleTree(ctx), treeAsBytes) != 0 { //nolint:gosimple
		store.Set(types.LatestMerkleTree, treeAsBytes)
	}
}

func (sk StateKeeper) GetNextMerkleTree(ctx sdk.Context) []byte {
	store := ctx.KVStore(sk.storeKey)
	return store.Get(types.NextMerkleTree)
}

func (sk StateKeeper) StoreNextMerkleTree(ctx sdk.Context, treeAsBytes []byte) {
	store := ctx.KVStore(sk.storeKey)
	if bytes.Compare(sk.GetNextMerkleTree(ctx), treeAsBytes) != 0 { //nolint:gosimple
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

func (sk *StateKeeper) prepareContext(ctx sdk.Context) {
	sk.StoreContextCidCount(ctx, sk.graphKeeper.GetCidsCount(ctx))
	sk.StoreContextLinkCount(ctx, sk.graphIndexedKeeper.GetLinksCount(ctx))
}

func (sk *StateKeeper) GetAccountCount(ctx sdk.Context) uint64 {
	return sk.stakeKeeper.GetNextAccountNumber(ctx)
}
