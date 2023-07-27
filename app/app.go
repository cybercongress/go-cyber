package app

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmclient "github.com/CosmWasm/wasmd/x/wasm/client"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	ibcclientclient "github.com/cosmos/ibc-go/v4/modules/core/02-client/client"
	"github.com/cybercongress/go-cyber/app/keepers"
	"github.com/cybercongress/go-cyber/app/params"
	upgrades "github.com/cybercongress/go-cyber/app/upgrades"
	_ "github.com/cybercongress/go-cyber/client/docs/statik"
	ctypes "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/utils"
	cyberbanktypes "github.com/cybercongress/go-cyber/x/cyberbank/types"
	"github.com/spf13/cast"
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const (
	appName = "BostromHub"
	// upgradeName = "cyberfrey" // TODO move to upgrades
)

// We pull these out so we can set them with LDFLAGS in the Makefile
var (
	NodeDir      = ".cyber"
	Bech32Prefix = "bostrom"

	// If EnabledSpecificProposals is "", and this is "true", then enable all x/wasm proposals.
	ProposalsEnabled        = "true"
	EnableSpecificProposals = ""

	Upgrades = []upgrades.Upgrade{}
)

// These constants are derived from the above variables.
// These are the ones we will want to use in the code, based on
// any overrides above
var (
	// DefaultNodeHome default home directories for Cyber
	DefaultNodeHome = os.ExpandEnv("$HOME/") + NodeDir
)

// GetEnabledProposals parses the ProposalsEnabled / EnableSpecificProposals values to
// produce a list of enabled proposals to pass into wasmd app.
func GetEnabledProposals() []wasm.ProposalType {
	if EnableSpecificProposals == "" {
		if ProposalsEnabled == "true" {
			return wasm.EnableAllProposals
		}
		return wasm.DisableAllProposals
	}
	chunks := strings.Split(EnableSpecificProposals, ",")
	proposals, err := wasm.ConvertToProposals(chunks)
	if err != nil {
		panic(err)
	}
	return proposals
}

func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler
	govProposalHandlers = wasmclient.ProposalHandlers

	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
		distrclient.ProposalHandler,
		upgradeclient.ProposalHandler,
		upgradeclient.CancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
	)

	return govProposalHandlers
}

// module accounts that are allowed to receive tokens
var allowedReceivingModAcc = map[string]bool{
	distrtypes.ModuleName: true,
}

var (
	_ simapp.App              = (*App)(nil)
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
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// the module manager
	mm *module.Manager
	sm *module.SimulationManager
}

// New returns a reference to an initialized Cyber Consensus Computer.
func NewApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig params.EncodingConfig,
	enabledProposals []wasm.ProposalType,
	appOpts servertypes.AppOptions,
	wasmOpts []wasm.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	appCodec, legacyAmino := encodingConfig.Marshaler, encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	app := &App{
		BaseApp:           bApp,
		aminoCodec:        legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
	}

	// Setup keepers
	app.AppKeepers = keepers.NewAppKeepers(
		appCodec,
		bApp,
		legacyAmino,
		maccPerms,
		app.ModuleAccountAddrs(),
		skipUpgradeHeights,
		homePath,
		invCheckPeriod,
		enabledProposals,
		appOpts,
		wasmOpts,
	)

	// upgrade handlers
	cfg := module.NewConfigurator(appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(appModules(app, encodingConfig, skipGenesisInvariants)...)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(orderBeginBlockers()...)

	app.mm.SetOrderEndBlockers(orderEndBlockers()...)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(orderInitBlockers()...)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)

	// TODO hack to initialize application
	// from sdk.bank/module.go: AppModule RegisterServices
	// 	m := keeper.NewMigrator(am.keeper.(keeper.BaseKeeper))
	bankMigModule := bank.NewAppModule(appCodec, *app.BankKeeper.GetBankKeeper(), app.AccountKeeper)
	proxy := app.mm.Modules[banktypes.ModuleName]
	app.mm.Modules[banktypes.ModuleName] = bankMigModule
	app.mm.RegisterServices(cfg)
	app.mm.Modules[banktypes.ModuleName] = proxy

	// TODO refactor bank reinitialization flow
	// NOTE custom implementation
	// reinitializeBank(app, cfg)

	// initialize stores
	app.MountKVStores(app.GetKVStoreKey())
	app.MountTransientStores(app.GetTransientStoreKey())
	app.MountMemoryStores(app.GetMemoryStoreKey())

	// register upgrade
	app.setupUpgradeHandlers(cfg)

	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerBaseOptions: HandlerBaseOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  app.FeeGrantKeeper,
				SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			BandwidthMeter:    app.BandwidthMeter,
			IBCKeeper:         app.IBCKeeper,
			TXCounterStoreKey: app.GetKey(wasm.StoreKey),
			WasmConfig:        &wasmConfig,
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

	// UPGRADES HANDLER SECTION
	// app.UpgradeKeeper.SetUpgradeHandler(
	//	upgradeName,
	//	func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	//
	//		//ctx.Logger().Info("start to init module...")
	//		//ctx.Logger().Info("start to run module migrations...")
	//
	//		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	//	},
	//)
	//
	// upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	//if err != nil {
	//	panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	//}
	//
	//if upgradeInfo.Name == upgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
	//	storeUpgrades := store.StoreUpgrades{}
	//
	//	// configure store loader that checks if version == upgradeHeight and applies store upgrades
	//	app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	//}
	//

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(fmt.Sprintf("failed to load latest version: %s", err))
		}
		ctx := app.BaseApp.NewContext(true, tmproto.Header{})

		// TODO refactor context load flow
		// NOTE custom implementation
		app.loadContexts(db)

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

// Name returns the name of the App
func (app *App) Name() string {
	return app.BaseApp.Name()
}

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	var bankGenesisState banktypes.GenesisState
	app.appCodec.MustUnmarshalJSON(genesisState["bank"], &bankGenesisState)
	app.BandwidthMeter.AddToDesirableBandwidth(ctx, bankGenesisState.Supply.AmountOf(ctypes.VOLT).Uint64())

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	resp := app.mm.InitGenesis(ctx, app.appCodec, genesisState)

	// because manager skips init genesis for modules with empty data (e.g null)
	app.mm.Modules[cyberbanktypes.ModuleName].InitGenesis(ctx, app.appCodec, nil)

	return resp
}

// LoadHeight loads a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *App) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
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
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	// Register legacy tx routes.
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
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
		if upgradeInfo.Name == upgrade.UpgradeName {
			app.SetStoreLoader(
				upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &upgrade.StoreUpgrades),
			)
		}
	}
}

func (app *App) setupUpgradeHandlers(cfg module.Configurator) {
	for _, upgrade := range Upgrades {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(
				app.mm,
				cfg,
				&app.AppKeepers,
			),
		)
	}
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
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

func (app *App) loadContexts(db dbm.DB) {
	freshCtx := app.BaseApp.NewContext(true, tmproto.Header{})
	freshCtx = freshCtx.WithBlockHeight(int64(app.RankKeeper.GetLatestBlockNumber(freshCtx)))
	start := time.Now()
	app.BaseApp.Logger().Info("Loading the brain state")

	if app.LastBlockHeight() >= 1 {
		calculationPeriod := app.RankKeeper.GetParams(freshCtx).CalculationPeriod
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
