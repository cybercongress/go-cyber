package app

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v3/modules/core"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/gravity-devs/liquidity/x/liquidity"
	liquiditytypes "github.com/gravity-devs/liquidity/x/liquidity/types"

	"github.com/cybercongress/go-cyber/v2/app/params"
	"github.com/cybercongress/go-cyber/v2/x/bandwidth"
	bandwidthtypes "github.com/cybercongress/go-cyber/v2/x/bandwidth/types"
	"github.com/cybercongress/go-cyber/v2/x/cyberbank"
	cyberbanktypes "github.com/cybercongress/go-cyber/v2/x/cyberbank/types"
	"github.com/cybercongress/go-cyber/v2/x/dmn"
	dmntypes "github.com/cybercongress/go-cyber/v2/x/dmn/types"
	"github.com/cybercongress/go-cyber/v2/x/graph"
	graphtypes "github.com/cybercongress/go-cyber/v2/x/graph/types"
	grid "github.com/cybercongress/go-cyber/v2/x/grid"
	gridtypes "github.com/cybercongress/go-cyber/v2/x/grid/types"
	"github.com/cybercongress/go-cyber/v2/x/rank"
	ranktypes "github.com/cybercongress/go-cyber/v2/x/rank/types"
	"github.com/cybercongress/go-cyber/v2/x/resources"
	resourcestypes "github.com/cybercongress/go-cyber/v2/x/resources/types"
	stakingwrap "github.com/cybercongress/go-cyber/v2/x/staking"

	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
)

// module account permissions
var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:     nil,
	distrtypes.ModuleName:          nil,
	minttypes.ModuleName:           {authtypes.Minter},
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:            {authtypes.Burner},
	ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	wasm.ModuleName:                {authtypes.Burner},
	liquiditytypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
	gridtypes.GridPoolName:         nil,
	resourcestypes.ResourcesName:   {authtypes.Minter, authtypes.Burner},
}

// ModuleBasics TODO add notes which modules have functional blockers
// ModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distr.AppModuleBasic{},
	gov.NewAppModuleBasic(getGovProposalHandlers()...),
	sdkparams.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	liquidity.AppModuleBasic{},
	wasm.AppModuleBasic{},
	bandwidth.AppModuleBasic{},
	cyberbank.AppModuleBasic{},
	graph.AppModuleBasic{},
	rank.AppModuleBasic{},
	grid.AppModuleBasic{},
	dmn.AppModuleBasic{},
	resources.AppModuleBasic{},
)

func appModules(
	app *App,
	encodingConfig params.EncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	appCodec := encodingConfig.Marshaler

	return []module.AppModule{
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
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.FeeGrantKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.StakingKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		sdkparams.NewAppModule(app.ParamsKeeper),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.interfaceRegistry),
		transfer.NewAppModule(app.TransferKeeper),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy),
		liquidity.NewAppModule(appCodec, app.LiquidityKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.DistrKeeper),
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
		stakingwrap.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy),
	}
}

// simulationModules returns modules for simulation manager
// define the order of the modules for deterministic simulationss
func simulationModules(
	app *App,
	encodingConfig params.EncodingConfig,
	_ bool,
) []module.AppModuleSimulation {
	appCodec := encodingConfig.Marshaler

	return []module.AppModuleSimulation{
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		stakingwrap.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		sdkparams.NewAppModule(app.ParamsKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		transfer.NewAppModule(app.TransferKeeper),
	}
}

// orderBeginBlockers tell the app's module manager how to set the order of
// BeginBlockers, which are run at the beginning of every block.
func orderBeginBlockers() []string {
	return []string{
		// upgrades should be run first
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		stakingtypes.ModuleName,
		liquiditytypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		evidencetypes.ModuleName,
		dmntypes.ModuleName,

		authtypes.ModuleName,
		authz.ModuleName,
		banktypes.ModuleName,
		bandwidthtypes.ModuleName,
		crisistypes.ModuleName,
		cyberbanktypes.ModuleName,
		feegrant.ModuleName,
		genutiltypes.ModuleName,
		govtypes.ModuleName,
		graphtypes.ModuleName,
		gridtypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		paramstypes.ModuleName,
		ranktypes.ModuleName,
		resourcestypes.ModuleName,
		vestingtypes.ModuleName,
		wasm.ModuleName,
	}
}

func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		liquiditytypes.ModuleName,
		cyberbanktypes.ModuleName,
		bandwidthtypes.ModuleName,
		graphtypes.ModuleName,
		ranktypes.ModuleName,

		authtypes.ModuleName,
		authz.ModuleName,
		banktypes.ModuleName,
		capabilitytypes.ModuleName,
		distrtypes.ModuleName,
		dmntypes.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		genutiltypes.ModuleName,
		gridtypes.ModuleName,
		ibchost.ModuleName,
		ibctransfertypes.ModuleName,
		minttypes.ModuleName,
		paramstypes.ModuleName,
		resourcestypes.ModuleName,
		slashingtypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		wasm.ModuleName,
	}
}

func orderInitBlockers() []string {
	return []string{
		capabilitytypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		evidencetypes.ModuleName,
		liquiditytypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		authtypes.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		wasm.ModuleName,
		bandwidthtypes.ModuleName,
		ranktypes.ModuleName,
		gridtypes.ModuleName,
		resourcestypes.ModuleName,
		dmntypes.ModuleName,
		graphtypes.ModuleName,
		cyberbanktypes.ModuleName, // cyberbank will be initialized directly in InitChainer
	}
}
