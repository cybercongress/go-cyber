package keepers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	liquiditykeeper "github.com/gravity-devs/liquidity/x/liquidity/keeper"
	liquiditytypes "github.com/gravity-devs/liquidity/x/liquidity/types"
	"github.com/spf13/cast"

	wasmplugins "github.com/cybercongress/go-cyber/v3/plugins"
	"github.com/cybercongress/go-cyber/v3/x/bandwidth"
	bandwidthkeeper "github.com/cybercongress/go-cyber/v3/x/bandwidth/keeper"
	bandwidthtypes "github.com/cybercongress/go-cyber/v3/x/bandwidth/types"
	cyberbankkeeper "github.com/cybercongress/go-cyber/v3/x/cyberbank/keeper"
	dmnkeeper "github.com/cybercongress/go-cyber/v3/x/dmn/keeper"
	dmntypes "github.com/cybercongress/go-cyber/v3/x/dmn/types"
	graphkeeper "github.com/cybercongress/go-cyber/v3/x/graph/keeper"
	graphtypes "github.com/cybercongress/go-cyber/v3/x/graph/types"
	gridkeeper "github.com/cybercongress/go-cyber/v3/x/grid/keeper"
	gridtypes "github.com/cybercongress/go-cyber/v3/x/grid/types"
	"github.com/cybercongress/go-cyber/v3/x/rank"
	rankkeeper "github.com/cybercongress/go-cyber/v3/x/rank/keeper"
	ranktypes "github.com/cybercongress/go-cyber/v3/x/rank/types"
	resourceskeeper "github.com/cybercongress/go-cyber/v3/x/resources/keeper"
	resourcestypes "github.com/cybercongress/go-cyber/v3/x/resources/types"
)

type AppKeepers struct {
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
	IBCKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper   evidencekeeper.Keeper
	TransferKeeper   ibctransferkeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	AuthzKeeper      authzkeeper.Keeper

	WasmKeeper      wasm.Keeper
	LiquidityKeeper liquiditykeeper.Keeper
	BandwidthMeter  *bandwidthkeeper.BandwidthMeter
	CyberbankKeeper *cyberbankkeeper.IndexedKeeper
	GraphKeeper     *graphkeeper.GraphKeeper
	IndexKeeper     *graphkeeper.IndexKeeper
	RankKeeper      *rankkeeper.StateKeeper
	GridKeeper      gridkeeper.Keeper
	DmnKeeper       *dmnkeeper.Keeper
	ResourcesKeeper resourceskeeper.Keeper

	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedIBCFeeKeeper   capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper     capabilitykeeper.ScopedKeeper
}

