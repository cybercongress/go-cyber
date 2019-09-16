package app

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkbank "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cybercongress/cyberd/store"
	"github.com/cybercongress/cyberd/types"
	cbd "github.com/cybercongress/cyberd/types"
	"github.com/cybercongress/cyberd/types/coin"
	"github.com/cybercongress/cyberd/util"
	"github.com/cybercongress/cyberd/x/bandwidth"
	bw "github.com/cybercongress/cyberd/x/bandwidth/types"
	cbdbank "github.com/cybercongress/cyberd/x/bank"
	"github.com/cybercongress/cyberd/x/debug"
	"github.com/cybercongress/cyberd/x/link"
	"github.com/cybercongress/cyberd/x/link/keeper"
	"github.com/cybercongress/cyberd/x/rank"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"os"
	"sort"
	"time"
)

const (
	APP            = "cyberd"
	appName        = "CyberdApp"
	DefaultKeyPass = "12345678" //todo remove from here
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.cyberdcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.cyberd")

	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		sdkbank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler),
		params.AppModuleBasic{},
        crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
	}

)

// CyberdApp implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore dbKeys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type CyberdApp struct {
	*baseapp.BaseApp
	cdc *codec.Codec

	txDecoder sdk.TxDecoder

	// bandwidth
	bandwidthMeter       bw.BandwidthMeter
	blockBandwidthKeeper bw.BlockSpentBandwidthKeeper

	// keys to access the multistore
	dbKeys CyberdAppDbKeys

	// keepers: manage getting and setting app data
	mainKeeper         store.MainKeeper
	accountKeeper      auth.AccountKeeper
	bankKeeper         *cbdbank.Keeper
	supplyKeeper       supply.Keeper
	stakingKeeper      staking.Keeper
	slashingKeeper     slashing.Keeper
	mintKeeper         mint.Keeper
	distrKeeper        distr.Keeper
	govKeeper          gov.Keeper
	paramsKeeper       params.Keeper
    crisisKeeper       crisis.Keeper
	accBandwidthKeeper bw.Keeper

	// cyberd storage
	// todo: move all processes with this storages to another file
	linkIndexedKeeper *keeper.LinkIndexedKeeper
	cidNumKeeper      keeper.CidNumberKeeper
	stakingIndex      *cbdbank.IndexedKeeper
	rankState         *rank.RankState

	latestBlockHeight int64

	debugState *debug.State
}

