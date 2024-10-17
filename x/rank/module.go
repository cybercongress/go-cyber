package rank

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cobra"

	"github.com/cybercongress/go-cyber/v5/x/rank/client/cli"
	"github.com/cybercongress/go-cyber/v5/x/rank/exported"
	"github.com/cybercongress/go-cyber/v5/x/rank/keeper"
	"github.com/cybercongress/go-cyber/v5/x/rank/types"
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

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

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

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (AppModuleBasic) GetTxCmd() *cobra.Command { return nil }

func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

type AppModule struct {
	AppModuleBasic
	cdc codec.Codec

	rk             *keeper.StateKeeper
	legacySubspace exported.Subspace
}

func NewAppModule(
	cdc codec.Codec,
	rankKeeper *keeper.StateKeeper,
	ss exported.Subspace,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		cdc:            cdc,
		rk:             rankKeeper,
		legacySubspace: ss,
	}
}

func AddModuleInitFlags(startCmd *cobra.Command) {
	startCmd.Flags().Bool(FlagComputeGPU, true, "Compute on GPU")
	startCmd.Flags().Bool(FlagSearchAPI, false, "Run search API")
}

func (AppModule) Name() string { return types.ModuleName }

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterQueryServer(cfg.QueryServer(), am.rk)
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(*am.rk))

	m := keeper.NewMigrator(*am.rk, am.legacySubspace)
	if err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 1 to 2: %v", types.ModuleName, err))
	}
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
	return 2
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.rk)
	return []abci.ValidatorUpdate{}
}
