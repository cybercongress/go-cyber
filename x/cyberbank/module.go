package cyberbank

import (
	"encoding/json"
	//"context"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	//"github.com/cosmos/cosmos-sdk/client/context"
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

type AppModuleBasic struct{
	cdc codec.Marshaler
}

func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

func (AppModuleBasic) DefaultGenesis(_ codec.JSONMarshaler) json.RawMessage {
	return nil
}

func (AppModuleBasic) ValidateGenesis(_ codec.JSONMarshaler, _ client.TxEncodingConfig, _ json.RawMessage) error {
	return nil
}

func (AppModuleBasic) RegisterRESTRoutes(_ client.Context, _ *mux.Router) {}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}

func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
}

func (AppModuleBasic) RegisterInterfaces(_ codectypes.InterfaceRegistry) {}

//____________________________________________________________________________

type AppModule struct {
	AppModuleBasic
	keeper        *keeper.IndexedKeeper
}

func (am AppModule) RegisterServices(_ module.Configurator) {}

func NewAppModule(ik *keeper.IndexedKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         ik,
	}
}

func (AppModule) Name() string {
	return types.ModuleName
}

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (AppModule) Route() sdk.Route {
	return sdk.Route{}
}

func (am AppModule) NewHandler() sdk.Handler {
	return nil
}

func (AppModule) QuerierRoute() string {
	return ""
}

func (am AppModule) LegacyQuerierHandler(_ *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONMarshaler, _ json.RawMessage) []abci.ValidatorUpdate {
	keeper.InitGenesis(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(_ sdk.Context, _ codec.JSONMarshaler) json.RawMessage {
	return nil
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}