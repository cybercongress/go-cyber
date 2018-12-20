package app

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkbank "github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	. "github.com/cybercongress/cyberd/app/genesis"
	cbd "github.com/cybercongress/cyberd/app/types"
	"github.com/cybercongress/cyberd/app/types/coin"
	"github.com/cybercongress/cyberd/store"
	"github.com/cybercongress/cyberd/types"
	"github.com/cybercongress/cyberd/x/bandwidth"
	bw "github.com/cybercongress/cyberd/x/bandwidth/types"
	cbdbank "github.com/cybercongress/cyberd/x/bank"
	"github.com/cybercongress/cyberd/x/link"
	"github.com/cybercongress/cyberd/x/link/keeper"
	"github.com/cybercongress/cyberd/x/mint"
	"github.com/cybercongress/cyberd/x/rank"
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
	DefaultKeyPass = "12345678" //todo remove from here
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.cyberdcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.cyberd")
)

type CyberdAppDbKeys struct {
	main          *sdk.KVStoreKey
	acc           *sdk.KVStoreKey
	accIndex      *sdk.KVStoreKey
	cidNum        *sdk.KVStoreKey
	cidNumReverse *sdk.KVStoreKey
	links         *sdk.KVStoreKey
	rank          *sdk.KVStoreKey
	stake         *sdk.KVStoreKey
	tStake        *sdk.TransientStoreKey
	fees          *sdk.KVStoreKey
	distr         *sdk.KVStoreKey
	slashing      *sdk.KVStoreKey
	params        *sdk.KVStoreKey
	tParams       *sdk.TransientStoreKey
	accBandwidth  *sdk.KVStoreKey
}

