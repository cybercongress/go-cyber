package app

import (
	"fmt"
	"io"
	"os"
	"time"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"

	v2 "github.com/cybercongress/go-cyber/v4/app/upgrades/v2"
	v3 "github.com/cybercongress/go-cyber/v4/app/upgrades/v3"
	v4 "github.com/cybercongress/go-cyber/v4/app/upgrades/v4"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"

	"github.com/cybercongress/go-cyber/v4/app/keepers"
	"github.com/cybercongress/go-cyber/v4/app/upgrades"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/spf13/cast"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/cosmos/cosmos-sdk/baseapp"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"

	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	tmjson "github.com/cometbft/cometbft/libs/json"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/cybercongress/go-cyber/v4/utils"

	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
)

const (
	appName = "BostromHub"
)

// We pull these out so we can set them with LDFLAGS in the Makefile
var (
	NodeDir      = ".cyber"
	Bech32Prefix = "bostrom"

	Upgrades = []upgrades.Upgrade{v2.Upgrade, v3.Upgrade, v4.Upgrade}
)

// These constants are derived from the above variables.
// These are the ones we will want to use in the code, based on
// any overrides above
var (
	// DefaultNodeHome default home directories for wasmd
	DefaultNodeHome = os.ExpandEnv("$HOME/") + NodeDir
)

func getGovProposalHandlers() []govclient.ProposalHandler {
	return []govclient.ProposalHandler{
		paramsclient.ProposalHandler,
		upgradeclient.LegacyProposalHandler,
		upgradeclient.LegacyCancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
	}
}

//// module accounts that are allowed to receive tokens
//var allowedReceivingModAcc = map[string]bool{
//	distrtypes.ModuleName: true,
//}

var (
	_ runtime.AppI            = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp
	keepers.AppKeepers

	aminoCodec        *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	ModuleManager *module.Manager

	sm *module.SimulationManager

	configurator module.Configurator
}

// New returns a reference to an initialized Cyber Consensus Computer.
func NewApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	wasmtypes.MaxWasmSize = 2 * 1024 * 1024 // 2MB
	encodingConfig := MakeEncodingConfig()

	appCodec, legacyAmino := encodingConfig.Marshaler, encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	bApp := baseapp.NewBaseApp(appName, logger, db, txConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	app := &App{
		BaseApp:           bApp,
		aminoCodec:        legacyAmino,
		appCodec:          appCodec,
		txConfig:          txConfig,
		interfaceRegistry: interfaceRegistry,
	}

	// Setup keepers
	app.AppKeepers = keepers.NewAppKeepers(
		appCodec,
		bApp,
		legacyAmino,
		keepers.GetMaccPerms(),
		appOpts,
		wasmOpts,
	)

	// load state streaming if enabled
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, logger, app.GetKVStoreKey()); err != nil {
		logger.Error("failed to load state streaming", "err", err)
		os.Exit(1)
	}

	// upgrade handlers
	app.configurator = module.NewConfigurator(appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.ModuleManager = module.NewManager(appModules(app, encodingConfig, skipGenesisInvariants)...)
	app.ModuleManager.RegisterServices(app.configurator)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.ModuleManager.SetOrderBeginBlockers(orderBeginBlockers()...)

	app.ModuleManager.SetOrderEndBlockers(orderEndBlockers()...)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.ModuleManager.SetOrderInitGenesis(orderInitBlockers()...)

	app.ModuleManager.RegisterInvariants(app.CrisisKeeper)

	// NOTE hack to register bank's module services because of custom wrapper
	// from sdk.bank/module.go: AppModule RegisterServices
	// 	m := keeper.NewMigrator(am.keeper.(keeper.BaseKeeper))
	// NOTE skip bank module from native services registration then initialize manually
	//delete(app.ModuleManager.Modules, banktypes.ModuleName)
	//app.ModuleManager.RegisterServices(cfg)
	//app.ModuleManager.Modules[banktypes.ModuleName] = bank.NewAppModule(encodingConfig.Marshaler, app.CyberbankKeeper.Proxy, app.AccountKeeper)
	//
	//banktypes.RegisterMsgServer(cfg.MsgServer(), bankkeeper.NewMsgServerImpl(app.CyberbankKeeper.Proxy))
	//banktypes.RegisterQueryServer(cfg.QueryServer(), app.CyberbankKeeper.Proxy)

	// initialize stores
	app.MountKVStores(app.GetKVStoreKey())
	app.MountTransientStores(app.GetTransientStoreKey())
	app.MountMemoryStores(app.GetMemoryStoreKey())

	// register upgrade
	app.setupUpgradeHandlers(app.configurator)

	// SDK v47 - since we do not use dep inject, this gives us access to newer gRPC services.
	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.ModuleManager.Modules))
	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.CyberbankKeeper.Proxy,
				FeegrantKeeper:  app.FeeGrantKeeper,
				SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			BandwidthMeter:    app.BandwidthMeter,
			IBCKeeper:         app.IBCKeeper,
			TXCounterStoreKey: app.GetKey(wasmtypes.StoreKey),
			WasmConfig:        &wasmConfig,
			WasmKeeper:        &app.WasmKeeper,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	// initialize BaseApp
	app.SetAnteHandler(anteHandler)
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	if manager := app.SnapshotManager(); manager != nil {
		err = manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.WasmKeeper),
		)
		if err != nil {
			panic("failed to register snapshot extension: " + err.Error())
		}
	}

	app.setupUpgradeStoreLoaders()

	app.setPostHandler()

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(fmt.Sprintf("failed to load latest version: %s", err))
		}
		ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})

		// TODO refactor context load flow
		// NOTE custom implementation
		app.loadContexts(db, ctx)

		if err := app.WasmKeeper.InitializePinnedCodes(ctx); err != nil {
			tmos.Exit(fmt.Sprintf("failed initialize pinned codes %s", err))
		}

		// Initialize and seal the capability keeper so all persistent capabilities
		// are loaded in-memory and prevent any further modules from creating scoped
		// sub-keepers.
		// This must be done during creation of baseapp rather than in InitChain so
		// that in-memory capabilities get regenerated on app restart.
		// Note that since this reads from the store, we can only perform it when
		// `loadLatest` is set to true.
		app.CapabilityKeeper.Seal()
	}

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(simulationModules(app, encodingConfig, skipGenesisInvariants)...)

	app.sm.RegisterStoreDecoders()

	return app
}

