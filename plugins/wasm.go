package plugins

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bandwidthtypes "github.com/cybercongress/go-cyber/v4/x/bandwidth/types"
	dmntypes "github.com/cybercongress/go-cyber/v4/x/dmn/types"
	graphtypes "github.com/cybercongress/go-cyber/v4/x/graph/types"
	gridtypes "github.com/cybercongress/go-cyber/v4/x/grid/types"
	ranktypes "github.com/cybercongress/go-cyber/v4/x/rank/types"
	resourcestypes "github.com/cybercongress/go-cyber/v4/x/resources/types"
	tokenfactorykeeper "github.com/cybercongress/go-cyber/v4/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/cybercongress/go-cyber/v4/x/tokenfactory/types"

	"github.com/cybercongress/go-cyber/v4/plugins/types"
	cyberbankkeeper "github.com/cybercongress/go-cyber/v4/x/cyberbank/keeper"
	dmnwasm "github.com/cybercongress/go-cyber/v4/x/dmn/wasm"
	resourceskeeper "github.com/cybercongress/go-cyber/v4/x/resources/keeper"

	bandwidthkeeper "github.com/cybercongress/go-cyber/v4/x/bandwidth/keeper"
	bandwidthwasm "github.com/cybercongress/go-cyber/v4/x/bandwidth/wasm"
	dmnkeeper "github.com/cybercongress/go-cyber/v4/x/dmn/keeper"
	graphkeeper "github.com/cybercongress/go-cyber/v4/x/graph/keeper"
	graphwasm "github.com/cybercongress/go-cyber/v4/x/graph/wasm"
	gridkeeper "github.com/cybercongress/go-cyber/v4/x/grid/keeper"
	gridwasm "github.com/cybercongress/go-cyber/v4/x/grid/wasm"
	rankkeeper "github.com/cybercongress/go-cyber/v4/x/rank/keeper"
	rankwasm "github.com/cybercongress/go-cyber/v4/x/rank/wasm"
	resourceswasm "github.com/cybercongress/go-cyber/v4/x/resources/wasm"
	tokenfactorywasm "github.com/cybercongress/go-cyber/v4/x/tokenfactory/wasm"
)

func RegisterCustomPlugins(
	rank *rankkeeper.StateKeeper,
	graph *graphkeeper.GraphKeeper,
	dmn *dmnkeeper.Keeper,
	grid *gridkeeper.Keeper,
	bandwidth *bandwidthkeeper.BandwidthMeter,
	resources *resourceskeeper.Keeper,
	graphIndex *graphkeeper.IndexKeeper,
	account *authkeeper.AccountKeeper,
	cyberbank *cyberbankkeeper.IndexedKeeper,
	bank *bankkeeper.Keeper,
	tokenFactory *tokenfactorykeeper.Keeper,
) []wasmkeeper.Option {
	rankQuerier := rankwasm.NewWasmQuerier(rank)
	graphQuerier := graphwasm.NewWasmQuerier(graph)
	dmnQuerier := dmnwasm.NewWasmQuerier(dmn)
	gridQuerier := gridwasm.NewWasmQuerier(grid)
	bandwidthQuerier := bandwidthwasm.NewWasmQuerier(bandwidth)
	tokenFactoryQuerier := tokenfactorywasm.NewWasmQuerier(*bank, tokenFactory)

	graphMessenger := graphwasm.NewMessenger(graph, graphIndex, account, cyberbank, bandwidth)
	dmnMessenger := dmnwasm.NewMessenger(dmn)
	gridMessenger := gridwasm.NewMessenger(grid)
	resourcesMessenger := resourceswasm.NewMessenger(resources)
	tokenFactoryMessenger := tokenfactorywasm.NewMessenger(*bank, tokenFactory)

	moduleQueriers := map[string]types.ModuleQuerier{
		ranktypes.ModuleName:         rankQuerier,
		graphtypes.ModuleName:        graphQuerier,
		dmntypes.ModuleName:          dmnQuerier,
		gridtypes.ModuleName:         gridQuerier,
		bandwidthtypes.ModuleName:    bandwidthQuerier,
		tokenfactorytypes.ModuleName: tokenFactoryQuerier,
	}

	wasmQueryPlugin := types.NewQueryPlugin(
		moduleQueriers,
		rank,
		graph,
		dmn,
		grid,
		bandwidth,
		bank,
		tokenFactory,
	)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: types.CustomQuerier(wasmQueryPlugin),
	})

	moduleMessengers := map[string]types.ModuleMessenger{
		graphtypes.ModuleName:        graphMessenger,
		dmntypes.ModuleName:          dmnMessenger,
		gridtypes.ModuleName:         gridMessenger,
		resourcestypes.ModuleName:    resourcesMessenger,
		tokenfactorytypes.ModuleName: tokenFactoryMessenger,
	}

	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		types.CustomMessageDecorator(
			moduleMessengers,
			graph,
			dmn,
			grid,
			resources,
			bank,
			tokenFactory,
		),
	)

	return []wasm.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
