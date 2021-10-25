package app

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/ibc-go/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/modules/core"
	ibcclient "github.com/cosmos/ibc-go/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/modules/core/keeper"
	_ "github.com/cybercongress/go-cyber/client/docs/statik"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	ctypes "github.com/cybercongress/go-cyber/types"
	"github.com/gorilla/mux"
	"github.com/gravity-devs/liquidity/x/liquidity"
	liquiditykeeper "github.com/gravity-devs/liquidity/x/liquidity/keeper"
	liquiditytypes "github.com/gravity-devs/liquidity/x/liquidity/types"
	"github.com/rakyll/statik/fs"

	wasmplugins "github.com/cybercongress/go-cyber/plugins"
	"github.com/cybercongress/go-cyber/x/dmn"
	"github.com/cybercongress/go-cyber/x/resources"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/spf13/cast"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cybercongress/go-cyber/utils"
	"github.com/cybercongress/go-cyber/x/bandwidth"
	"github.com/cybercongress/go-cyber/x/cyberbank"
	cyberbankkeeper "github.com/cybercongress/go-cyber/x/cyberbank/keeper"
	cyberbanktypes "github.com/cybercongress/go-cyber/x/cyberbank/types"
	"github.com/cybercongress/go-cyber/x/graph"

	bandwidthkeeper "github.com/cybercongress/go-cyber/x/bandwidth/keeper"
	bandwidthtypes "github.com/cybercongress/go-cyber/x/bandwidth/types"
	graphkeeper "github.com/cybercongress/go-cyber/x/graph/keeper"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
	"github.com/cybercongress/go-cyber/x/rank"
	rankkeeper "github.com/cybercongress/go-cyber/x/rank/keeper"
	ranktypes "github.com/cybercongress/go-cyber/x/rank/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmclient "github.com/CosmWasm/wasmd/x/wasm/client"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	bandwidthwasm "github.com/cybercongress/go-cyber/x/bandwidth/wasm"
	dmnwasm "github.com/cybercongress/go-cyber/x/dmn/wasm"
	graphwasm "github.com/cybercongress/go-cyber/x/graph/wasm"
	gridwasm "github.com/cybercongress/go-cyber/x/grid/wasm"
	rankwasm "github.com/cybercongress/go-cyber/x/rank/wasm"
	resourceswasm "github.com/cybercongress/go-cyber/x/resources/wasm"

	grid "github.com/cybercongress/go-cyber/x/grid"
	gridkeeper "github.com/cybercongress/go-cyber/x/grid/keeper"
	gridtypes "github.com/cybercongress/go-cyber/x/grid/types"

	dmnkeeper "github.com/cybercongress/go-cyber/x/dmn/keeper"
	dmntypes "github.com/cybercongress/go-cyber/x/dmn/types"

	resourceskeeper "github.com/cybercongress/go-cyber/x/resources/keeper"
	resourcestypes "github.com/cybercongress/go-cyber/x/resources/types"
	stakingwrap "github.com/cybercongress/go-cyber/x/staking"

	"github.com/cybercongress/go-cyber/app/params"
)

const (
	appName = "BostromHub"
)

// We pull these out so we can set them with LDFLAGS in the Makefile
var (
	NodeDir      = ".cyber"
	Bech32Prefix = "bostrom"

	// DefaultBondDenom is the denomination of coin to use for bond/staking
	DefaultBondDenom = "boot"
	// DefaultFeeDenom is the denomination of coin to use for fees
	DefaultFeeDenom = "boot"
	// DefaultReDnmString is the allowed denom regex expression
	DefaultReDnmString = `[a-zA-Z][a-zA-Z0-9/\-\.]{2,127}`

	// If EnabledSpecificProposals is "", and this is "true", then enable all x/wasm proposals.
	// If EnabledSpecificProposals is "", and this is not "true", then disable all x/wasm proposals.
	ProposalsEnabled = "true"
	// If set to non-empty string it must be comma-separated list of values that are all a subset
	// of "EnableAllProposals" (takes precedence over ProposalsEnabled)
	// https://github.com/CosmWasm/wasmd/blob/02a54d33ff2c064f3539ae12d75d027d9c665f05/x/wasm/internal/types/proposal.go#L28-L34
	//EnableSpecificProposals = "StoreCodeProposal,InstantiateContractProposal,MigrateContractProposal,ClearAdminProposal,PinCodes,UnpinCodes"
	EnableSpecificProposals = "PinCodes,UnpinCodes"
)

