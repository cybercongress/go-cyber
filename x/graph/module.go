package graph

import (
	"encoding/json"
	//"fmt"

	"context"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	abci "github.com/tendermint/tendermint/abci/types"
	//"github.com/gogo/protobuf/codec"
	"github.com/cosmos/cosmos-sdk/codec"
	//"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	//"github.com/tendermint/tendermint/abci/types"

	cyberbankkeeper "github.com/cybercongress/go-cyber/x/cyberbank/keeper"
	"github.com/cybercongress/go-cyber/x/graph/client/cli"
	"github.com/cybercongress/go-cyber/x/graph/client/rest"
	"github.com/cybercongress/go-cyber/x/graph/keeper"
	"github.com/cybercongress/go-cyber/x/graph/types"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
)

type AppModuleBasic struct{
	cdc codec.Marshaler
}

func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the reports module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the reports module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return nil
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, _ client.TxEncodingConfig, bz json.RawMessage) error {
	return nil
}

// RegisterRESTRoutes registers the REST routes for the reports module.
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	rest.RegisterRoutes(clientCtx, rtr)
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	_ = types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}

func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

//____________________________________________________________________________

type AppModule struct {
	AppModuleBasic
	gk keeper.GraphKeeper
	ik *keeper.IndexKeeper
	ak authkeeper.AccountKeeper
	bk *cyberbankkeeper.IndexedKeeper
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.gk, am.ik, am.ak, am.bk))
	types.RegisterQueryServer(cfg.QueryServer(), am.gk)
}

func NewAppModule(
	cdc 			codec.Marshaler,
	graphKeeper     keeper.GraphKeeper,
	indexKeeper 	*keeper.IndexKeeper,
	accountKeeper 	authkeeper.AccountKeeper,
	bankKeeper      *cyberbankkeeper.IndexedKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc}, // TODO remove cdc?
		gk:             graphKeeper,
		ik:             indexKeeper,
		ak:             accountKeeper,
		bk:				bankKeeper,
	}
}

func (AppModule) Name() string {
	return types.ModuleName
}

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(am.gk, am.ik, am.ak, am.bk))
}

func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.gk, am.ik, am.ak, am.bk)
}

func (am AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return keeper.NewQuerier(am.gk, legacyQuerierCdc)
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	_ = keeper.InitGenesis(ctx, am.gk, am.ik)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	_ = keeper.WriteGenesis(ctx, am.gk, am.ik)
	return nil
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {
}

func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

