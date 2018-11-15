package app

import (
	"crypto/sha256"
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	. "github.com/cybercongress/cyberd/cosmos/poc/app/bank"
	"github.com/cybercongress/cyberd/cosmos/poc/app/rank"
	. "github.com/cybercongress/cyberd/cosmos/poc/app/storage"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"math"
	"os"
	"time"
)

const (
	APP     = "cyberd"
	appName = "CyberdApp"
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
	accStorage          auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper

	// cyberd storages
	persistStorages CyberdPersistentStorages
	memStorage      *InMemoryStorage

	latestRankHash []byte
}

// NewBasecoinApp returns a reference to a new CyberdApp given a
// logger and
// database. Internally, a codec is created along with all the necessary dbKeys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewCyberdApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*baseapp.BaseApp)) *CyberdApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()

	dbKeys := CyberdAppDbKeys{
		main:     sdk.NewKVStoreKey("main"),
		acc:      sdk.NewKVStoreKey("acc"),
		cidIndex: sdk.NewKVStoreKey("cid_index"),
		links:    sdk.NewKVStoreKey("links"),
		rank:     sdk.NewKVStoreKey("rank"),
	}

	ms := NewMainStorage(dbKeys.main)
	storages := CyberdPersistentStorages{
		CidIndex: NewCidIndexStorage(ms, dbKeys.cidIndex),
		Links:    NewLinksStorage(dbKeys.links, cdc),
		Rank:     NewRankStorage(ms, dbKeys.rank),
	}

	// create your application type
	var app = &CyberdApp{
		cdc:             cdc,
		BaseApp:         baseapp.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...),
		dbKeys:          dbKeys,
		persistStorages: storages,
		mainStorage:     ms,
	}

	// define and attach the mappers and keepers
	app.accStorage = auth.NewAccountKeeper(app.cdc, dbKeys.acc, NewAccount)
	app.coinKeeper = bank.NewBaseKeeper(app.accStorage)
	app.memStorage = &InMemoryStorage{}

	// register message routes
	app.Router().
		AddRoute("bank", NewBankHandler(app.coinKeeper, app.memStorage)).
		AddRoute("link", NewLinksHandler(storages.CidIndex, storages.Links, app.memStorage))

	// perform initialization logic
	app.SetInitChainer(NewGenesisApplier(app.memStorage, app.cdc, app.accStorage))
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accStorage, app.feeCollectionKeeper))

	// mount the multistore and load the latest state
	app.MountStoresIAVL(dbKeys.main, dbKeys.acc, dbKeys.cidIndex, dbKeys.links, dbKeys.rank)
	err := app.LoadLatestVersion(dbKeys.main)
	if err != nil {
		cmn.Exit(err.Error())
	}
	ctx := app.BaseApp.NewContext(true, abci.Header{})
	start := time.Now()
	app.BaseApp.Logger.Info("Loading mem state")
	app.memStorage.Load(ctx, storages, app.accStorage)
	app.BaseApp.Logger.Info("App loaded", "time", time.Since(start))
	app.latestRankHash = ms.GetAppHash(ctx)
	app.Seal()
	return app
}

// BeginBlocker reflects logic to run before any TXs application are processed
// by the application.
func (app *CyberdApp) BeginBlocker(_ sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

// Calculates cyber.Rank for block N, and returns Hash of result as app state.
// Calculated app state will be included in N+1 block header, thus influence on block hash.
// App state is consensus driven state.
func (app *CyberdApp) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {

	start := time.Now()
	app.BaseApp.Logger.Info("Calculating rank")
	newRank, steps := rank.CalculateRank(app.memStorage)
	app.BaseApp.Logger.Info("Rank calculated", "steps", steps, "time", time.Since(start))

	rankAsBytes := make([]byte, 8*len(newRank))
	for i, f64 := range newRank {
		binary.LittleEndian.PutUint64(rankAsBytes[i*8:i*8+8], math.Float64bits(f64))
	}

	hash := sha256.Sum256(rankAsBytes)
	app.latestRankHash = hash[:]
	app.memStorage.UpdateRank(newRank)
	app.mainStorage.StoreAppHash(ctx, hash[:])
	return abci.ResponseEndBlock{}
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