// GetEnabledProposals parses the ProposalsEnabled / EnableSpecificProposals values to
// produce a list of enabled proposals to pass into cyber app.
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

// These constants are derived from the above variables.
// These are the ones we will want to use in the code, based on
// any overrides above
var (
	// DefaultNodeHome default home directories for wasmd
	DefaultNodeHome = os.ExpandEnv("$HOME/") + NodeDir
)

var (
	//DefaultPowerReduction = sdk.NewIntFromUint64(1000000000)

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(append(
			wasmclient.ProposalHandlers,
			paramsclient.ProposalHandler,
			distrclient.ProposalHandler,
			upgradeclient.ProposalHandler,
			upgradeclient.CancelProposalHandler,
			ibcclientclient.UpdateClientProposalHandler,
			ibcclientclient.UpgradeProposalHandler,
		)...,
		),
		sdkparams.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		vesting.AppModuleBasic{},

		ibc.AppModuleBasic{},
		transfer.AppModuleBasic{},

		bandwidth.AppModuleBasic{},
		cyberbank.AppModuleBasic{},
		graph.AppModuleBasic{},
		rank.AppModuleBasic{},
		grid.AppModuleBasic{},
		dmn.AppModuleBasic{},
		resources.AppModuleBasic{},
		wasm.AppModuleBasic{},
		liquidity.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		gridtypes.GridPoolName:         nil,
		resourcestypes.ResourcesName:   {authtypes.Minter, authtypes.Burner},
		wasm.ModuleName:                {authtypes.Burner},
		liquiditytypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
	}

    // module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName: true,
	}
)

