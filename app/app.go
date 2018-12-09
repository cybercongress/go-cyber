package app

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	. "github.com/cybercongress/cyberd/app/bank"
	"github.com/cybercongress/cyberd/app/coin"
	"github.com/cybercongress/cyberd/app/rank"
	. "github.com/cybercongress/cyberd/app/storage"
	cbd "github.com/cybercongress/cyberd/app/types"
	"github.com/cybercongress/cyberd/types"
	"github.com/cybercongress/cyberd/x/bandwidth"
	"github.com/cybercongress/cyberd/x/mint"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"math"
	"os"
	"sort"
	"time"
)

const (
	APP            = "cyberd"
	appName        = "CyberdApp"
	DefaultKeyPass = "12345678"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.cyberdcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.cyberd")
)

type CyberdAppDbKeys struct {
	main         *sdk.KVStoreKey
	acc          *sdk.KVStoreKey
	accIndex     *sdk.KVStoreKey
	cidIndex     *sdk.KVStoreKey
	links        *sdk.KVStoreKey
	rank         *sdk.KVStoreKey
	stake        *sdk.KVStoreKey
	tStake       *sdk.TransientStoreKey
	fees         *sdk.KVStoreKey
	keyDistr     *sdk.KVStoreKey
	slashing     *sdk.KVStoreKey
	params       *sdk.KVStoreKey
	tParams      *sdk.TransientStoreKey
	accBandwidth *sdk.KVStoreKey
}

// CyberdApp implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore dbKeys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type CyberdApp struct {
	*baseapp.BaseApp
	cdc *codec.Codec

	txDecoder        sdk.TxDecoder
	bandwidthHandler types.BandwidthHandler

	// keys to access the multistore
	dbKeys CyberdAppDbKeys

	// manage getting and setting accounts
	mainStorage         MainStorage
	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          bank.Keeper
	stakeKeeper         stake.Keeper
	slashingKeeper      slashing.Keeper
	distrKeeper         distr.Keeper
	paramsKeeper        params.Keeper
	accBandwidthKeeper  bandwidth.AccountBandwidthKeeper

	//inflation
	minter mint.Minter

	// cyberd storages
	persistStorages CyberdPersistentStorages
	memStorage      *InMemoryStorage

	latestRankHash     []byte
	latestBlockHeight  int64
	currentCreditPrice float64

	computeUnit rank.ComputeUnit
}

