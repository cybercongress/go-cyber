package rank

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	//"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cybercongress/go-cyber/x/rank/client/cli"
	"github.com/cybercongress/go-cyber/x/rank/client/rest"

	"github.com/cybercongress/go-cyber/x/rank/keeper"
	"github.com/cybercongress/go-cyber/x/rank/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

const (
	FlagComputeGPU = "compute-gpu"
	FlagSearchAPI  = "search-api"
)

type AppModuleBasic struct {
	cdc codec.Codec
}

func (AppModuleBasic) Name() string { return types.ModuleName }

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return types.ValidateGenesis(&data)
}

func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	rest.RegisterRoutes(clientCtx, rtr)
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (AppModuleBasic) GetTxCmd() *cobra.Command { return nil }

func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

func (AppModuleBasic) RegisterInterfaces(_ codectypes.InterfaceRegistry) {}

type AppModule struct {
	AppModuleBasic

	rk *keeper.StateKeeper
}

func NewAppModule(
	cdc codec.Codec, rankKeeper *keeper.StateKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		rk:             rankKeeper,
	}
}

func AddModuleInitFlags(startCmd *cobra.Command) {
	startCmd.Flags().Bool(FlagComputeGPU, true, "Compute on GPU")
	startCmd.Flags().Bool(FlagSearchAPI, false, "Run search API")
}

func (AppModule) Name() string { return types.ModuleName }

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (am AppModule) Route() sdk.Route {
	return sdk.Route{}
}

func (am AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return keeper.NewQuerier(am.rk, legacyQuerierCdc)
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterQueryServer(cfg.QueryServer(), am.rk)
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	keeper.InitGenesis(ctx, *am.rk, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := keeper.ExportGenesis(ctx, *am.rk)
	return cdc.MustMarshalJSON(gs)
}

func (am AppModule) ConsensusVersion() uint64 {
	return 1
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.rk)
	return []abci.ValidatorUpdate{}
}
