package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	APP     = "cyberd"
	appName = "CyberdApp"
)

type CyberdAppDbKeys struct {
	main     *sdk.KVStoreKey
	acc      *sdk.KVStoreKey
	cidIndex *sdk.KVStoreKey
	cidIns   *sdk.KVStoreKey
	cidOuts  *sdk.KVStoreKey
	rank     *sdk.KVStoreKey
}

// CyberdApp implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore dbKeys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type CyberdApp struct {
	*baseapp.BaseApp
	cdc *wire.Codec

	// keys to access the multistore
	dbKeys CyberdAppDbKeys

	// manage getting and setting accounts
	accountMapper       auth.AccountMapper
	cidIndexMapper      CidIndexStorage
	inCidsMapper        LinksStorage
	outCidsMapper       LinksStorage
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	ibcMapper           ibc.Mapper
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
		cidIns:   sdk.NewKVStoreKey("cid_ins"),
		cidOuts:  sdk.NewKVStoreKey("cid_outs"),
		rank:     sdk.NewKVStoreKey("rank"),
	}

	// create your application type
	var app = &CyberdApp{
		cdc:     cdc,
		BaseApp: baseapp.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...),
		dbKeys:  dbKeys,
	}

	// define and attach the mappers and keepers
	app.accountMapper = auth.NewAccountMapper(cdc, dbKeys.acc, NewAccount)
	app.coinKeeper = bank.NewKeeper(app.accountMapper)
	app.cidIndexMapper = CidIndexStorage{mainStoreKey: dbKeys.main, indexKey: dbKeys.cidIndex, cdc: cdc}
	app.inCidsMapper = LinksStorage{key: dbKeys.cidIns, cdc: cdc}
	app.outCidsMapper = LinksStorage{key: dbKeys.cidOuts, cdc: cdc}

	// register message routes
	app.Router().
		AddRoute("bank", bank.NewHandler(app.coinKeeper)).
		AddRoute("link", NewLinksHandler(app.cidIndexMapper, app.inCidsMapper, app.outCidsMapper))

	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))

	// mount the multistore and load the latest state
	app.MountStoresIAVL(dbKeys.main, dbKeys.acc, dbKeys.cidIndex, dbKeys.cidIns, dbKeys.cidOuts, dbKeys.rank)
	err := app.LoadLatestVersion(dbKeys.main)
	if err != nil {
		cmn.Exit(err.Error())
	}

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
	return abci.ResponseEndBlock{}
}