func NewAppKeepers(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	cdc *codec.LegacyAmino,
	maccPerms map[string][]string,
	blockedAddress map[string]bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	enabledProposals []wasm.ProposalType,
	appOpts servertypes.AppOptions,
	wasmOpts []wasm.Option,
) AppKeepers {
	appKeepers := AppKeepers{}

	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.GenerateKeys()

	appKeepers.ParamsKeeper = initParamsKeeper(
		appCodec,
		cdc,
		appKeepers.keys[paramstypes.StoreKey],
		appKeepers.tkeys[paramstypes.TStoreKey],
	)

	// set the BaseApp's parameter store
	bApp.SetParamStore(appKeepers.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		appKeepers.keys[capabilitytypes.StoreKey],
		appKeepers.memKeys[capabilitytypes.MemStoreKey],
	)

	scopedIBCKeeper := appKeepers.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedWasmKeeper := appKeepers.CapabilityKeeper.ScopeToModule(wasm.ModuleName)

	// add keepers
	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		appKeepers.keys[authtypes.StoreKey],
		appKeepers.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
	)

	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		appKeepers.keys[banktypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.GetSubspace(banktypes.ModuleName),
		blockedAddress,
	)

	// Cyber uses custom bank module wrapped around SDK's bank module
	appKeepers.CyberbankKeeper = cyberbankkeeper.NewIndexedKeeper(
		appCodec,
		appKeepers.keys[authtypes.StoreKey],
		cyberbankkeeper.Wrap(appKeepers.BankKeeper),
		appKeepers.AccountKeeper,
	)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[stakingtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		appKeepers.GetSubspace(stakingtypes.ModuleName),
	)

	appKeepers.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[minttypes.StoreKey],
		appKeepers.GetSubspace(minttypes.ModuleName),
		&stakingKeeper,
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		authtypes.FeeCollectorName,
	)

	appKeepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[distrtypes.StoreKey],
		appKeepers.GetSubspace(distrtypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		&stakingKeeper,
		authtypes.FeeCollectorName,
		blockedAddress,
	)

	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[slashingtypes.StoreKey],
		&stakingKeeper,
		appKeepers.GetSubspace(slashingtypes.ModuleName),
	)

	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appKeepers.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		appKeepers.CyberbankKeeper.Proxy,
		authtypes.FeeCollectorName,
	)

	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		appKeepers.keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		nil,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	appKeepers.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			appKeepers.DistrKeeper.Hooks(),
			appKeepers.SlashingKeeper.Hooks(),
		),
	)

	// Start cyber's keepers configuration

	appKeepers.CyberbankKeeper.Proxy.AddHook(
		bandwidth.CollectAddressesWithStakeChange(),
	)

	appKeepers.BandwidthMeter = bandwidthkeeper.NewBandwidthMeter(
		appCodec,
		appKeepers.keys[bandwidthtypes.StoreKey],
		appKeepers.CyberbankKeeper.Proxy,
		appKeepers.GetSubspace(bandwidthtypes.ModuleName),
	)

	appKeepers.GraphKeeper = graphkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[graphtypes.ModuleName],
		appKeepers.tkeys[paramstypes.TStoreKey],
	)

	appKeepers.IndexKeeper = graphkeeper.NewIndexKeeper(
		*appKeepers.GraphKeeper,
		appKeepers.tkeys[paramstypes.TStoreKey],
	)

	computeUnit := ranktypes.ComputeUnit(cast.ToInt(appOpts.Get(rank.FlagComputeGPU)))
	searchAPI := cast.ToBool(appOpts.Get(rank.FlagSearchAPI))
	appKeepers.RankKeeper = rankkeeper.NewKeeper(
		appKeepers.keys[ranktypes.ModuleName],
		appKeepers.GetSubspace(ranktypes.ModuleName),
		searchAPI,
		appKeepers.CyberbankKeeper,
		appKeepers.IndexKeeper,
		appKeepers.GraphKeeper,
		appKeepers.AccountKeeper,
		computeUnit,
	)

	appKeepers.GridKeeper = gridkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[gridtypes.ModuleName],
		appKeepers.CyberbankKeeper.Proxy,
		appKeepers.AccountKeeper,
		appKeepers.GetSubspace(gridtypes.ModuleName),
	)
	appKeepers.CyberbankKeeper.SetGridKeeper(&appKeepers.GridKeeper)
	appKeepers.CyberbankKeeper.SetAccountKeeper(appKeepers.AccountKeeper)

	appKeepers.ResourcesKeeper = resourceskeeper.NewKeeper(
		appCodec,
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		appKeepers.BandwidthMeter,
		appKeepers.GetSubspace(resourcestypes.ModuleName),
	)

	appKeepers.DmnKeeper = dmnkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[dmntypes.StoreKey],
		appKeepers.CyberbankKeeper.Proxy,
		appKeepers.AccountKeeper,
		appKeepers.GetSubspace(dmntypes.ModuleName),
	)

	appKeepers.LiquidityKeeper = liquiditykeeper.NewKeeper(
		appCodec,
		appKeepers.keys[liquiditytypes.StoreKey],
		appKeepers.GetSubspace(liquiditytypes.ModuleName),
		appKeepers.CyberbankKeeper.Proxy,
		appKeepers.AccountKeeper,
		appKeepers.DistrKeeper,
	)

	// End cyber's keepers configuration

	// Create IBC Keeper
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibchost.StoreKey],
		appKeepers.GetSubspace(ibchost.ModuleName),
		appKeepers.StakingKeeper,
		appKeepers.UpgradeKeeper,
		scopedIBCKeeper,
	)

	appKeepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[feegrant.StoreKey],
		appKeepers.AccountKeeper,
	)

	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(
		appKeepers.keys[authzkeeper.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
	)

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, sdkparams.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(appKeepers.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(appKeepers.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper))

	// Create Transfer Keepers
	appKeepers.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		scopedTransferKeeper,
	)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[evidencetypes.StoreKey],
		&appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	appKeepers.EvidenceKeeper = *evidenceKeeper

	// TODO update later to data (move wasm dir inside data directory)
	wasmDir := filepath.Join(homePath, "wasm")

	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}

	cyberOpts := wasmplugins.RegisterCustomPlugins(
		appKeepers.RankKeeper,
		appKeepers.GraphKeeper,
		appKeepers.DmnKeeper,
		appKeepers.GridKeeper,
		appKeepers.BandwidthMeter,
		appKeepers.LiquidityKeeper,
	)
	wasmOpts = append(wasmOpts, cyberOpts...)

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	availableCapabilities := strings.Join(AllCapabilities(), ",")
	appKeepers.WasmKeeper = wasm.NewKeeper(
		appCodec,
		appKeepers.keys[wasm.StoreKey],
		appKeepers.GetSubspace(wasm.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		appKeepers.StakingKeeper,
		appKeepers.DistrKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		appKeepers.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		availableCapabilities,
		wasmOpts...,
	)

	appKeepers.DmnKeeper.SetWasmKeeper(appKeepers.WasmKeeper)

	if len(enabledProposals) != 0 {
		govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(appKeepers.WasmKeeper, enabledProposals))
	}

	transferIBCModule := transfer.NewIBCModule(appKeepers.TransferKeeper)
	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.
		AddRoute(ibctransfertypes.ModuleName, transferIBCModule).
		AddRoute(wasm.ModuleName, wasm.NewIBCHandler(appKeepers.WasmKeeper, appKeepers.IBCKeeper.ChannelKeeper))
	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	appKeepers.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[govtypes.StoreKey],
		appKeepers.GetSubspace(govtypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		&stakingKeeper,
		govRouter,
	)

	appKeepers.ScopedIBCKeeper = scopedIBCKeeper
	appKeepers.ScopedTransferKeeper = scopedTransferKeeper
	appKeepers.ScopedWasmKeeper = scopedWasmKeeper

	return appKeepers
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

// GetSubspace returns a param subspace for a given module name.
func (appKeepers *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// AllCapabilities returns all capabilities available with the current wasmvm
// See https://github.com/CosmWasm/cosmwasm/blob/main/docs/CAPABILITIES-BUILT-IN.md
// This functionality is going to be moved upstream: https://github.com/CosmWasm/wasmvm/issues/425
func AllCapabilities() []string {
	return []string{
		"iterator",
		"staking",
		"stargate",
		"cyber",
	}
}
