package link

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/abci/types"

	"github.com/cybercongress/go-cyber/x/bandwidth/exported"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) { RegisterCodec(cdc) }

func (AppModuleBasic) DefaultGenesis() json.RawMessage { return nil }

func (AppModuleBasic) ValidateGenesis(json.RawMessage) error { return nil }

func (AppModuleBasic) RegisterRESTRoutes(context.CLIContext, *mux.Router) {}

func (AppModuleBasic) GetTxCmd(*codec.Codec) *cobra.Command { return nil }

func (AppModuleBasic) GetQueryCmd(*codec.Codec) *cobra.Command { return nil }

type AppModule struct {
	AppModuleBasic

	cidNumberKeeper CidNumberKeeper
	indexedKeeper   IndexedKeeper
	accountKeeper   auth.AccountKeeper
	accountBandwidthKeeper exported.BaseAccountBandwidthKeeper
	meter exported.Meter
}

func NewAppModule(cidNumberKeeper CidNumberKeeper, indexedKeeper IndexedKeeper,
	accountKeeper auth.AccountKeeper, accountBandwidthKeeper exported.BaseAccountBandwidthKeeper, meter exported.Meter) AppModule {

	return AppModule{
		AppModuleBasic:  AppModuleBasic{},
		cidNumberKeeper: cidNumberKeeper,
		indexedKeeper:   indexedKeeper,
		accountKeeper:   accountKeeper,
		accountBandwidthKeeper: accountBandwidthKeeper,
		meter: meter,
	}
}

func (AppModule) Name() string {
	return ModuleName
}

func (am AppModule) InitGenesis(_ sdk.Context, _ json.RawMessage) []types.ValidatorUpdate { return nil }

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (am AppModule) Route() string { return RouterKey }

func (am AppModule) NewHandler() sdk.Handler {
	return NewLinksHandler(am.cidNumberKeeper, am.indexedKeeper, am.accountKeeper, am.accountBandwidthKeeper, am.meter)
}

func (am AppModule) QuerierRoute() string { return RouterKey }

func (am AppModule) NewQuerierHandler() sdk.Querier { return nil }

func (am AppModule) ExportGenesis(_ sdk.Context) json.RawMessage { return nil }

func (am AppModule) BeginBlock(sdk.Context, types.RequestBeginBlock) {}

func (am AppModule) EndBlock(sdk.Context, types.RequestEndBlock) []types.ValidatorUpdate { return nil }
