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

	"github.com/cybercongress/go-cyber/x/link/client/cli"
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

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

type AppModule struct {
	AppModuleBasic

	graphKeeper     GraphKeeper
	indexKeeper   	*IndexKeeper
	accountKeeper   auth.AccountKeeper
}

func NewAppModule(
	graphKeeper     GraphKeeper,
	indexKeeper 	*IndexKeeper,
	accountKeeper 	auth.AccountKeeper,
	) AppModule {

	return AppModule{
		AppModuleBasic:  AppModuleBasic{},
		graphKeeper: 			graphKeeper,
		indexKeeper:   			indexKeeper,
		accountKeeper:   		accountKeeper,
	}
}

func (AppModule) Name() string {
	return ModuleName
}


func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (am AppModule) Route() string { return RouterKey }

func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.graphKeeper, am.indexKeeper, am.accountKeeper)
}

func (am AppModule) QuerierRoute() string { return QuerierRoute }

func (am AppModule) NewQuerierHandler() sdk.Querier { return nil }

func (am AppModule) InitGenesis(ctx sdk.Context, _ json.RawMessage) []types.ValidatorUpdate {
	InitGenesis(ctx, am.graphKeeper, am.indexKeeper)
	return []types.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	WriteGenesis(ctx, am.graphKeeper, am.indexKeeper)
	return nil
}

func (am AppModule) BeginBlock(sdk.Context, types.RequestBeginBlock) {}

func (am AppModule) EndBlock(sdk.Context, types.RequestEndBlock) []types.ValidatorUpdate {
	return []types.ValidatorUpdate{}
}
