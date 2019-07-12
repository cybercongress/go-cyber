package app

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/mint"
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
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
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

	// manage getting and setting app data
	mainKeeper          store.MainKeeper
	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          *cbdbank.Keeper
	stakingKeeper       staking.Keeper
	slashingKeeper      slashing.Keeper
	distrKeeper         distr.Keeper
	govKeeper           gov.Keeper
	paramsKeeper        params.Keeper
	accBandwidthKeeper  bw.Keeper
	mintKeeper     		mint.Keeper

	// cyberd storage
	// todo: move all processes with this storages to another file
	linkIndexedKeeper *keeper.LinkIndexedKeeper
	cidNumKeeper      keeper.CidNumberKeeper
	stakingIndex      *cbdbank.IndexedKeeper
	rankState         *rank.RankState

	latestBlockHeight int64

	debugState *debug.State
}

// NewBasecoinApp returns a reference to a new CyberdApp given a
// logger and
// database. Internally, a codec is created along with all the necessary dbKeys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewCyberdApp(
	logger log.Logger, db dbm.DB, opts Options, baseAppOptions ...func(*baseapp.BaseApp),
) *CyberdApp {

	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()
	dbKeys := NewCyberdAppDbKeys()
	ms := store.NewMainKeeper(dbKeys.main)
	var txDecoder = auth.DefaultTxDecoder(cdc)

	// create your application type
	var app = &CyberdApp{
		cdc:        cdc,
		txDecoder:  txDecoder,
		BaseApp:    baseapp.NewBaseApp(appName, logger, db, txDecoder, baseAppOptions...),
		dbKeys:     dbKeys,
		mainKeeper: ms,
	}

	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(app.cdc, dbKeys.fees)
	app.paramsKeeper = params.NewKeeper(app.cdc, dbKeys.params, dbKeys.tParams)
	app.accBandwidthKeeper = bandwidth.NewAccBandwidthKeeper(dbKeys.accBandwidth)
	app.blockBandwidthKeeper = bandwidth.NewBlockSpentBandwidthKeeper(dbKeys.blockBandwidth)

	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc, dbKeys.acc,
		app.paramsKeeper.Subspace(auth.DefaultParamspace), cbd.NewCyberdAccount,
	)

	var stakingKeeper staking.Keeper
	app.bankKeeper = cbdbank.NewBankKeeper(
		app.accountKeeper, &stakingKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
	)
	app.bankKeeper.AddHook(bandwidth.CollectAddressesWithStakeChange())

	stakingKeeper = staking.NewKeeper(
		app.cdc, dbKeys.stake,
		dbKeys.tStake, app.bankKeeper,
		app.paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc, dbKeys.slashing,
		&stakingKeeper, app.paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)
	app.distrKeeper = distr.NewKeeper(
		app.cdc, dbKeys.distr,
		app.paramsKeeper.Subspace(distr.DefaultParamspace),
		app.bankKeeper, &stakingKeeper, app.feeCollectionKeeper,
		distr.DefaultCodespace,
	)

	app.govKeeper = gov.NewKeeper(
		app.cdc, dbKeys.gov,
		app.paramsKeeper, app.paramsKeeper.Subspace(gov.DefaultParamspace), app.bankKeeper, &stakingKeeper,
		gov.DefaultCodespace,
	)

	app.mintKeeper = mint.NewKeeper(
		app.cdc, dbKeys.mint, app.paramsKeeper.Subspace(mint.DefaultParamspace), &stakingKeeper,
		app.feeCollectionKeeper,
	)

	// cyberd keepers
	app.linkIndexedKeeper = keeper.NewLinkIndexedKeeper(keeper.NewBaseLinkKeeper(ms, dbKeys.links))
	app.cidNumKeeper = keeper.NewBaseCidNumberKeeper(ms, dbKeys.cidNum, dbKeys.cidNumReverse)
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

	app.bandwidthMeter = bandwidth.NewBaseMeter(
		app.mainKeeper, app.accountKeeper, app.bankKeeper, app.accBandwidthKeeper, bandwidth.MsgBandwidthCosts,
		app.blockBandwidthKeeper,
	)

	// register message routes
	app.Router().
		AddRoute("link", link.NewLinksHandler(app.cidNumKeeper, app.linkIndexedKeeper, app.accountKeeper)).
		AddRoute(bank.RouterKey, bank.NewHandler(app.bankKeeper)).
		AddRoute(staking.RouterKey, staking.NewHandler(app.stakingKeeper)).
		AddRoute(distr.RouterKey, distr.NewHandler(app.distrKeeper)).
		AddRoute(slashing.RouterKey, slashing.NewHandler(app.slashingKeeper)).
		AddRoute(gov.RouterKey, gov.NewHandler(app.govKeeper))

	app.QueryRouter().
		AddRoute(auth.QuerierRoute, auth.NewQuerier(app.accountKeeper)).
		AddRoute(distr.QuerierRoute, distr.NewQuerier(app.distrKeeper)).
		AddRoute(gov.QuerierRoute, gov.NewQuerier(app.govKeeper)).
		AddRoute(slashing.QuerierRoute, slashing.NewQuerier(app.slashingKeeper, app.cdc)).
		AddRoute(staking.QuerierRoute, staking.NewQuerier(app.stakingKeeper, app.cdc))

	// mount the multistore and load the latest state
	app.MountStores(
		dbKeys.main, dbKeys.acc, dbKeys.cidNum, dbKeys.cidNumReverse, dbKeys.links, dbKeys.rank, dbKeys.stake,
		dbKeys.slashing, dbKeys.gov, dbKeys.params, dbKeys.distr, dbKeys.fees, dbKeys.accBandwidth,
		dbKeys.blockBandwidth, dbKeys.tDistr, dbKeys.tParams, dbKeys.tStake, dbKeys.mint,
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
	app.latestBlockHeight = int64(ms.GetLatestBlockNumber(ctx))
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
	stateJSON := req.AppStateBytes
	var genesisState GenesisState
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err)
	}

	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		app.accountKeeper.GetNextAccountNumber(ctx) //increment for future accs being started from right number
		app.accountKeeper.SetAccount(ctx, acc)
		app.stakingIndex.UpdateStake(types.AccNumber(acc.AccountNumber), acc.Coins.AmountOf(coin.CYB).Int64())
	}

	bank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)
	// initialize distribution (must happen before staking)
	distr.InitGenesis(ctx, app.distrKeeper, genesisState.DistrData)

	validators, err := staking.InitGenesis(ctx, app.stakingKeeper, genesisState.StakingData)
	if err != nil {
		panic(err)
	}

	auth.InitGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper, genesisState.AuthData)

	slashing.InitGenesis(
		ctx, app.slashingKeeper, genesisState.SlashingData,
		genesisState.StakingData.Validators.ToSDKValidators(),
	)
	gov.InitGenesis(ctx, app.govKeeper, genesisState.GovData)
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
			bz := app.cdc.MustMarshalBinaryLengthPrefixed(tx)
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