// NewCyberdApp returns a reference to a new CyberdApp given a
// logger and database. Internally, a codec is created along with all the necessary dbKeys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewCyberdApp(
	logger log.Logger, db dbm.DB, opts Options, baseAppOptions ...func(*baseapp.BaseApp)) *CyberdApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()
	txDecoder := auth.DefaultTxDecoder(cdc)
	baseApp := baseapp.NewBaseApp(appName, logger, db, txDecoder, baseAppOptions...)
	dbKeys := NewCyberdAppDbKeys()
	mainKeeper := store.NewMainKeeper(dbKeys.main)

	// create your application type
	var app = &CyberdApp{
		cdc:        cdc,
		txDecoder:  txDecoder,
		BaseApp:    baseApp,
		dbKeys:     dbKeys,
		mainKeeper: mainKeeper,
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(app.cdc, dbKeys.params, dbKeys.tParams, params.DefaultCodespace)
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(sdkbank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	mintSubspace := app.paramsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	crisisSubspace := app.paramsKeeper.Subspace(crisis.DefaultParamspace)
	govSubspace := app.paramsKeeper.Subspace(gov.DefaultParamspace)

	blacklistedAddrs := make(map[string]bool)
	// add keepers

	app.accountKeeper = auth.NewAccountKeeper(app.cdc, dbKeys.acc, authSubspace, cbd.NewCyberdAccount)
	app.bankKeeper = cbdbank.NewKeeper(app.accountKeeper, bankSubspace, sdkbank.DefaultCodespace, blacklistedAddrs)
	app.supplyKeeper = supply.NewKeeper(app.cdc, dbKeys.supply, app.accountKeeper, app.bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(
		app.cdc, dbKeys.stake, dbKeys.tStake, app.supplyKeeper, stakingSubspace, staking.DefaultCodespace)
	app.bankKeeper.SetStakingKeeper(&stakingKeeper)
	app.bankKeeper.SetSupplyKeeper(app.supplyKeeper)
	app.mintKeeper = mint.NewKeeper(
		app.cdc, dbKeys.mint, mintSubspace, &stakingKeeper, app.supplyKeeper, auth.FeeCollectorName)
	app.distrKeeper = distr.NewKeeper(
		app.cdc, dbKeys.distr, distrSubspace, &stakingKeeper, app.supplyKeeper,
		distr.DefaultCodespace, auth.FeeCollectorName, blacklistedAddrs)
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc, dbKeys.slashing, &stakingKeeper, slashingSubspace, slashing.DefaultCodespace)
	app.crisisKeeper = crisis.NewKeeper(crisisSubspace, uint(1), app.supplyKeeper, auth.FeeCollectorName)

	app.accBandwidthKeeper = bandwidth.NewAccBandwidthKeeper(dbKeys.accBandwidth)
	app.blockBandwidthKeeper = bandwidth.NewBlockSpentBandwidthKeeper(dbKeys.blockBandwidth)

	govRouter := gov.NewRouter()
    govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
        AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
        AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper))

    app.govKeeper = gov.NewKeeper(
		app.cdc, dbKeys.gov, app.paramsKeeper, govSubspace,
		app.supplyKeeper, &stakingKeeper, gov.DefaultCodespace, govRouter)

	// cyberd keepers
	app.linkIndexedKeeper = keeper.NewLinkIndexedKeeper(keeper.NewBaseLinkKeeper(mainKeeper, dbKeys.links))
	app.cidNumKeeper = keeper.NewBaseCidNumberKeeper(mainKeeper, dbKeys.cidNum, dbKeys.cidNumReverse)
	app.stakingIndex = cbdbank.NewIndexedKeeper(app.bankKeeper, app.accountKeeper)
	app.rankState = rank.NewRankState(
		opts.AllowSearch, app.mainKeeper, app.stakingIndex,
		app.linkIndexedKeeper, app.cidNumKeeper, opts.ComputeUnit,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above are passed by reference,
	// so that it can be modified like below:
	app.stakingKeeper = *stakingKeeper.SetHooks(
		NewStakeHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)
	app.bankKeeper.AddHook(bandwidth.CollectAddressesWithStakeChange())

	app.bandwidthMeter = bandwidth.NewBaseMeter(
		app.mainKeeper, app.accountKeeper, app.bankKeeper, app.accBandwidthKeeper, bandwidth.MsgBandwidthCosts,
		app.blockBandwidthKeeper,
	)

	// register message routes
	app.Router().
		AddRoute("link", link.NewLinksHandler(app.cidNumKeeper, app.linkIndexedKeeper, app.accountKeeper)).
		AddRoute(sdkbank.RouterKey, sdkbank.NewHandler(app.bankKeeper)).
		AddRoute(staking.RouterKey, staking.NewHandler(app.stakingKeeper)).
		AddRoute(distr.RouterKey, distr.NewHandler(app.distrKeeper)).
		AddRoute(slashing.RouterKey, slashing.NewHandler(app.slashingKeeper)).
		AddRoute(gov.RouterKey, gov.NewHandler(app.govKeeper))

	app.QueryRouter().
		AddRoute(auth.QuerierRoute, auth.NewQuerier(app.accountKeeper)).
		AddRoute(distr.QuerierRoute, distr.NewQuerier(app.distrKeeper)).
		AddRoute(gov.QuerierRoute, gov.NewQuerier(app.govKeeper)).
		AddRoute(slashing.QuerierRoute, slashing.NewQuerier(app.slashingKeeper)).
		AddRoute(staking.QuerierRoute, staking.NewQuerier(app.stakingKeeper))

	// mount the multistore and load the latest state
	app.MountStores(
		dbKeys.main, dbKeys.acc, dbKeys.cidNum, dbKeys.cidNumReverse, dbKeys.links, dbKeys.rank, dbKeys.stake,
		dbKeys.slashing, dbKeys.gov, dbKeys.params, dbKeys.distr, dbKeys.fees, dbKeys.accBandwidth,
		dbKeys.blockBandwidth, dbKeys.tParams, dbKeys.tStake, dbKeys.mint, dbKeys.supply,
	)

	app.SetInitChainer(app.applyGenesis)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(NewAnteHandler(app.accountKeeper))

	// perform initialization logic
	err := app.LoadLatestVersion(dbKeys.main)
	if err != nil {
		cmn.Exit(err.Error())
	}

	ctx := app.BaseApp.NewContext(true, abci.Header{})
	app.latestBlockHeight = int64(mainKeeper.GetLatestBlockNumber(ctx))
	ctx = ctx.WithBlockHeight(app.latestBlockHeight)

	// build context for current rank calculation round
	rankRoundBlockNumber := (app.latestBlockHeight / rank.CalculationPeriod) * rank.CalculationPeriod
	if rankRoundBlockNumber == 0 && app.latestBlockHeight >= 1 {
		rankRoundBlockNumber = 1 // special case cause tendermint blocks start from 1
	}
	rankCtx, err := util.NewContextWithMSVersion(db, rankRoundBlockNumber, dbKeys.GetStoreKeys()...)

	// load in-memory data
	start := time.Now()
	app.BaseApp.Logger().Info("Loading mem state")
	app.linkIndexedKeeper.Load(rankCtx, ctx)
	app.stakingIndex.Load(rankCtx, ctx)
	app.BaseApp.Logger().Info("App loaded", "time", time.Since(start))

	// BANDWIDTH LOAD
	app.bandwidthMeter.Load(ctx)

	// RANK PARAMS
	app.rankState.Load(ctx, app.Logger())
	app.debugState = &debug.State{Opts: opts.Debug, StartupBlock: app.latestBlockHeight}
	app.Seal()
	return app
}