var (
	_ CosmosApp               = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

// SdkCoinDenomRegex returns a new sdk base denom regex string
func SdkCoinDenomRegex() string {
	return DefaultReDnmString
}


// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp

	appName string

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	EvidenceKeeper   evidencekeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	AuthzKeeper      authzkeeper.Keeper

	IBCKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	TransferKeeper   ibctransferkeeper.Keeper

	BandwidthMeter   *bandwidthkeeper.BandwidthMeter
	CyberbankKeeper  *cyberbankkeeper.IndexedKeeper
	GraphKeeper      *graphkeeper.GraphKeeper
	IndexKeeper      *graphkeeper.IndexKeeper
	RankKeeper 		 *rankkeeper.StateKeeper
	GridKeeper 		 gridkeeper.Keeper
	DmnKeeper  		 *dmnkeeper.Keeper
	ResourcesKeeper  resourceskeeper.Keeper
	WasmKeeper       wasm.Keeper
	LiquidityKeeper  liquiditykeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper

	// the module manager
	mm *module.Manager

	configurator module.Configurator
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
	wasmOpts []wasm.Option,
	enabledProposals []wasm.ProposalType,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	config := sdk.NewConfig()
	config.Seal()

	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	sdk.SetCoinDenomRegex(SdkCoinDenomRegex)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, liquiditytypes.StoreKey, ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey, feegrant.StoreKey, authzkeeper.StoreKey,
		bandwidthtypes.StoreKey,
		graphtypes.StoreKey,
		ranktypes.StoreKey,
		gridtypes.StoreKey,
		dmntypes.StoreKey,
		wasm.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(
		paramstypes.TStoreKey,
		graphtypes.TStoreKey,
	)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &App{
		BaseApp:           bApp,
		appName:           appName,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)

	// grant capabilities for the ibc and ibc-transfer modules
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedWasmKeeper := app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		app.AccountKeeper,
		app.GetSubspace(banktypes.ModuleName),
		app.ModuleAccountAddrs(),
	)
	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		app.BaseApp.MsgServiceRouter(),
	)
	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		keys[feegrant.StoreKey],
		app.AccountKeeper,
	)
	app.CyberbankKeeper = cyberbankkeeper.NewIndexedKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		cyberbankkeeper.Wrap(app.BankKeeper),
		app.AccountKeeper,
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingtypes.StoreKey],
		app.AccountKeeper,
		app.CyberbankKeeper.Proxy,
		app.GetSubspace(stakingtypes.ModuleName),
	)
	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		keys[minttypes.StoreKey],
		app.GetSubspace(minttypes.ModuleName),
		&stakingKeeper,
		app.AccountKeeper,
		app.CyberbankKeeper.Proxy,
		authtypes.FeeCollectorName,
	)
	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		keys[distrtypes.StoreKey],
		app.GetSubspace(distrtypes.ModuleName),
		app.AccountKeeper,
		app.CyberbankKeeper.Proxy,
		&stakingKeeper,
		authtypes.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		keys[slashingtypes.StoreKey],
		&stakingKeeper,
		app.GetSubspace(slashingtypes.ModuleName),
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		app.CyberbankKeeper.Proxy,
		authtypes.FeeCollectorName,
	)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.DistrKeeper.Hooks(),
			app.SlashingKeeper.Hooks(),
		),
	)

	app.CyberbankKeeper.Proxy.AddHook(
		bandwidth.CollectAddressesWithStakeChange(),
	)

	app.BandwidthMeter = bandwidthkeeper.NewBandwidthMeter(
		app.appCodec,
		keys[bandwidthtypes.StoreKey],
		app.CyberbankKeeper.Proxy,
		app.GetSubspace(bandwidthtypes.ModuleName),
	)

	app.GraphKeeper = graphkeeper.NewKeeper(
		appCodec,
		keys[graphtypes.ModuleName],
		tkeys[paramstypes.TStoreKey],
	)

	app.IndexKeeper = graphkeeper.NewIndexKeeper(
		*app.GraphKeeper,
		tkeys[paramstypes.TStoreKey],
	)

	computeUnit := ranktypes.ComputeUnit(cast.ToInt(appOpts.Get(rank.FlagComputeGPU)))
	searchAPI := cast.ToBool(appOpts.Get(rank.FlagSearchAPI))

	app.RankKeeper = rankkeeper.NewKeeper(
		keys[ranktypes.ModuleName],
		app.GetSubspace(ranktypes.ModuleName),
		searchAPI,
		app.CyberbankKeeper,
		app.IndexKeeper,
		app.GraphKeeper,
		app.AccountKeeper,
		computeUnit,
	)

	app.GridKeeper = gridkeeper.NewKeeper(
		appCodec,
		keys[gridtypes.ModuleName],
		app.CyberbankKeeper.Proxy,
		app.AccountKeeper,
		app.GetSubspace(gridtypes.ModuleName),
	)
	app.CyberbankKeeper.Proxy.SetGridKeeper(&app.GridKeeper)
	app.CyberbankKeeper.Proxy.SetAccountKeeper(app.AccountKeeper)

	app.ResourcesKeeper = resourceskeeper.NewKeeper(
		appCodec,
		app.AccountKeeper,
		app.CyberbankKeeper.Proxy,
		app.BandwidthMeter,
		app.GetSubspace(resourcestypes.ModuleName),
	)

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibchost.StoreKey],
		app.GetSubspace(ibchost.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.CyberbankKeeper.Proxy,
		scopedTransferKeeper,
	)

	app.DmnKeeper = dmnkeeper.NewKeeper(
		appCodec,
		keys[dmntypes.StoreKey],
		app.CyberbankKeeper.Proxy,
		app.AccountKeeper,
		app.GetSubspace(dmntypes.ModuleName),
	)

	// Initialize CosmWasm
	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	// Initialize Cyber's query integrations
	querier := wasmplugins.NewQuerier()
	queries := map[string]wasmplugins.WasmQuerierInterface{
		wasmplugins.WasmQueryRouteRank:      rankwasm.NewWasmQuerier(app.RankKeeper),
		wasmplugins.WasmQueryRouteGraph:     graphwasm.NewWasmQuerier(*app.GraphKeeper),
		wasmplugins.WasmQueryRouteDmn:       dmnwasm.NewWasmQuerier(*app.DmnKeeper),
		wasmplugins.WasmQueryRouteGrid:      gridwasm.NewWasmQuerier(app.GridKeeper),
		wasmplugins.WasmQueryRouteBandwidth: bandwidthwasm.NewWasmQuerier(app.BandwidthMeter),
	}
	querier.Queriers = queries
	queryPlugins := &wasm.QueryPlugins{
		Custom: querier.QueryCustom,
	}

	// Initialize Cyber's encoder integrations
	parser := wasmplugins.NewMsgParser()
	parsers := map[string]wasmplugins.WasmMsgParserInterface{
		wasmplugins.WasmMsgParserRouteGraph:     graphwasm.NewWasmMsgParser(),
		wasmplugins.WasmMsgParserRouteDmn:       dmnwasm.NewWasmMsgParser(),
		wasmplugins.WasmMsgParserRouteGrid:      gridwasm.NewWasmMsgParser(),
		wasmplugins.WasmMsgParserRouteResources: resourceswasm.NewWasmMsgParser(),
	}
	parser.Parsers = parsers
	customEncoders := &wasm.MessageEncoders{
		Custom: parser.ParseCustom,
	}

	supportedFeatures := "iterator,staking,stargate,cyber"
	wasmOpts = append(wasmOpts, wasmkeeper.WithQueryPlugins(queryPlugins))
	wasmOpts = append(wasmOpts, wasmkeeper.WithMessageEncoders(customEncoders))

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	app.WasmKeeper = wasm.NewKeeper(
		appCodec,
		keys[wasm.StoreKey],
		app.GetSubspace(wasm.ModuleName),
		app.AccountKeeper,
		app.CyberbankKeeper.Proxy,
		app.StakingKeeper,
		app.DistrKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		app.TransferKeeper,
		app.Router(),
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		wasmOpts...,
	)

	app.DmnKeeper.SetWasmKeeper(app.WasmKeeper)

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, sdkparams.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))

	if len(enabledProposals) != 0 {
		govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.WasmKeeper, enabledProposals))
	}

	app.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		app.GetSubspace(govtypes.ModuleName),
		app.AccountKeeper,
		app.CyberbankKeeper.Proxy,
		&stakingKeeper,
		govRouter,
	)

	transferModule := transfer.NewAppModule(app.TransferKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(
		ibctransfertypes.ModuleName,
		transferModule,
	)
	ibcRouter.AddRoute(
		wasm.ModuleName,
		wasm.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper),
	)
	app.IBCKeeper.SetRouter(ibcRouter)

	app.LiquidityKeeper = liquiditykeeper.NewKeeper(
		appCodec,
		keys[liquiditytypes.StoreKey],
		app.GetSubspace(liquiditytypes.ModuleName),
		app.CyberbankKeeper.Proxy,
		app.AccountKeeper,
		app.DistrKeeper,
	)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		keys[evidencetypes.StoreKey],
		&app.StakingKeeper,
		app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	var skipGenesisInvariants = cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper,
			app.StakingKeeper,
			app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		vesting.NewAppModule(app.AccountKeeper, app.CyberbankKeeper.Proxy),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.StakingKeeper),
		stakingwrap.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.interfaceRegistry),
		sdkparams.NewAppModule(app.ParamsKeeper),

		ibc.NewAppModule(app.IBCKeeper),
		transferModule,

		liquidity.NewAppModule(appCodec, app.LiquidityKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.DistrKeeper),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper),

		cyberbank.NewAppModule(appCodec, app.CyberbankKeeper),
		bandwidth.NewAppModule(appCodec, app.AccountKeeper, app.BandwidthMeter),
		graph.NewAppModule(
			appCodec, app.GraphKeeper, app.IndexKeeper,
			app.AccountKeeper, app.CyberbankKeeper, app.BandwidthMeter,
		),
		rank.NewAppModule(appCodec, app.RankKeeper),
		grid.NewAppModule(appCodec, app.GridKeeper),
		dmn.NewAppModule(appCodec, *app.DmnKeeper),
		resources.NewAppModule(appCodec, app.ResourcesKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		liquiditytypes.ModuleName,
		ibchost.ModuleName,
		dmntypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		liquiditytypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		cyberbanktypes.ModuleName,
		bandwidthtypes.ModuleName,
		graphtypes.ModuleName,
		ranktypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		ibchost.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		liquiditytypes.ModuleName,
		ibctransfertypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,

		// graph and cyberbank will be initialized directly in InitChainer
		wasm.ModuleName,
		bandwidthtypes.ModuleName,
		ranktypes.ModuleName,
		gridtypes.ModuleName,
		resourcestypes.ModuleName,
		dmntypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())

	// TODO This needs to be refactored
	// Bank module fails on register services using passed Proxy
	// register all other modules services and then manually register bank's things
	delete(app.mm.Modules, "bank")
	app.mm.RegisterServices(app.configurator)
	app.mm.Modules["bank"] = bank.NewAppModule(appCodec, app.CyberbankKeeper.Proxy, app.AccountKeeper)

	banktypes.RegisterMsgServer(app.configurator.MsgServer(), bankkeeper.NewMsgServerImpl(app.CyberbankKeeper.Proxy))
	banktypes.RegisterQueryServer(app.configurator.QueryServer(), app.CyberbankKeeper.Proxy)
	// migration wouldn't be called because bank's consensus version is 2
	m := bankkeeper.NewMigrator(app.BankKeeper.(bankkeeper.BaseKeeper))
	app.configurator.RegisterMigration(banktypes.ModuleName, 1, m.Migrate1to2)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerBaseOptions: HandlerBaseOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.CyberbankKeeper.Proxy,
				FeegrantKeeper:  app.FeeGrantKeeper,
				BandwidthMeter:  app.BandwidthMeter,
				SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			IBCChannelkeeper:    app.IBCKeeper.ChannelKeeper,
			txCounterStoreKey:   keys[wasm.StoreKey],
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	app.SetAnteHandler(anteHandler)
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(fmt.Sprintf("failed to load latest version: %s", err))
		}

		// TODO refactor load flow as updated state sync in sdk v45.X will be released
		freshCtx := app.BaseApp.NewContext(true, tmproto.Header{})
		freshCtx = freshCtx.WithBlockHeight(int64(app.RankKeeper.GetLatestBlockNumber(freshCtx)))

		bandwidthParams := bandwidthtypes.DefaultParams()
		rankParams := ranktypes.DefaultParams()
		if app.LastBlockHeight() >= 1 {
			bandwidthParams = app.BandwidthMeter.GetParams(freshCtx)
			rankParams = app.RankKeeper.GetParams(freshCtx)
		}
		app.BandwidthMeter.SetParams(freshCtx, bandwidthParams)
		app.RankKeeper.SetParams(freshCtx, rankParams)

		calculationPeriod := app.RankKeeper.GetParams(freshCtx).CalculationPeriod
		rankRoundBlockNumber := (freshCtx.BlockHeight() / calculationPeriod) * calculationPeriod
		if rankRoundBlockNumber == 0 && freshCtx.BlockHeight() >= 1 {
			rankRoundBlockNumber = 1
		}

		rankCtx, err := utils.NewContextWithMSVersion(db, rankRoundBlockNumber, keys)
		if err != nil {
			tmos.Exit(err.Error())
		}

		start := time.Now()
		app.BaseApp.Logger().Info("Loading the brain state")
		app.CyberbankKeeper.LoadState(rankCtx, freshCtx)
		// TODO update index state load to one context as we store cyberlink' block now
		app.IndexKeeper.LoadState(rankCtx, freshCtx)
		app.BandwidthMeter.LoadState(freshCtx)
		app.GraphKeeper.LoadNeudeg(rankCtx, freshCtx)
		app.RankKeeper.LoadState(freshCtx)

		app.CapabilityKeeper.Seal()

		app.BaseApp.Logger().Info(
			"Cyber Consensus Supercomputer is started!",
			"duration", time.Since(start).String(),
		)
	}


	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	return app
}

// Name returns the name of the App
func (app *App) Name() string { return app.BaseApp.Name() }

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
	app.legacyAmino.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

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
	return app.legacyAmino
}

// AppCodec returns Cyber's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Cyber's InterfaceRegistry
func (app *App) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *App) GetMemKey(storeKey string) *sdk.MemoryStoreKey {
	return app.memKeys[storeKey]
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

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(apiSvr.Router)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app *App) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)

	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)

	paramsKeeper.Subspace(liquiditytypes.ModuleName)

	paramsKeeper.Subspace(wasm.ModuleName)
	paramsKeeper.Subspace(bandwidthtypes.ModuleName)
	paramsKeeper.Subspace(ranktypes.ModuleName)
	paramsKeeper.Subspace(gridtypes.ModuleName)
	paramsKeeper.Subspace(dmntypes.ModuleName)
	paramsKeeper.Subspace(resourcestypes.ModuleName)

	return paramsKeeper
}

// MakeCodecs constructs the *std.Codec and *codec.LegacyAmino instances used by
// Cyber's app. It is useful for tests and clients who do not want to construct the
// full Cyber app
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	config := MakeEncodingConfig()
	return config.Marshaler, config.Amino
}