// CyberdApp implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore dbKeys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type CyberdApp struct {
	*baseapp.BaseApp
	cdc *codec.Codec

	txDecoder sdk.TxDecoder

	// bandwidth
	bandwidthMeter bw.BandwidthMeter
	//todo 3 fields below should be in bw meter
	curBlockSpentBandwidth  uint64 //resets every block
	lastTotalSpentBandwidth uint64 //resets every bandwidth price adjustment interval
	currentCreditPrice      float64

	// keys to access the multistore
	dbKeys CyberdAppDbKeys

	// manage getting and setting app data
	mainKeeper          store.MainKeeper
	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          cbdbank.Keeper
	stakeKeeper         stake.Keeper
	slashingKeeper      slashing.Keeper
	distrKeeper         distr.Keeper
	paramsKeeper        params.Keeper
	accBandwidthKeeper  bw.Keeper

	//inflation
	minter mint.Minter

	// cyberd storages
	linkIndexedKeeper keeper.LinkIndexedKeeper
	cidNumKeeper      keeper.CidNumberKeeper
	stakeIndex        *cbdbank.IndexedKeeper
	rankState         *rank.RankState

	latestRankHash    []byte
	latestBlockHeight int64

	computeUnit rank.ComputeUnit

	rankCalculationFinished bool
	lastCidCount            int64
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
		main:          sdk.NewKVStoreKey("main"),
		acc:           sdk.NewKVStoreKey("acc"),
		cidNum:        sdk.NewKVStoreKey("cid_index"),
		cidNumReverse: sdk.NewKVStoreKey("cid_index_reverse"),
		links:         sdk.NewKVStoreKey("links"),
		rank:          sdk.NewKVStoreKey("rank"),
		stake:         sdk.NewKVStoreKey("stake"),
		fees:          sdk.NewKVStoreKey("fee"),
		tStake:        sdk.NewTransientStoreKey("transient_stake"),
		distr:         sdk.NewKVStoreKey("distr"),
		slashing:      sdk.NewKVStoreKey("slashing"),
		params:        sdk.NewKVStoreKey("params"),
		tParams:       sdk.NewTransientStoreKey("transient_params"),
		accBandwidth:  sdk.NewKVStoreKey("acc_bandwidth"),
	}

	ms := store.NewMainKeeper(dbKeys.main)

	var txDecoder = auth.DefaultTxDecoder(cdc)

	// create your application type
	var app = &CyberdApp{
		cdc:         cdc,
		txDecoder:   txDecoder,
		BaseApp:     baseapp.NewBaseApp(appName, logger, db, txDecoder, baseAppOptions...),
		dbKeys:      dbKeys,
		mainKeeper:  ms,
		computeUnit: computeUnit,
	}

	// define and attach the mappers and keepers
	app.accountKeeper = auth.NewAccountKeeper(app.cdc, dbKeys.acc, cbd.NewCyberdAccount)
	app.paramsKeeper = params.NewKeeper(app.cdc, dbKeys.params, dbKeys.tParams)
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(app.cdc, dbKeys.fees)

	var stakeKeeper stake.Keeper

	app.bankKeeper = cbdbank.NewBankKeeper(app.accountKeeper, &stakeKeeper)
	app.bankKeeper.AddHook(bandwidth.CollectAddressesWithStakeChange())

	app.accBandwidthKeeper = bandwidth.NewAccBandwidthKeeper(dbKeys.accBandwidth)

	stakeKeeper = stake.NewKeeper(
		app.cdc, dbKeys.stake,
		dbKeys.tStake, app.bankKeeper,
		app.paramsKeeper.Subspace(stake.DefaultParamspace),
		stake.DefaultCodespace,
	)
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc, dbKeys.slashing,
		stakeKeeper, app.paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)
	app.distrKeeper = distr.NewKeeper(
		app.cdc, dbKeys.distr,
		app.paramsKeeper.Subspace(distr.DefaultParamspace),
		app.bankKeeper, stakeKeeper, app.feeCollectionKeeper,
		distr.DefaultCodespace,
	)
	app.minter = mint.NewMinter(app.feeCollectionKeeper, stakeKeeper)

	// cyberd keepers
	app.linkIndexedKeeper = keeper.NewLinkIndexedKeeper(keeper.NewBaseLinkKeeper(ms, dbKeys.links))
	app.cidNumKeeper = keeper.NewBaseCidNumberKeeper(ms, dbKeys.cidNum, dbKeys.cidNumReverse)
	app.stakeIndex = cbdbank.NewIndexedKeeper(&app.bankKeeper, app.accountKeeper)
	app.rankState = rank.NewRankState(&app.linkIndexedKeeper)

	// register the staking hooks
	// NOTE: stakeKeeper above are passed by reference,
	// so that it can be modified like below:
	app.stakeKeeper = *stakeKeeper.SetHooks(
		NewStakeHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	app.bandwidthMeter = bandwidth.NewBaseMeter(
		app.accountKeeper, app.bankKeeper, app.accBandwidthKeeper, bandwidth.MsgBandwidthCosts,
	)

	// register message routes
	app.Router().
		AddRoute("bank", sdkbank.NewHandler(app.bankKeeper)).
		AddRoute("link", link.NewLinksHandler(app.cidNumKeeper, &app.linkIndexedKeeper, app.accountKeeper)).
		AddRoute("stake", stake.NewHandler(app.stakeKeeper)).
		AddRoute("distr", distr.NewHandler(app.distrKeeper)).
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
		dbKeys.main, dbKeys.acc, dbKeys.cidNum, dbKeys.cidNumReverse, dbKeys.links, dbKeys.rank, dbKeys.stake,
		dbKeys.slashing, dbKeys.params, dbKeys.distr, dbKeys.fees, dbKeys.accBandwidth,
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
	app.linkIndexedKeeper.Load(ctx)
	app.stakeIndex.Load(ctx)
	app.BaseApp.Logger.Info("App loaded", "time", time.Since(start))
	app.latestRankHash = ms.GetAppHash(ctx)
	app.lastTotalSpentBandwidth = ms.GetSpentBandwidth(ctx)
	app.currentCreditPrice = math.Float64frombits(ms.GetBandwidthPrice(ctx, bandwidth.BaseCreditPrice))
	app.curBlockSpentBandwidth = 0
	app.latestBlockHeight = int64(ms.GetLatestBlockNumber(ctx))
	app.rankCalculationFinished = true

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
		app.stakeIndex.UpdateStake(types.AccNumber(acc.AccountNumber), acc.Coins.AmountOf(coin.CBD).Int64())
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

	bandwidth.InitGenesis(ctx, app.bandwidthMeter, app.accBandwidthKeeper, genesisState.Accounts)
	return abci.ResponseInitChain{
		Validators: validators,
	}
}

