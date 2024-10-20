package cyberbank

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cybercongress/go-cyber/v5/x/cyberbank/keeper"
	"github.com/cybercongress/go-cyber/v5/x/cyberbank/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct {
	cdc codec.Codec
}

func (AppModuleBasic) Name() string { return types.ModuleName }

func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

func (AppModuleBasic) DefaultGenesis(_ codec.JSONCodec) json.RawMessage {
	return nil
}

func (AppModuleBasic) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, _ json.RawMessage) error {
	return nil
}

func (AppModuleBasic) RegisterRESTRoutes(_ client.Context, _ *mux.Router) {}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}

func (AppModuleBasic) GetTxCmd() *cobra.Command { return nil }

func (AppModuleBasic) GetQueryCmd() *cobra.Command { return nil }

func (AppModuleBasic) RegisterInterfaces(_ codectypes.InterfaceRegistry) {}

type AppModule struct {
	AppModuleBasic

	keeper *keeper.IndexedKeeper
}

func NewAppModule(
	cdc codec.Codec,
	ik *keeper.IndexedKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         ik,
	}
}

func (AppModule) Name() string { return types.ModuleName }

func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (AppModule) RegisterServices(_ module.Configurator) {}

func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, _ json.RawMessage) []abci.ValidatorUpdate {
	// this will be called by the app module manager because of null module params in genesis, used direct call in InitChainer
	// am.keeper.InitGenesis(ctx)
	return []abci.ValidatorUpdate{}
}

func (AppModule) ExportGenesis(_ sdk.Context, _ codec.JSONCodec) json.RawMessage { return nil }

func (am AppModule) ConsensusVersion() uint64 {
	return 1
}

func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}
