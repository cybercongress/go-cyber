package app

import (
	simappparams "cosmossdk.io/simapp/params"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
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
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
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
	"github.com/cosmos/cosmos-sdk/x/nft"
	nftmodule "github.com/cosmos/cosmos-sdk/x/nft/module"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	pfmrouter "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward"
	pfmroutertypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	solomachine "github.com/cosmos/ibc-go/v7/modules/light-clients/06-solomachine"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/cybercongress/go-cyber/v4/x/clock"
	"github.com/cybercongress/go-cyber/v4/x/tokenfactory"
	tokenfactorytypes "github.com/cybercongress/go-cyber/v4/x/tokenfactory/types"

	"github.com/cybercongress/go-cyber/v4/x/liquidity"
	liquiditytypes "github.com/cybercongress/go-cyber/v4/x/liquidity/types"

	"github.com/cybercongress/go-cyber/v4/x/bandwidth"
	bandwidthtypes "github.com/cybercongress/go-cyber/v4/x/bandwidth/types"
	"github.com/cybercongress/go-cyber/v4/x/cyberbank"
	cyberbanktypes "github.com/cybercongress/go-cyber/v4/x/cyberbank/types"
	"github.com/cybercongress/go-cyber/v4/x/dmn"
	dmntypes "github.com/cybercongress/go-cyber/v4/x/dmn/types"
	"github.com/cybercongress/go-cyber/v4/x/graph"
	graphtypes "github.com/cybercongress/go-cyber/v4/x/graph/types"
	grid "github.com/cybercongress/go-cyber/v4/x/grid"
	gridtypes "github.com/cybercongress/go-cyber/v4/x/grid/types"
	"github.com/cybercongress/go-cyber/v4/x/rank"
	ranktypes "github.com/cybercongress/go-cyber/v4/x/rank/types"
	"github.com/cybercongress/go-cyber/v4/x/resources"
	resourcestypes "github.com/cybercongress/go-cyber/v4/x/resources/types"
	stakingwrap "github.com/cybercongress/go-cyber/v4/x/staking"

	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	clocktypes "github.com/cybercongress/go-cyber/v4/x/clock/types"
)

// ModuleBasics TODO add notes which modules have functional blockers
// ModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distr.AppModuleBasic{},
	gov.NewAppModuleBasic(getGovProposalHandlers()),
	sdkparams.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	vesting.AppModuleBasic{},
	pfmrouter.AppModuleBasic{},
	nftmodule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	ibcfee.AppModuleBasic{},
	ica.AppModuleBasic{},
	transfer.AppModuleBasic{},
	consensus.AppModuleBasic{},
	liquidity.AppModuleBasic{},
	wasm.AppModuleBasic{},
	bandwidth.AppModuleBasic{},
	cyberbank.AppModuleBasic{},
	graph.AppModuleBasic{},
	rank.AppModuleBasic{},
	grid.AppModuleBasic{},
	dmn.AppModuleBasic{},
	resources.AppModuleBasic{},
	tokenfactory.AppModuleBasic{},
	clock.AppModuleBasic{},
	// https://github.com/cosmos/ibc-go/blob/main/docs/docs/05-migrations/08-v6-to-v7.md
	ibctm.AppModuleBasic{},
	solomachine.AppModuleBasic{},
)

func appModules(
	app *App,
	encodingConfig simappparams.EncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	appCodec := encodingConfig.Codec

	return []module.AppModule{
		genutil.NewAppModule(
			app.AccountKeeper,
			app.StakingKeeper,
			app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.AccountKeeper, app.CyberbankKeeper.Proxy),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.FeeGrantKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, &app.GovKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		sdkparams.NewAppModule(app.ParamsKeeper),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.interfaceRegistry),
		nftmodule.NewAppModule(appCodec, app.AppKeepers.NFTKeeper, app.AppKeepers.AccountKeeper, app.CyberbankKeeper.Proxy, app.interfaceRegistry),
		crisis.NewAppModule(app.AppKeepers.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		consensus.NewAppModule(appCodec, app.AppKeepers.ConsensusParamsKeeper),

		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		liquidity.NewAppModule(appCodec, app.LiquidityKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.DistrKeeper, app.GetSubspace(liquiditytypes.ModuleName)),
		cyberbank.NewAppModule(appCodec, app.CyberbankKeeper),
		bandwidth.NewAppModule(appCodec, app.AccountKeeper, app.BandwidthMeter, app.GetSubspace(bandwidthtypes.ModuleName)),
		graph.NewAppModule(
			appCodec, app.GraphKeeper, app.IndexKeeper,
			app.AccountKeeper, app.CyberbankKeeper, app.BandwidthMeter,
		),
		rank.NewAppModule(appCodec, app.RankKeeper, app.GetSubspace(ranktypes.ModuleName)),
		grid.NewAppModule(appCodec, app.GridKeeper, app.GetSubspace(gridtypes.ModuleName)),
		dmn.NewAppModule(appCodec, *app.DmnKeeper, app.GetSubspace(dmntypes.ModuleName)),
		resources.NewAppModule(appCodec, app.ResourcesKeeper, app.GetSubspace(resourcestypes.ModuleName)),
		stakingwrap.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.CyberbankKeeper.Proxy, app.GetSubspace(stakingtypes.ModuleName)),
		// NOTE add Bank Proxy here to resolve issue when new neuron is created during token-factory token transfer
		// TODO add storage listener to update neurons memory index out of cyberbank proxy
		tokenfactory.NewAppModule(app.AppKeepers.TokenFactoryKeeper, app.AppKeepers.AccountKeeper, app.CyberbankKeeper.Proxy, app.GetSubspace(tokenfactorytypes.ModuleName)),
		clock.NewAppModule(appCodec, app.AppKeepers.ClockKeeper),

		ibc.NewAppModule(app.IBCKeeper),
		transfer.NewAppModule(app.TransferKeeper),
		ibcfee.NewAppModule(app.IBCFeeKeeper),
		ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),
		pfmrouter.NewAppModule(app.PacketForwardKeeper, app.GetSubspace(pfmroutertypes.ModuleName)),
	}
}

