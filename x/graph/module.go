package graph

import (
	"context"
	"encoding/json"

	//"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	abci "github.com/tendermint/tendermint/abci/types"

	bandwidthkeeper "github.com/cybercongress/go-cyber/x/bandwidth/keeper"

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

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct {
	cdc codec.Codec
}

func (AppModuleBasic) Name() string { return types.ModuleName }

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis(_ codec.JSONCodec) json.RawMessage { return nil }

func (AppModuleBasic) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, _ json.RawMessage) error {
	return nil
}

func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	rest.RegisterRoutes(clientCtx, rtr)
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
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

type AppModule struct {
	AppModuleBasic

	gk *keeper.GraphKeeper
	ik *keeper.IndexKeeper
	ak authkeeper.AccountKeeper
	bk *cyberbankkeeper.IndexedKeeper
	bm *bandwidthkeeper.BandwidthMeter
}

func NewAppModule(
	cdc codec.Codec,
	graphKeeper *keeper.GraphKeeper,
	indexKeeper *keeper.IndexKeeper,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper *cyberbankkeeper.IndexedKeeper,
	bandwidthKeeper *bandwidthkeeper.BandwidthMeter,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		gk:             graphKeeper,
		ik:             indexKeeper,
		ak:             accountKeeper,
		bk:             bankKeeper,
		bm:             bandwidthKeeper,
	}
}

func (AppModule) Name() string { return types.ModuleName }

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (am AppModule) Route() sdk.Route {
	return sdk.Route{}
}

func (am AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.gk, am.ik, am.ak, am.bk, am.bm))
	types.RegisterQueryServer(cfg.QueryServer(), am.gk)
}

func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return keeper.NewQuerier(*am.gk, legacyQuerierCdc)
}

func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, _ json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONCodec) json.RawMessage {
	err := keeper.WriteGenesis(ctx, *am.gk, am.ik)
	if err != nil {
		panic(err)
	}
	return nil
}

func (am AppModule) ConsensusVersion() uint64 {
	return 1
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.gk, am.ik)
	return []abci.ValidatorUpdate{}
}
