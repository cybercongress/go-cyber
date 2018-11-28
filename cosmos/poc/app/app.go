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
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	. "github.com/cybercongress/cyberd/cosmos/poc/app/bank"
	"github.com/cybercongress/cyberd/cosmos/poc/app/coin"
	"github.com/cybercongress/cyberd/cosmos/poc/app/rank"
	. "github.com/cybercongress/cyberd/cosmos/poc/app/storage"
	cbd "github.com/cybercongress/cyberd/cosmos/poc/app/types"
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
	APP     = "cyberd"
	appName = "CyberdApp"
	DefaultKeyPass = "12345678"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.cyberdcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.cyberd")
)

type CyberdAppDbKeys struct {
	main     *sdk.KVStoreKey
	acc      *sdk.KVStoreKey
	accIndex *sdk.KVStoreKey
	cidIndex *sdk.KVStoreKey
	links    *sdk.KVStoreKey
	rank     *sdk.KVStoreKey
	stake    *sdk.KVStoreKey
	tStake   *sdk.TransientStoreKey
	slashing *sdk.KVStoreKey
	params   *sdk.KVStoreKey
	tParams  *sdk.TransientStoreKey
}

// CyberdApp implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore dbKeys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type CyberdApp struct {
	*baseapp.BaseApp
	cdc *codec.Codec

	// keys to access the multistore
	dbKeys CyberdAppDbKeys

	// manage getting and setting accounts
	mainStorage         MainStorage
	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	stakeKeeper         stake.Keeper
	slashingKeeper      slashing.Keeper
	paramsKeeper        params.Keeper

	// cyberd storages
	persistStorages CyberdPersistentStorages
	memStorage      *InMemoryStorage

	latestRankHash []byte

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
		tStake:   sdk.NewTransientStoreKey("transient_stake"),
		slashing: sdk.NewKVStoreKey("slashing"),
		params:   sdk.NewKVStoreKey("params"),
		tParams:  sdk.NewTransientStoreKey("transient_params"),
	}

	ms := NewMainStorage(dbKeys.main)
	storages := CyberdPersistentStorages{
		CidIndex: NewCidIndexStorage(ms, dbKeys.cidIndex),
		Links:    NewLinksStorage(ms, dbKeys.links),
		Rank:     NewRankStorage(ms, dbKeys.rank),
	}

	// create your application type
	var app = &CyberdApp{
		cdc:             cdc,
		BaseApp:         baseapp.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...),
		dbKeys:          dbKeys,
		persistStorages: storages,
		mainStorage:     ms,
		computeUnit:     computeUnit,
	}

	// define and attach the mappers and keepers
	app.accountKeeper = auth.NewAccountKeeper(app.cdc, dbKeys.acc, NewAccount)
	app.coinKeeper = bank.NewBaseKeeper(app.accountKeeper)
	app.paramsKeeper = params.NewKeeper(app.cdc, dbKeys.params, dbKeys.tParams)
	stakeKeeper := stake.NewKeeper(
		app.cdc, dbKeys.stake,
		dbKeys.tStake, app.coinKeeper,
		app.paramsKeeper.Subspace(stake.DefaultParamspace),
		app.RegisterCodespace(stake.DefaultCodespace),
	)
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		dbKeys.slashing,
		&stakeKeeper, app.paramsKeeper.Subspace(slashing.DefaultParamspace),
		app.RegisterCodespace(slashing.DefaultCodespace),
	)
	app.memStorage = &InMemoryStorage{}

	// register the staking hooks
	// NOTE: stakeKeeper above are passed by reference,
	// so that it can be modified like below:
	app.stakeKeeper = *stakeKeeper.SetHooks(NewHooks(app.slashingKeeper.Hooks()))

	// register message routes
	app.Router().
		AddRoute("bank", NewBankHandler(app.coinKeeper, app.memStorage, app.accountKeeper)).
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
	app.MountStoresIAVL(dbKeys.main, dbKeys.acc, dbKeys.cidIndex, dbKeys.links, dbKeys.rank, dbKeys.stake,
		dbKeys.slashing, dbKeys.params)
	app.MountStoresTransient(dbKeys.tParams, dbKeys.tStake)
	err := app.LoadLatestVersion(dbKeys.main)
	if err != nil {
		cmn.Exit(err.Error())
	}
	ctx := app.BaseApp.NewContext(true, abci.Header{})
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
	// TODO is this now the whole genesis file?

	var genesisState GenesisState
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
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
		panic(err) // TODO find a way to do this w/o panics
	}

	// initialize module-specific stores
	slashing.InitGenesis(ctx, app.slashingKeeper, genesisState.SlashingData, genesisState.StakeData)
	err = CyberdValidateGenesisState(genesisState)
	if err != nil {
		panic(err) // TODO find a way to do this w/o panics
	}

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
	tags := slashing.BeginBlocker(ctx, req, app.slashingKeeper)

	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
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

	validatorUpdates := stake.EndBlocker(ctx, app.stakeKeeper)
	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
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