// NewBasecoinApp returns a reference to a new CyberdApp given a
// logger and
// database. Internally, a codec is created along with all the necessary dbKeys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewCyberdApp(
	logger log.Logger, db dbm.DB, computeUnit rank.ComputeUnit, baseAppOptions ...func(*baseapp.BaseApp),
) *CyberdApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()

	dbKeys := CyberdAppDbKeys{
		main:     sdk.NewKVStoreKey("main"),
		acc:      sdk.NewKVStoreKey("acc"),
		cidIndex: sdk.NewKVStoreKey("cid_index"),
		links:    sdk.NewKVStoreKey("links"),
		rank:     sdk.NewKVStoreKey("rank"),
		stake:    sdk.NewKVStoreKey("stake"),
		fees:     sdk.NewKVStoreKey("fee"),
		tStake:   sdk.NewTransientStoreKey("transient_stake"),
		keyDistr: sdk.NewKVStoreKey("distr"),
		slashing: sdk.NewKVStoreKey("slashing"),
		params:   sdk.NewKVStoreKey("params"),
		tParams:  sdk.NewTransientStoreKey("transient_params"),
		accBandwidth: sdk.NewKVStoreKey("acc_bandwidth"),
	}

	ms := NewMainStorage(dbKeys.main)
	storages := CyberdPersistentStorages{
		CidIndex: NewCidIndexStorage(ms, dbKeys.cidIndex),
		Links:    NewLinksStorage(ms, dbKeys.links),
		Rank:     NewRankStorage(ms, dbKeys.rank),
	}

	var txDecoder = auth.DefaultTxDecoder(cdc)

	// create your application type
	var app = &CyberdApp{
		cdc:             cdc,
		txDecoder:       txDecoder,
		BaseApp:         baseapp.NewBaseApp(appName, logger, db, txDecoder, baseAppOptions...),
		dbKeys:          dbKeys,
		persistStorages: storages,
		mainStorage:     ms,
		computeUnit:     computeUnit,
	}

	// define and attach the mappers and keepers
	app.accountKeeper = auth.NewAccountKeeper(app.cdc, dbKeys.acc, NewAccount)
	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper)
	app.paramsKeeper = params.NewKeeper(app.cdc, dbKeys.params, dbKeys.tParams)
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(app.cdc, dbKeys.fees)
	app.accBandwidthKeeper = bandwidth.NewAccountBandwidthKeeper(dbKeys.accBandwidth)

	stakeKeeper := stake.NewKeeper(
		app.cdc, dbKeys.stake,
		dbKeys.tStake, app.bankKeeper,
		app.paramsKeeper.Subspace(stake.DefaultParamspace),
		stake.DefaultCodespace,
	)
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc, dbKeys.slashing,
		&stakeKeeper, app.paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)
	app.distrKeeper = distr.NewKeeper(
		app.cdc, dbKeys.keyDistr,
		app.paramsKeeper.Subspace(distr.DefaultParamspace),
		app.bankKeeper, &stakeKeeper, app.feeCollectionKeeper,
		distr.DefaultCodespace,
	)
	app.minter = mint.NewMinter(app.feeCollectionKeeper, stakeKeeper)

	app.memStorage = &InMemoryStorage{}

	// register the staking hooks
	// NOTE: stakeKeeper above are passed by reference,
	// so that it can be modified like below:
	app.stakeKeeper = *stakeKeeper.SetHooks(NewHooks(app.slashingKeeper.Hooks()))

	app.bandwidthHandler = bandwidth.NewBandwidthHandler(app.stakeKeeper, app.accountKeeper, app.accBandwidthKeeper)

	// register message routes
	app.Router().
		AddRoute("bank", NewBankHandler(app.bankKeeper, app.memStorage, app.accountKeeper)).
		AddRoute("link", NewLinksHandler(storages.CidIndex, storages.Links, app.memStorage, app.accountKeeper)).
		AddRoute("stake", stake.NewHandler(app.stakeKeeper)).
		AddRoute("slashing", slashing.NewHandler(app.slashingKeeper))

	app.QueryRouter().
		AddRoute("stake", stake.NewQuerier(app.stakeKeeper, app.cdc))

	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))

	// mount the multistore and load the latest state
	app.MountStores(
		dbKeys.main, dbKeys.acc, dbKeys.cidIndex, dbKeys.links, dbKeys.rank, dbKeys.stake,
		dbKeys.slashing, dbKeys.params, dbKeys.keyDistr, dbKeys.fees,
	)
	app.MountStoresTransient(dbKeys.tParams, dbKeys.tStake)
	err := app.LoadLatestVersion(dbKeys.main)
	if err != nil {
		cmn.Exit(err.Error())
	}
	ctx := app.BaseApp.NewContext(true, abci.Header{})

	// load in-memory data
	start := time.Now()
	app.BaseApp.Logger.Info("Loading mem state")
	app.memStorage.Load(ctx, storages, app.accountKeeper)
	app.BaseApp.Logger.Info("App loaded", "time", time.Since(start))
	app.latestRankHash = ms.GetAppHash(ctx)
	app.Seal()
	return app
}

// initChainer implements the custom application logic that the BaseApp will
// invoke upon initialization. In this case, it will take the application's
// state provided by 'req' and attempt to deserialize said state. The state
// should contain all the genesis accounts. These accounts will be added to the
// application's account mapper.
func (app *CyberdApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {

	stateJSON := req.AppStateBytes
	var genesisState GenesisState
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err)
	}

	// sort by account number to maintain consistency
	sort.Slice(genesisState.Accounts, func(i, j int) bool {
		return genesisState.Accounts[i].AccountNumber < genesisState.Accounts[j].AccountNumber
	})
	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, acc)
		app.memStorage.UpdateStake(cbd.AccountNumber(acc.AccountNumber), acc.Coins.AmountOf(coin.CBD).Int64())
	}

	// load the initial stake information
	validators, err := stake.InitGenesis(ctx, app.stakeKeeper, genesisState.StakeData)
	if err != nil {
		panic(err)
	}

	// initialize module-specific stores
	slashing.InitGenesis(ctx, app.slashingKeeper, genesisState.SlashingData, genesisState.StakeData)
	distr.InitGenesis(ctx, app.distrKeeper, genesisState.DistrData)
	err = CyberdValidateGenesisState(genesisState)
	if err != nil {
		panic(err)
	}

	// deliver genesis txes
	if len(genesisState.GenTxs) > 0 {
		for _, genTx := range genesisState.GenTxs {
			var tx auth.StdTx
			err = app.cdc.UnmarshalJSON(genTx, &tx)
			if err != nil {
				panic(err)
			}
			bz := app.cdc.MustMarshalBinaryLengthPrefixed(tx)
			res := app.BaseApp.DeliverTx(bz)
			if !res.IsOK() {
				panic(res.Log)
			}
		}

		validators = app.stakeKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}

	// sanity check
	if len(req.Validators) > 0 {
		if len(req.Validators) != len(validators) {
			panic(fmt.Errorf("len(RequestInitChain.Validators) != len(validators) (%d != %d)",
				len(req.Validators), len(validators)))
		}
		sort.Sort(abci.ValidatorUpdates(req.Validators))
		sort.Sort(abci.ValidatorUpdates(validators))
		for i, val := range validators {
			if !val.Equal(req.Validators[i]) {
				panic(fmt.Errorf("validators[%d] != req.Validators[%d] ", i, i))
			}
		}
	}

	return abci.ResponseInitChain{
		Validators: validators,
	}
}

