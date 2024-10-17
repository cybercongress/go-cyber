package clock

import (
	"context"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/cometbft/cometbft/abci/types"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cybercongress/go-cyber/v5/x/clock/client/cli"
	"github.com/cybercongress/go-cyber/v5/x/clock/keeper"
	"github.com/cybercongress/go-cyber/v5/x/clock/types"
)

const (
	ModuleName = types.ModuleName

	// ConsensusVersion defines the current x/clock module consensus version.
	ConsensusVersion = 1
)

var (
	_ module.AppModuleBasic   = AppModuleBasic{}
	_ module.AppModuleGenesis = AppModule{}
	_ module.AppModule        = AppModule{}
)

// AppModuleBasic defines the basic application module used by the wasm module.
type AppModuleBasic struct {
	cdc codec.Codec
}

func (a AppModuleBasic) Name() string {
	return types.ModuleName
}

func (a AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(&types.GenesisState{
		Params: types.DefaultParams(),
	})
}

func (a AppModuleBasic) ValidateGenesis(marshaler codec.JSONCodec, _ client.TxEncodingConfig, message json.RawMessage) error {
	var data types.GenesisState
	err := marshaler.UnmarshalJSON(message, &data)
	if err != nil {
		return err
	}
	if err := data.Params.Validate(); err != nil {
		return errorsmod.Wrap(err, "params")
	}
	return nil
}

func (a AppModuleBasic) RegisterRESTRoutes(_ client.Context, _ *mux.Router) {
}

func (a AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
	if err != nil {
		// same behavior as in cosmos-sdk
		panic(err)
	}
}

func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

func (a AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (a AppModuleBasic) RegisterInterfaces(r codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(r)
}

type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper
}

// NewAppModule constructor
func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
) *AppModule {
	return &AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
	}
}

func (a AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(a.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQuerier(a.keeper))
}

func (a AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {
}

func (a AppModule) InitGenesis(ctx sdk.Context, marshaler codec.JSONCodec, message json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	marshaler.MustUnmarshalJSON(message, &genesisState)
	_ = a.keeper.SetParams(ctx, genesisState.Params)
	return nil
}

func (a AppModule) ExportGenesis(ctx sdk.Context, marshaler codec.JSONCodec) json.RawMessage {
	genState := ExportGenesis(ctx, a.keeper)
	return marshaler.MustMarshalJSON(genState)
}

// ConsensusVersion is a sequence number for state-breaking change of the
// module. It should be incremented on each consensus-breaking change
// introduced by the module. To avoid wrong/empty versions, the initial version
// should be set to 1.
func (a AppModule) ConsensusVersion() uint64 {
	return ConsensusVersion
}

func (a AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlocker(ctx, a.keeper)
}

func (a AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, a.keeper)
	return nil
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}
