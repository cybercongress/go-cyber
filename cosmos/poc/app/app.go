package app

import (
	"encoding/json"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	appName = "CyberdApp"
)

// CyberdApp implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore keys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type CyberdApp struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the multistore
	keyMain    *sdk.KVStoreKey
	keyAccount *sdk.KVStoreKey
	keyIBC     *sdk.KVStoreKey
	keyLink    *sdk.KVStoreKey

	// manage getting and setting accounts
	accountMapper       auth.AccountMapper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	ibcMapper           ibc.Mapper
}

// NewBasecoinApp returns a reference to a new CyberdApp given a logger and
// database. Internally, a codec is created along with all the necessary keys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewCyberdApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *CyberdApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()

	// create your application type
	var app = &CyberdApp{
		cdc:        cdc,
		BaseApp:    bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...),
		keyMain:    sdk.NewKVStoreKey("main"),
		keyAccount: sdk.NewKVStoreKey("acc"),
		keyIBC:     sdk.NewKVStoreKey("ibc"),
		keyLink:    sdk.NewKVStoreKey("link"),
	}

	// define and attach the mappers and keepers
	app.accountMapper = auth.NewAccountMapper(
		cdc,
		app.keyAccount, // target store
		func() auth.Account {
			return &auth.BaseAccount{}
		},
	)
	app.coinKeeper = bank.NewKeeper(app.accountMapper)
	app.ibcMapper = ibc.NewMapper(app.cdc, app.keyIBC, app.RegisterCodespace(ibc.DefaultCodespace))

	// register message routes
	app.Router().
		AddRoute("bank", bank.NewHandler(app.coinKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.coinKeeper)).
		AddRoute("link", NewHandler(app.keyLink))

	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))

	// mount the multistore and load the latest state
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC, app.keyLink)
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	app.Seal()

	return app
}

// MakeCodec creates a new wire codec and registers all the necessary types
// with the codec.
func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()

	wire.RegisterCrypto(cdc)
	sdk.RegisterWire(cdc)
	bank.RegisterWire(cdc)
	ibc.RegisterWire(cdc)
	auth.RegisterWire(cdc)
	RegisterWire(cdc)

	cdc.Seal()

	return cdc
}

// BeginBlocker reflects logic to run before any TXs application are processed
// by the application.
func (app *CyberdApp) BeginBlocker(_ sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

// EndBlocker reflects logic to run after all TXs are processed by the
// application.
func (app *CyberdApp) EndBlocker(_ sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {
	return abci.ResponseEndBlock{}
}

// initChainer implements the custom application logic that the BaseApp will
// invoke upon initialization. In this case, it will take the application's
// state provided by 'req' and attempt to deserialize said state. The state
// should contain all the genesis accounts. These accounts will be added to the
// application's account mapper.
func (app *CyberdApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		// TODO: https://github.com/cosmos/cosmos-sdk/issues/468
		panic(err)
	}

	for _, gacc := range genesisState.Accounts {
		acc, err := gacc.ToBaseAccount()
		if err != nil {
			// TODO: https://github.com/cosmos/cosmos-sdk/issues/468
			panic(err)
		}

		acc.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)
		app.accountMapper.SetAccount(ctx, acc)
	}

	return abci.ResponseInitChain{}
}

// ExportAppStateAndValidators implements custom application logic that exposes
// various parts of the application's state and set of validators. An error is
// returned if any step getting the state or set of validators fails.
func (app *CyberdApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	accounts := []*GenesisAccount{}

	appendAccountsFn := func(acc auth.Account) bool {
		account := &GenesisAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}

		accounts = append(accounts, account)
		return false
	}

	app.accountMapper.IterateAccounts(ctx, appendAccountsFn)

	genState := GenesisState{Accounts: accounts}
	appState, err = wire.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	return appState, validators, err
}
