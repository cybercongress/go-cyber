package keepers

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/nft"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	tokenfactorykeeper "github.com/cybercongress/go-cyber/v4/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/cybercongress/go-cyber/v4/x/tokenfactory/types"
	"path/filepath"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
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
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	//"github.com/cybercongress/go-cyber/v4/app"

	"github.com/spf13/cast"

	liquiditykeeper "github.com/cybercongress/go-cyber/v4/x/liquidity/keeper"
	liquiditytypes "github.com/cybercongress/go-cyber/v4/x/liquidity/types"

	wasmplugins "github.com/cybercongress/go-cyber/v4/plugins"
	"github.com/cybercongress/go-cyber/v4/x/bandwidth"
	bandwidthkeeper "github.com/cybercongress/go-cyber/v4/x/bandwidth/keeper"
	bandwidthtypes "github.com/cybercongress/go-cyber/v4/x/bandwidth/types"
	cyberbankkeeper "github.com/cybercongress/go-cyber/v4/x/cyberbank/keeper"
	dmnkeeper "github.com/cybercongress/go-cyber/v4/x/dmn/keeper"
	dmntypes "github.com/cybercongress/go-cyber/v4/x/dmn/types"
	graphkeeper "github.com/cybercongress/go-cyber/v4/x/graph/keeper"
	graphtypes "github.com/cybercongress/go-cyber/v4/x/graph/types"
	gridkeeper "github.com/cybercongress/go-cyber/v4/x/grid/keeper"
	gridtypes "github.com/cybercongress/go-cyber/v4/x/grid/types"
	"github.com/cybercongress/go-cyber/v4/x/rank"
	rankkeeper "github.com/cybercongress/go-cyber/v4/x/rank/keeper"
	ranktypes "github.com/cybercongress/go-cyber/v4/x/rank/types"
	resourceskeeper "github.com/cybercongress/go-cyber/v4/x/resources/keeper"
	resourcestypes "github.com/cybercongress/go-cyber/v4/x/resources/types"

	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"

	govv1beta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

var (
	wasmCapabilities = "iterator,staking,stargate,cyber,cosmwasm_1_1,cosmwasm_1_2,cosmwasm_1_3"

	tokenFactoryCapabilities = []string{
		tokenfactorytypes.EnableBurnFrom,
		tokenfactorytypes.EnableForceTransfer,
		tokenfactorytypes.EnableSetMetadata,
	}
)

// module account permissions
var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:     nil,
	distrtypes.ModuleName:          nil,
	minttypes.ModuleName:           {authtypes.Minter},
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:            {authtypes.Burner},
	nft.ModuleName:                 nil,
	ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	ibcfeetypes.ModuleName:         nil,
	wasmtypes.ModuleName:           {authtypes.Burner},
	liquiditytypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
	gridtypes.GridPoolName:         nil,
	resourcestypes.ResourcesName:   {authtypes.Minter, authtypes.Burner},
	tokenfactorytypes.ModuleName:   {authtypes.Minter, authtypes.Burner},
}