func (app *CyberdApp) applyGenesis(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	start := time.Now()
	app.Logger().Info("Applying genesis")
	var genesisState GenesisState
	err := app.cdc.UnmarshalJSON(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}

	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		app.accountKeeper.GetNextAccountNumber(ctx) //increment for future accs being started from right number
		app.accountKeeper.SetAccount(ctx, acc)
		app.stakingIndex.UpdateStake(
		    types.AccNumber(acc.GetAccountNumber()),
		    acc.GetCoins().AmountOf(coin.CYB).Int64())
	}

	cbdbank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)
	// initialize distribution (must happen before staking)
	distr.InitGenesis(ctx, app.distrKeeper, app.supplyKeeper, genesisState.DistrData)
	supply.InitGenesis(ctx, app.supplyKeeper, app.accountKeeper, genesisState.SupplyData)

	validators := staking.InitGenesis(
		ctx, app.stakingKeeper, app.accountKeeper, app.supplyKeeper, genesisState.StakingData)

	auth.InitGenesis(ctx, app.accountKeeper, genesisState.AuthData)
	slashing.InitGenesis(ctx, app.slashingKeeper, app.stakingKeeper, genesisState.SlashingData)
	gov.InitGenesis(ctx, app.govKeeper, app.supplyKeeper, genesisState.GovData)
	mint.InitGenesis(ctx, app.mintKeeper, genesisState.MintData)
	bandwidth.InitGenesis(ctx, app.bandwidthMeter, app.accBandwidthKeeper, genesisState.GetAddresses())

	err = validateGenesisState(genesisState)
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
			bz := abci.RequestDeliverTx{Tx: app.cdc.MustMarshalBinaryLengthPrefixed(tx)}
			res := app.BaseApp.DeliverTx(bz)
			if !res.IsOK() {
				panic(res.Log)
			}
		}

		validators = app.stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
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

	err = link.InitGenesis(ctx, app.cidNumKeeper, app.linkIndexedKeeper, app.Logger())
	if err != nil {
		panic(err)
	}

	app.Logger().Info("Genesis applied", "time", time.Since(start))
	return abci.ResponseInitChain{
		Validators: validators,
	}
}