func (app *App) setPostHandler() {
	postHandler, err := posthandler.NewPostHandler(
		posthandler.HandlerOptions{},
	)
	if err != nil {
		panic(err)
	}

	app.SetPostHandler(postHandler)
}

// Name returns the name of the App
func (app *App) Name() string {
	return app.BaseApp.Name()
}

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.ModuleManager.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.ModuleManager.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	// custom initialization
	// var bankGenesisState banktypes.GenesisState
	// app.appCodec.MustUnmarshalJSON(genesisState["bank"], &bankGenesisState)
	// app.BandwidthMeter.AddToDesirableBandwidth(ctx, bankGenesisState.Supply.AmountOf(ctypes.VOLT).Uint64())

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap())
	resp := app.ModuleManager.InitGenesis(ctx, app.appCodec, genesisState)

	// custom initialization
	for _, account := range app.AccountKeeper.GetAllAccounts(ctx) {
		app.CyberbankKeeper.InitializeStakeAmpere(
			account.GetAccountNumber(),
			uint64(app.CyberbankKeeper.Proxy.GetAccountTotalStakeAmper(ctx, account.GetAddress())),
		)
	}

	return resp
}

// LoadHeight loads a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range keepers.GetMaccPerms() {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.aminoCodec
}

// AppCodec returns Juno's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Juno's InterfaceRegistry
func (app *App) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

func (app *App) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// configure store loader that checks if version == upgradeHeight and applies store upgrades
func (app *App) setupUpgradeStoreLoaders() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic("failed to read upgrade info from disk" + err.Error())
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, upgrade := range Upgrades {
		storeUpgrades := upgrade.StoreUpgrades
		if upgradeInfo.Name == upgrade.UpgradeName {
			app.SetStoreLoader(
				upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades),
			)
		}
	}
}

func (app *App) setupUpgradeHandlers(cfg module.Configurator) {
	for _, upgrade := range Upgrades {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(
				app.ModuleManager,
				cfg,
				&app.AppKeepers,
			),
		)
	}
}

// SimulationManager implements the SimulationApp interface
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}

// MakeCodecs constructs the *std.Codec and *codec.LegacyAmino instances used by
// Cyber's app. It is useful for tests and clients who do not want to construct the
// full Cyber app
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	config := MakeEncodingConfig()
	return config.Marshaler, config.Amino
}

func (app *App) loadContexts(db dbm.DB, ctx sdk.Context) {
	freshCtx := ctx.WithBlockHeight(int64(app.RankKeeper.GetLatestBlockNumber(ctx)))
	start := time.Now()
	app.BaseApp.Logger().Info("Loading the brain state")

	if app.LastBlockHeight() >= 1 {
		calculationPeriod := app.RankKeeper.GetParams(freshCtx).CalculationPeriod
		// TODO remove this after upgrade to v4 because on network upgrade block cannot access rank params
		if calculationPeriod == 0 {
			calculationPeriod = int64(5)
		}

		rankRoundBlockNumber := (freshCtx.BlockHeight() / calculationPeriod) * calculationPeriod
		if rankRoundBlockNumber == 0 && freshCtx.BlockHeight() >= 1 {
			rankRoundBlockNumber = 1
		}

		rankCtx, err := utils.NewContextWithMSVersion(db, rankRoundBlockNumber, app.GetKVStoreKey())
		if err != nil {
			tmos.Exit(err.Error())
		}

		app.CyberbankKeeper.LoadState(rankCtx, freshCtx)
		// TODO update index state load to one context as we store cyberlink' block now
		app.IndexKeeper.LoadState(rankCtx, freshCtx)
		app.BandwidthMeter.LoadState(freshCtx)
		app.GraphKeeper.LoadNeudeg(rankCtx, freshCtx)
		app.RankKeeper.LoadState(freshCtx)
		app.RankKeeper.StartRankCalculation(freshCtx)
	} else {
		// genesis case
		app.CyberbankKeeper.LoadState(freshCtx, freshCtx)
		// TODO update index state load to one context as we store cyberlink' block now
		app.IndexKeeper.LoadState(freshCtx, freshCtx)
		app.BandwidthMeter.InitState()
		app.GraphKeeper.LoadNeudeg(freshCtx, freshCtx)
		app.RankKeeper.LoadState(freshCtx)
	}

	app.BaseApp.Logger().Info(
		"Cyber Consensus Supercomputer is started!",
		"duration", time.Since(start).String(),
	)
}