type AppKeepers struct {
	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    *stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     *crisiskeeper.Keeper
	UpgradeKeeper    *upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCFeeKeeper     ibcfeekeeper.Keeper
	EvidenceKeeper   evidencekeeper.Keeper
	TransferKeeper   ibctransferkeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	NFTKeeper        nftkeeper.Keeper
	AuthzKeeper      authzkeeper.Keeper

	ConsensusParamsKeeper consensusparamkeeper.Keeper

	TokenFactoryKeeper tokenfactorykeeper.Keeper
	WasmKeeper         wasmkeeper.Keeper
	LiquidityKeeper    liquiditykeeper.Keeper
	BandwidthMeter     *bandwidthkeeper.BandwidthMeter
	CyberbankKeeper    *cyberbankkeeper.IndexedKeeper
	GraphKeeper        *graphkeeper.GraphKeeper
	IndexKeeper        *graphkeeper.IndexKeeper
	RankKeeper         *rankkeeper.StateKeeper
	GridKeeper         gridkeeper.Keeper
	DmnKeeper          *dmnkeeper.Keeper
	ResourcesKeeper    resourceskeeper.Keeper

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
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
) AppKeepers {
	appKeepers := AppKeepers{}

	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.GenerateKeys()
	keys := appKeepers.GetKVStoreKey()
	tkeys := appKeepers.GetTransientStoreKey()

	appKeepers.ParamsKeeper = initParamsKeeper(
		appCodec,
		cdc,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)

	govModAddress := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// set the BaseApp's parameter store
	appKeepers.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		keys[consensusparamtypes.StoreKey],
		govModAddress,
	)
	bApp.SetParamStore(&appKeepers.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		appKeepers.memKeys[capabilitytypes.MemStoreKey],
	)

	scopedIBCKeeper := appKeepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedTransferKeeper := appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedWasmKeeper := appKeepers.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)

	// add keepers
	Bech32Prefix := "bostrom"
	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		authtypes.ProtoBaseAccount,
		maccPerms,
		Bech32Prefix,
		govModAddress,
	)

	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		appKeepers.AccountKeeper,
		BlockedAddresses(),
		govModAddress,
	)

	// Cyber uses custom bank module wrapped around SDK's bank module
	appKeepers.CyberbankKeeper = cyberbankkeeper.NewIndexedKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		cyberbankkeeper.Wrap(appKeepers.BankKeeper),
		appKeepers.AccountKeeper,
	)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		govModAddress,
	)

	appKeepers.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		keys[minttypes.StoreKey],
		stakingKeeper,
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		authtypes.FeeCollectorName,
		govModAddress,
	)

	appKeepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		keys[distrtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		stakingKeeper,
		authtypes.FeeCollectorName,
		govModAddress,
	)

	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		cdc,
		keys[slashingtypes.StoreKey],
		stakingKeeper,
		govModAddress,
	)

	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))
	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		keys[crisistypes.StoreKey],
		invCheckPeriod,
		appKeepers.CyberbankKeeper.Proxy,
		authtypes.FeeCollectorName,
		govModAddress,
	)

	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		bApp,
		govModAddress,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			appKeepers.DistrKeeper.Hooks(),
			appKeepers.SlashingKeeper.Hooks(),
		),
	)
	appKeepers.StakingKeeper = stakingKeeper

	// Start cyber's keepers configuration

	appKeepers.CyberbankKeeper.Proxy.AddHook(
		bandwidth.CollectAddressesWithStakeChange(),
	)

	appKeepers.BandwidthMeter = bandwidthkeeper.NewBandwidthMeter(
		appCodec,
		keys[bandwidthtypes.StoreKey],
		tkeys[bandwidthtypes.TStoreKey],
		appKeepers.CyberbankKeeper.Proxy,
		govModAddress,
	)

	appKeepers.GraphKeeper = graphkeeper.NewKeeper(
		appCodec,
		keys[graphtypes.ModuleName],
		// TODO fixed appKeepers.tkeys[paramstypes.TStoreKey], need test
		tkeys[graphtypes.TStoreKey],
	)

	appKeepers.IndexKeeper = graphkeeper.NewIndexKeeper(
		*appKeepers.GraphKeeper,
		// TODO fixed appKeepers.tkeys[paramstypes.TStoreKey], need test
		tkeys[graphtypes.TStoreKey],
	)

	computeUnit := ranktypes.ComputeUnit(cast.ToInt(appOpts.Get(rank.FlagComputeGPU)))
	searchAPI := cast.ToBool(appOpts.Get(rank.FlagSearchAPI))
	appKeepers.RankKeeper = rankkeeper.NewKeeper(
		appCodec,
		keys[ranktypes.ModuleName],
		searchAPI,
		appKeepers.CyberbankKeeper,
		appKeepers.IndexKeeper,
		appKeepers.GraphKeeper,
		appKeepers.AccountKeeper,
		computeUnit,
		govModAddress,
	)

	appKeepers.GridKeeper = gridkeeper.NewKeeper(
		appCodec,
		keys[gridtypes.ModuleName],
		appKeepers.CyberbankKeeper.Proxy,
		appKeepers.AccountKeeper,
		govModAddress,
	)
	appKeepers.CyberbankKeeper.SetGridKeeper(&appKeepers.GridKeeper)
	appKeepers.CyberbankKeeper.SetAccountKeeper(appKeepers.AccountKeeper)

	appKeepers.ResourcesKeeper = resourceskeeper.NewKeeper(
		appCodec,
		keys[gridtypes.ModuleName],
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		appKeepers.BandwidthMeter,
		govModAddress,
	)

	appKeepers.DmnKeeper = dmnkeeper.NewKeeper(
		appCodec,
		keys[dmntypes.StoreKey],
		appKeepers.CyberbankKeeper.Proxy,
		appKeepers.AccountKeeper,
		govModAddress,
	)

	appKeepers.LiquidityKeeper = liquiditykeeper.NewKeeper(
		appCodec,
		keys[liquiditytypes.StoreKey],
		appKeepers.CyberbankKeeper.Proxy,
		appKeepers.AccountKeeper,
		appKeepers.DistrKeeper,
		govModAddress,
	)

	// End cyber's keepers configuration

	// Create IBC Keeper
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibcexported.StoreKey],
		appKeepers.GetSubspace(ibcexported.ModuleName),
		stakingKeeper,
		appKeepers.UpgradeKeeper,
		scopedIBCKeeper,
	)

	appKeepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		keys[feegrant.StoreKey],
		appKeepers.AccountKeeper,
	)

	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
	)

	// register the proposal types
	govRouter := govv1beta.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govv1beta.ProposalHandler).
		AddRoute(paramproposal.RouterKey, sdkparams.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(appKeepers.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper))

	govConfig := govtypes.DefaultConfig()
	govConfig.MaxMetadataLen = 10200

	govKeeper := govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		appKeepers.StakingKeeper,
		bApp.MsgServiceRouter(),
		govtypes.DefaultConfig(),
		govModAddress,
	)

	appKeepers.NFTKeeper = nftkeeper.NewKeeper(keys[nftkeeper.StoreKey], appCodec, appKeepers.AccountKeeper, appKeepers.CyberbankKeeper.Proxy)

	appKeepers.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register governance hooks
		),
	)

	appKeepers.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		appCodec,
		keys[ibcfeetypes.StoreKey],
		appKeepers.IBCKeeper.ChannelKeeper, // may be replaced with IBC middleware
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
	)

	// Create Transfer Keepers
	appKeepers.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.IBCFeeKeeper, // ISC4 Wrapper: fee IBC middleware
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		scopedTransferKeeper,
	)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		keys[evidencetypes.StoreKey],
		stakingKeeper,
		appKeepers.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	appKeepers.EvidenceKeeper = *evidenceKeeper

	appKeepers.TokenFactoryKeeper = tokenfactorykeeper.NewKeeper(
		appCodec,
		appKeepers.keys[tokenfactorytypes.StoreKey],
		maccPerms,
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		appKeepers.DistrKeeper,
		tokenFactoryCapabilities,
		govModAddress,
	)

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
		&appKeepers.GridKeeper,
		appKeepers.BandwidthMeter,
		&appKeepers.ResourcesKeeper,
		appKeepers.IndexKeeper,
		&appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper,
		&appKeepers.BankKeeper,
		&appKeepers.TokenFactoryKeeper,
	)
	wasmOpts = append(wasmOpts, cyberOpts...)

	appKeepers.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		keys[wasmtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.CyberbankKeeper.Proxy,
		stakingKeeper,
		distrkeeper.NewQuerier(appKeepers.DistrKeeper),
		appKeepers.IBCFeeKeeper, // ISC4 Wrapper: fee IBC middleware
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		appKeepers.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		wasmCapabilities,
		govModAddress,
		wasmOpts...,
	)

	appKeepers.DmnKeeper.SetWasmKeeper(appKeepers.WasmKeeper)

	// register wasm gov proposal types
	// The gov proposal types can be individually enabled
	govRouter.AddRoute(wasm.RouterKey, wasmkeeper.NewLegacyWasmProposalHandler(appKeepers.WasmKeeper, wasmtypes.EnableAllProposals))
	// Set legacy router for backwards compatibility with gov v1beta1
	appKeepers.GovKeeper.SetLegacyRouter(govRouter)

	// Create Transfer Stack
	var transferStack porttypes.IBCModule
	transferStack = transfer.NewIBCModule(appKeepers.TransferKeeper)
	transferStack = ibcfee.NewIBCMiddleware(transferStack, appKeepers.IBCFeeKeeper)

	// Create fee enabled wasm ibc Stack
	var wasmStack porttypes.IBCModule
	wasmStack = wasm.NewIBCHandler(appKeepers.WasmKeeper, appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCFeeKeeper)
	wasmStack = ibcfee.NewIBCMiddleware(wasmStack, appKeepers.IBCFeeKeeper)

	ibcRouter := porttypes.NewRouter().
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(wasmtypes.ModuleName, wasmStack)
	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	appKeepers.ScopedIBCKeeper = scopedIBCKeeper
	appKeepers.ScopedTransferKeeper = scopedTransferKeeper
	appKeepers.ScopedWasmKeeper = scopedWasmKeeper

	return appKeepers
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(wasmtypes.ModuleName)
	paramsKeeper.Subspace(bandwidthtypes.ModuleName)
	paramsKeeper.Subspace(ranktypes.ModuleName)
	paramsKeeper.Subspace(gridtypes.ModuleName)
	paramsKeeper.Subspace(dmntypes.ModuleName)
	paramsKeeper.Subspace(resourcestypes.ModuleName)
	paramsKeeper.Subspace(liquiditytypes.ModuleName)
	paramsKeeper.Subspace(tokenfactorytypes.ModuleName)

	return paramsKeeper
}

// GetSubspace returns a param subspace for a given module name.
func (appKeepers *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// BlockedAddresses returns all the app's blocked account addresses.
func BlockedAddresses() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range GetMaccPerms() {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	// allow the following addresses to receive funds
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	return modAccAddrs
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}