// simulationModules returns modules for simulation manager
// define the order of the modules for deterministic simulationss
func simulationModules(
	app *App,
	encodingConfig simappparams.EncodingConfig,
	_ bool,
) []module.AppModuleSimulation {
	appCodec := encodingConfig.Codec

	return []module.AppModuleSimulation{
		auth.NewAppModule(appCodec, app.AppKeepers.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		bank.NewAppModule(appCodec, app.AppKeepers.BankKeeper, app.AppKeepers.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.AppKeepers.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.AppKeepers.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AppKeepers.AuthzKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, &app.AppKeepers.GovKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, app.AppKeepers.MintKeeper, app.AppKeepers.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		stakingwrap.NewAppModule(appCodec, app.AppKeepers.StakingKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.AppKeepers.DistrKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.AppKeepers.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.AppKeepers.SlashingKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.AppKeepers.StakingKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		sdkparams.NewAppModule(app.AppKeepers.ParamsKeeper),
		evidence.NewAppModule(app.AppKeepers.EvidenceKeeper),
		wasm.NewAppModule(appCodec, &app.AppKeepers.WasmKeeper, app.AppKeepers.StakingKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		ibc.NewAppModule(app.AppKeepers.IBCKeeper),
		transfer.NewAppModule(app.AppKeepers.TransferKeeper),
		ibcfee.NewAppModule(app.IBCFeeKeeper),
		ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),
	}
}

// orderBeginBlockers tell the app's module manager how to set the order of
// BeginBlockers, which are run at the beginning of every block.
// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
func orderBeginBlockers() []string {
	return []string{
		// upgrades should be run first
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		pfmroutertypes.ModuleName,
		ibcfeetypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		nft.ModuleName,
		// additional modules
		liquiditytypes.ModuleName,
		dmntypes.ModuleName,
		clocktypes.ModuleName,
		bandwidthtypes.ModuleName,
		cyberbanktypes.ModuleName,
		graphtypes.ModuleName,
		gridtypes.ModuleName,
		ranktypes.ModuleName,
		resourcestypes.ModuleName,
		tokenfactorytypes.ModuleName,
		wasm.ModuleName,
	}
}

func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		pfmroutertypes.ModuleName,
		capabilitytypes.ModuleName,
		ibcfeetypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		nft.ModuleName,
		// additional modules
		clocktypes.ModuleName,
		tokenfactorytypes.ModuleName,
		dmntypes.ModuleName,
		gridtypes.ModuleName,
		resourcestypes.ModuleName,
		liquiditytypes.ModuleName,
		wasm.ModuleName,
		// TODO check end blocks
		cyberbanktypes.ModuleName,
		bandwidthtypes.ModuleName,
		graphtypes.ModuleName,
		ranktypes.ModuleName,
	}
}

// NOTE: The genutils module must occur after staking so that pools are
// properly initialized with tokens from genesis accounts.
// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
// NOTE: Capability module must occur first so that it can initialize any capabilities
// so that other modules that want to create or claim capabilities afterwards in InitChain
// can do so safely.
// NOTE: wasm module should be at the end as it can call other module functionality direct or via message dispatching during
// genesis phase. For example bank transfer, auth account check, staking, ...
func orderInitBlockers() []string {
	return []string{
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		pfmroutertypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		nft.ModuleName,
		consensusparamtypes.ModuleName,
		// additional modules
		liquiditytypes.ModuleName,
		ibcfeetypes.ModuleName,
		clocktypes.ModuleName,
		tokenfactorytypes.ModuleName,
		bandwidthtypes.ModuleName,
		ranktypes.ModuleName,
		gridtypes.ModuleName,
		resourcestypes.ModuleName,
		dmntypes.ModuleName,
		graphtypes.ModuleName,
		// NOTE: cyberbank will be initialized directly in InitChainer
		cyberbanktypes.ModuleName,
		wasm.ModuleName,
	}
}
