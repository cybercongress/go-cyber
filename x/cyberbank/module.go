package cyberbank

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cybercongress/go-cyber/x/cyberbank/keeper"
	"github.com/cybercongress/go-cyber/x/cyberbank/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
)

type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {}

func (AppModuleBasic) DefaultGenesis() json.RawMessage { return nil }

func (AppModuleBasic) ValidateGenesis(json.RawMessage) error { return nil }

func (AppModuleBasic) RegisterRESTRoutes(_ context.CLIContext, _ *mux.Router) {}

func (AppModuleBasic) GetTxCmd(_ *codec.Codec) *cobra.Command { return nil }

func (AppModuleBasic) GetQueryCmd(_ *codec.Codec) *cobra.Command { return nil }

//____________________________________________________________________________

type AppModule struct {
	AppModuleBasic

	keeper        *keeper.IndexedKeeper
}


func NewAppModule(ik *keeper.IndexedKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         ik,
	}
}

func (AppModule) Name() string {
	return ModuleName
}

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (AppModule) Route() string {
	return ""
}

func (am AppModule) NewHandler() sdk.Handler {
	return nil
}

func (AppModule) QuerierRoute() string {
	return ""
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) { }

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) InitGenesis(ctx sdk.Context, _ json.RawMessage) []abci.ValidatorUpdate {
	InitGenesis(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(sdk.Context) json.RawMessage {
	return nil
}