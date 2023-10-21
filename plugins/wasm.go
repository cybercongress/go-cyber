package plugins

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	liquiditykeeper "github.com/gravity-devs/liquidity/x/liquidity/keeper"

	bandwidthkeeper "github.com/cybercongress/go-cyber/v2/x/bandwidth/keeper"
	bandwidthwasm "github.com/cybercongress/go-cyber/v2/x/bandwidth/wasm"
	dmnkeeper "github.com/cybercongress/go-cyber/v2/x/dmn/keeper"
	dmnwasm "github.com/cybercongress/go-cyber/v2/x/dmn/wasm"
	graphkeeper "github.com/cybercongress/go-cyber/v2/x/graph/keeper"
	graphwasm "github.com/cybercongress/go-cyber/v2/x/graph/wasm"
	gridkeeper "github.com/cybercongress/go-cyber/v2/x/grid/keeper"
	gridwasm "github.com/cybercongress/go-cyber/v2/x/grid/wasm"
	rankkeeper "github.com/cybercongress/go-cyber/v2/x/rank/keeper"
	rankwasm "github.com/cybercongress/go-cyber/v2/x/rank/wasm"
	resourceswasm "github.com/cybercongress/go-cyber/v2/x/resources/wasm"
)

func RegisterCustomPlugins(
	rank *rankkeeper.StateKeeper,
	graph *graphkeeper.GraphKeeper,
	dmn *dmnkeeper.Keeper,
	grid gridkeeper.Keeper,
	bandwidth *bandwidthkeeper.BandwidthMeter,
	liquidity liquiditykeeper.Keeper,
) []wasmkeeper.Option {
	// Initialize Cyber's query integrations
	querier := NewQuerier()
	queries := map[string]WasmQuerierInterface{
		WasmQueryRouteRank:      rankwasm.NewWasmQuerier(rank),
		WasmQueryRouteGraph:     graphwasm.NewWasmQuerier(*graph),
		WasmQueryRouteDmn:       dmnwasm.NewWasmQuerier(*dmn),
		WasmQueryRouteGrid:      gridwasm.NewWasmQuerier(grid),
		WasmQueryRouteBandwidth: bandwidthwasm.NewWasmQuerier(bandwidth),
		WasmQueryRouteLiquidity: NewWasmQuerier(liquidity),
	}
	querier.Queriers = queries

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasm.QueryPlugins{
		Custom: querier.QueryCustom,
	})

	// Initialize Cyber's encoder integrations
	parser := NewMsgParser()
	parsers := map[string]WasmMsgParserInterface{
		WasmMsgParserRouteGraph:     graphwasm.NewWasmMsgParser(),
		WasmMsgParserRouteDmn:       dmnwasm.NewWasmMsgParser(),
		WasmMsgParserRouteGrid:      gridwasm.NewWasmMsgParser(),
		WasmMsgParserRouteResources: resourceswasm.NewWasmMsgParser(),
		WasmMsgParserLiquidity:      NewWasmMsgParser(),
	}
	parser.Parsers = parsers

	messengerEncodersOpt := wasmkeeper.WithMessageEncoders(&wasm.MessageEncoders{
		Custom: parser.ParseCustom,
	})

	return []wasm.Option{
		queryPluginOpt,
		messengerEncodersOpt,
	}
}