func (app *CyberdApp) CheckTx(txBytes []byte) (res abci.ResponseCheckTx) {

	ctx := app.NewContext(true, abci.Header{Height: app.latestBlockHeight})
	tx, acc, err := app.decodeTxAndAcc(ctx, txBytes)

	if err == nil {

		txCost := app.bandwidthMeter.GetPricedTxCost(tx)
		accBw := app.bandwidthMeter.GetCurrentAccBandwidth(ctx, acc)

		if !accBw.HasEnoughRemained(txCost) {
			err = types.ErrNotEnoughBandwidth()
		} else {

			resp := app.BaseApp.CheckTx(txBytes)
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
		Tags:      result.Tags,
	}
}

func (app *CyberdApp) DeliverTx(txBytes []byte) (res abci.ResponseDeliverTx) {

	ctx := app.NewContext(false, abci.Header{Height: app.latestBlockHeight})
	tx, acc, err := app.decodeTxAndAcc(ctx, txBytes)

	if err == nil {

		txCost := app.bandwidthMeter.GetPricedTxCost(tx)
		accBw := app.bandwidthMeter.GetCurrentAccBandwidth(ctx, acc)

		if !accBw.HasEnoughRemained(txCost) {
			err = types.ErrNotEnoughBandwidth()
		} else {

			resp := app.BaseApp.DeliverTx(txBytes)
			app.bandwidthMeter.ConsumeAccBandwidth(ctx, accBw, txCost)
			app.bandwidthMeter.AddToBlockBandwidth(app.bandwidthMeter.GetTxCost(tx))

			return abci.ResponseDeliverTx{
				Code:      uint32(resp.Code),
				Codespace: string(resp.Codespace),
				Data:      resp.Data,
				Log:       resp.Log,
				GasWanted: int64(resp.GasWanted),
				GasUsed:   int64(resp.GasUsed),
				Tags:      append(resp.Tags, getSignersTags(tx)...),
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
		Tags:      result.Tags,
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

func getSignersTags(tx sdk.Tx) sdk.Tags {

	signers := make(map[string]struct{})
	for _, msg := range tx.GetMsgs() {
		for _, s := range msg.GetSigners() {
			signers[s.String()] = struct{}{}
		}
	}

	var tags = make(sdk.Tags, 0)
	for signer := range signers {
		tags.AppendTag("signer", signer)
	}
	return tags
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
	tags := slashing.BeginBlocker(ctx, req, app.slashingKeeper)
	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
}

// Calculates cyber.Rank for block N, and returns Hash of result as app state.
// Calculated app state will be included in N+1 block header, thus influence on block hash.
// App state is consensus driven state.
func (app *CyberdApp) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {

	app.latestBlockHeight = ctx.BlockHeight()
	app.mainKeeper.StoreLatestBlockNumber(ctx, uint64(ctx.BlockHeight()))

	tags := gov.EndBlocker(ctx, app.govKeeper)
	validatorUpdates, stakingTags := staking.EndBlocker(ctx, app.stakingKeeper)
	app.stakingIndex.EndBlocker(ctx)

	bandwidth.EndBlocker(ctx, app.bandwidthMeter)

	// RANK CALCULATION
	app.rankState.EndBlocker(ctx, app.Logger())

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Tags:             append(tags, stakingTags...),
	}
}

// Implements ABCI
func (app *CyberdApp) Commit() (res abci.ResponseCommit) {
	err := app.linkIndexedKeeper.Commit(uint64(app.latestBlockHeight))
	if err != nil {
		panic(err)
	}
	app.BaseApp.Commit()
	return abci.ResponseCommit{Data: app.rankState.GetNetworkRankHash()}
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
	return app.rankState.GetNetworkRankHash()
}