func (app *CyberdApp) CheckTx(req abci.RequestCheckTx) (res abci.ResponseCheckTx) {

	ctx := app.NewContext(true, abci.Header{Height: app.latestBlockHeight})
	tx, acc, err := app.decodeTxAndAcc(ctx, req.GetTx())

	if err == nil {

		txCost := app.bandwidthMeter.GetPricedTxCost(tx)
		accBw := app.bandwidthMeter.GetCurrentAccBandwidth(ctx, acc)

		if !accBw.HasEnoughRemained(txCost) {
			err = types.ErrNotEnoughBandwidth()
		} else {
			resp := app.BaseApp.CheckTx(req)
			if resp.Code == 0 {
				app.bandwidthMeter.ConsumeAccBandwidth(ctx, accBw, txCost)
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
		Events:    result.Events.ToABCIEvents(),
	}
}

func (app *CyberdApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {

	ctx := app.NewContext(false, abci.Header{Height: app.latestBlockHeight})
	tx, acc, err := app.decodeTxAndAcc(ctx, req.GetTx())

	if err == nil {

		txCost := app.bandwidthMeter.GetPricedTxCost(tx)
		accBw := app.bandwidthMeter.GetCurrentAccBandwidth(ctx, acc)

		if !accBw.HasEnoughRemained(txCost) {
			err = types.ErrNotEnoughBandwidth()
		} else {
			resp := app.BaseApp.DeliverTx(req)
			app.bandwidthMeter.ConsumeAccBandwidth(ctx, accBw, txCost)
			app.bandwidthMeter.AddToBlockBandwidth(app.bandwidthMeter.GetTxCost(tx))

			return abci.ResponseDeliverTx{
				Code:      resp.GetCode(),
				Codespace: resp.GetCodespace(),
				Data:      resp.GetData(),
				Log:       resp.GetLog(),
				GasWanted: resp.GetGasWanted(),
				GasUsed:   resp.GetGasUsed(),
				Events:    append(resp.GetEvents(), getSignersEvents(tx).ToABCIEvents()...),
			}
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
		Events:    result.Events.ToABCIEvents(),
	}
}

func (app *CyberdApp) decodeTxAndAcc(ctx sdk.Context, txBytes []byte) (auth.StdTx, sdk.AccAddress, sdk.Error) {

	decoded, err := app.txDecoder(txBytes)
	if err != nil {
		return auth.StdTx{}, nil, err
	}

	tx := decoded.(auth.StdTx)
	if tx.GetMsgs() == nil || len(tx.GetMsgs()) == 0 {
		return tx, nil, sdk.ErrInternal("Tx.GetMsgs() must return at least one message in list")
	}

	if err := tx.ValidateBasic(); err != nil {
		return tx, nil, err
	}

	// signers acc [0] bandwidth will be consumed
	account := tx.GetSigners()[0]
	acc := app.accountKeeper.GetAccount(ctx, account)
	if acc == nil {
		return tx, nil, sdk.ErrUnknownAddress(account.String())
	}

	return tx, account, nil
}

func getSignersEvents(tx sdk.Tx) sdk.Events {
	events := sdk.EmptyEvents()

	for _, msg := range tx.GetMsgs() {
		for _, accAddress := range msg.GetSigners() {
			attribute := sdk.NewAttribute("address", accAddress.String())
			event := sdk.NewEvent("signer", attribute)
			events.AppendEvent(event)
		}
	}

	return events
}

func (app *CyberdApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {

	debug.BeginBlocker(app.debugState, req, app.Logger())

	// mint new tokens for the previous block
	mint.BeginBlocker(ctx, app.mintKeeper)
	// distribute rewards for the previous block
	distr.BeginBlocker(ctx, req, app.distrKeeper)

	// slash anyone who double signed.
	// NOTE: This should happen after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool,
	// so as to keep the CanWithdrawInvariant invariant.
	slashing.BeginBlocker(ctx, req, app.slashingKeeper)
	return abci.ResponseBeginBlock{Events: ctx.EventManager().ABCIEvents()}
}

// Calculates cyber.Rank for block N, and returns Hash of result as app state.
// Calculated app state will be included in N+1 block header, thus influence on block hash.
// App state is consensus driven state.
func (app *CyberdApp) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {

	app.latestBlockHeight = ctx.BlockHeight()
	app.mainKeeper.StoreLatestBlockNumber(ctx, uint64(ctx.BlockHeight()))

	gov.EndBlocker(ctx, app.govKeeper)
	validatorUpdates := staking.EndBlocker(ctx, app.stakingKeeper)
	app.stakingIndex.EndBlocker(ctx)

	bandwidth.EndBlocker(ctx, app.bandwidthMeter)

	// RANK CALCULATION
	app.rankState.EndBlocker(ctx, app.Logger())

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Events:           ctx.EventManager().ABCIEvents(),
	}
}

// Implements ABCI
func (app *CyberdApp) Commit() (res abci.ResponseCommit) {
	err := app.linkIndexedKeeper.Commit(uint64(app.latestBlockHeight))
	if err != nil {
		panic(err)
	}
	app.BaseApp.Commit()
	return abci.ResponseCommit{Data: app.appHash()}
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

	linkHash := app.linkIndexedKeeper.GetNetworkLinkHash()
	rankHash := app.rankState.GetNetworkRankHash()

	result := make([]byte, len(linkHash))

	for i := 0; i < len(linkHash); i++ {
		result[i] = linkHash[i] ^ rankHash[i]
	}

	return result
}