func (app *CyberdApp) CheckTx(txBytes []byte) (res abci.ResponseCheckTx) {

	ctx := app.NewContext(true, abci.Header{Height: app.latestBlockHeight})
	tx, acc, err := app.decodeTxAndAcc(ctx, txBytes)

	if err == nil {

		txCost := app.bandwidthMeter.GetTxCost(ctx, app.currentCreditPrice, tx)
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

		txCost := app.bandwidthMeter.GetTxCost(ctx, app.currentCreditPrice, tx)
		accBw := app.bandwidthMeter.GetCurrentAccBandwidth(ctx, acc)

		if !accBw.HasEnoughRemained(txCost) {
			err = types.ErrNotEnoughBandwidth()
		} else {

			resp := app.BaseApp.DeliverTx(txBytes)
			app.bandwidthMeter.ConsumeAccBandwidth(ctx, accBw, txCost)
			app.curBlockSpentBandwidth = app.curBlockSpentBandwidth + uint64(txCost)

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
		tags.AppendTag("signer", []byte(signer))
	}
	return tags
}

func (app *CyberdApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {

	// mint new tokens for the previous block
	mint.BeginBlocker(ctx, app.minter)
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

var rankChan = make(chan []float64, 1)

// Calculates cyber.Rank for block N, and returns Hash of result as app state.
// Calculated app state will be included in N+1 block header, thus influence on block hash.
// App state is consensus driven state.
func (app *CyberdApp) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {

	app.latestBlockHeight = ctx.BlockHeight()
	app.mainKeeper.StoreLatestBlockNumber(ctx, uint64(ctx.BlockHeight()))

	validatorUpdates, tags := stake.EndBlocker(ctx, app.stakeKeeper)
	app.stakeIndex.EndBlocker(ctx)

	newPrice, totalSpentBandwidth := bandwidth.EndBlocker(
		ctx, app.lastTotalSpentBandwidth+app.curBlockSpentBandwidth, app.currentCreditPrice,
		app.mainKeeper, app.bandwidthMeter,
	)
	app.lastTotalSpentBandwidth = totalSpentBandwidth
	app.currentCreditPrice = newPrice
	app.curBlockSpentBandwidth = 0

	// START RANK CALCULATION
	//linksCount := app.mainKeeper.GetLinksCount(ctx)

	//start := time.Now()
	// TODO: Deal with hashes (new, old)
	// TODO: Run reindexation in parallel
	if ctx.BlockHeight()%rank.CalculationPeriod == 0 || ctx.BlockHeight() == 1 {

		if !app.rankCalculationFinished {
			newRank := <-rankChan
			app.rankState.SetNewRank(newRank)
			app.rankCalculationFinished = true
		}

		app.rankState.ApplyNewRank(app.lastCidCount)
		app.linkIndexedKeeper.MergeNewLinks()

		rankHash := app.rankState.CalculateHash()
		app.latestRankHash = rankHash
		app.mainKeeper.StoreAppHash(ctx, rankHash)

		app.rankCalculationFinished = false

		app.lastCidCount = int64(app.mainKeeper.GetCidsCount(ctx))
		calcCtx := rank.NewCalcContext(ctx, app.linkIndexedKeeper, app.cidNumKeeper, app.stakeIndex)
		go rank.CalculateRank(calcCtx, rankChan, app.computeUnit, app.BaseApp.Logger)

	} else if !app.rankCalculationFinished {

		select {
		case newRank := <-rankChan:
			app.rankState.SetNewRank(newRank)
			app.rankCalculationFinished = true
		default:
		}

	}
	//app.BaseApp.Logger.Info(
	//	"Rank calculated", "steps", steps, "time", time.Since(start), "links", linksCount, "cids", cidsCount,
	//)

	// END RANK CALCULATION

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