// BeginBlocker reflects logic to run before any TXs application are processed
// by the application.
func (app *CyberdApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	mint.BeginBlocker(ctx, app.minter)
	tags := slashing.BeginBlocker(ctx, req, app.slashingKeeper)

	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
}

func (app *CyberdApp) CheckTx(txBytes []byte) (res abci.ResponseCheckTx) {

	cachedCtx, writeCache := app.NewContext(true, abci.Header{Height: app.latestBlockHeight}).CacheContext()

	var tx, err = app.txDecoder(txBytes)
	if err == nil {

		err = app.bandwidthHandler(cachedCtx, tx)
		if err == nil {

			resp := app.BaseApp.CheckTx(txBytes)
			if resp.Code == 0 {
				writeCache()
			}

			return resp
		}
	}

	result := err.Result()

	return abci.ResponseCheckTx{
		Code:      uint32(result.Code),
		Data:      result.Data,
		Log:       result.Log,
		GasWanted: int64(result.GasWanted),
		GasUsed:   int64(result.GasUsed),
		Tags:      result.Tags,
	}

}

func (app *CyberdApp) DeliverTx(txBytes []byte) (res abci.ResponseDeliverTx) {

	cachedCtx, writeCache := app.NewContext(false, abci.Header{Height: app.latestBlockHeight}).CacheContext()

	var tx, err = app.txDecoder(txBytes)
	if err == nil {

		err = app.bandwidthHandler(cachedCtx, tx)
		if err == nil {

			resp := app.BaseApp.DeliverTx(txBytes)
			if resp.Code == 0 {
				writeCache()
			}

			return resp
		}
	}

	result := err.Result()

	return abci.ResponseDeliverTx{
		Code:      uint32(result.Code),
		Codespace: string(result.Codespace),
		Data:      result.Data,
		Log:       result.Log,
		GasWanted: int64(result.GasWanted),
		GasUsed:   int64(result.GasUsed),
		Tags:      result.Tags,
	}
}

// Calculates cyber.Rank for block N, and returns Hash of result as app state.
// Calculated app state will be included in N+1 block header, thus influence on block hash.
// App state is consensus driven state.
func (app *CyberdApp) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {

	linksCount := app.mainStorage.GetLinksCount(ctx)
	cidsCount := app.mainStorage.GetCidsCount(ctx)

	start := time.Now()
	newRank, steps := rank.CalculateRank(app.memStorage, app.computeUnit, app.BaseApp.Logger)
	app.BaseApp.Logger.Info(
		"Rank calculated", "steps", steps, "time", time.Since(start), "links", linksCount, "cids", cidsCount,
	)

	rankAsBytes := make([]byte, 8*len(newRank))
	for i, f64 := range newRank {
		binary.LittleEndian.PutUint64(rankAsBytes[i*8:i*8+8], math.Float64bits(f64))
	}

	hash := sha256.Sum256(rankAsBytes)
	app.latestRankHash = hash[:]
	app.memStorage.UpdateRank(newRank)
	app.mainStorage.StoreAppHash(ctx, hash[:])

	validatorUpdates, tags := stake.EndBlocker(ctx, app.stakeKeeper)
	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Tags:             tags,
	}
}

// Implements ABCI
func (app *CyberdApp) Commit() (res abci.ResponseCommit) {

	app.BaseApp.Commit()
	return abci.ResponseCommit{Data: app.latestRankHash}
}

// Implements ABCI
func (app *CyberdApp) Info(req abci.RequestInfo) abci.ResponseInfo {

	return abci.ResponseInfo{
		Data:             app.BaseApp.Name(),
		LastBlockHeight:  app.LastBlockHeight(),
		LastBlockAppHash: app.appHash(),
	}
}

func (app *CyberdApp) appHash() []byte {

	if app.LastBlockHeight() == 0 {
		return make([]byte, 0)
	}
	return app.latestRankHash
}
