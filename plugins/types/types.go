package types

import (
	"errors"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bandwidthkeeper "github.com/cybercongress/go-cyber/v4/x/bandwidth/keeper"
	bandwidthtypes "github.com/cybercongress/go-cyber/v4/x/bandwidth/types"
	dmnkeeper "github.com/cybercongress/go-cyber/v4/x/dmn/keeper"
	dmntypes "github.com/cybercongress/go-cyber/v4/x/dmn/types"
	graphkeeper "github.com/cybercongress/go-cyber/v4/x/graph/keeper"
	graphtypes "github.com/cybercongress/go-cyber/v4/x/graph/types"
	gridkeeper "github.com/cybercongress/go-cyber/v4/x/grid/keeper"
	gridtypes "github.com/cybercongress/go-cyber/v4/x/grid/types"
	rankkeeper "github.com/cybercongress/go-cyber/v4/x/rank/keeper"
	ranktypes "github.com/cybercongress/go-cyber/v4/x/rank/types"
	resourceskeeper "github.com/cybercongress/go-cyber/v4/x/resources/keeper"
	resourcestypes "github.com/cybercongress/go-cyber/v4/x/resources/types"
)

type ModuleQuerier interface {
	HandleQuery(ctx sdk.Context, query CyberQuery) ([]byte, error)
}

type ModuleMessenger interface {
	HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg CyberMsg) ([]sdk.Event, [][]byte, error)
}

var ErrHandleQuery = errors.New("error handle query")

var ErrHandleMsg = errors.New("error handle message")

type QueryPlugin struct {
	moduleQueriers []ModuleQuerier
	rankKeeper     *rankkeeper.StateKeeper
	graphKeeper    *graphkeeper.GraphKeeper
	dmnKeeper      *dmnkeeper.Keeper
	gridKeeper     *gridkeeper.Keeper
	bandwidthMeter *bandwidthkeeper.BandwidthMeter
}

func NewQueryPlugin(
	moduleQueriers []ModuleQuerier,
	rank *rankkeeper.StateKeeper,
	graph *graphkeeper.GraphKeeper,
	dmn *dmnkeeper.Keeper,
	grid *gridkeeper.Keeper,
	bandwidth *bandwidthkeeper.BandwidthMeter,
) *QueryPlugin {
	return &QueryPlugin{
		moduleQueriers: moduleQueriers,
		rankKeeper:     rank,
		graphKeeper:    graph,
		dmnKeeper:      dmn,
		gridKeeper:     grid,
		bandwidthMeter: bandwidth,
	}
}

type CyberQuery struct {
	// rankKeeper queries
	ParticleRank *ranktypes.QueryRankRequest `json:"particle_rank,omitempty"`
	// TODO add IsLinkExist, IsAnyLinkExist

	// graph queries
	GraphStats *graphtypes.QueryGraphStatsRequest `json:"graph_stats,omitempty"`

	// dmn queries
	Thought      *dmntypes.QueryThoughtParamsRequest `json:"thought,omitempty"`
	ThoughtStats *dmntypes.QueryThoughtParamsRequest `json:"thought_stats,omitempty"`
	ThoughtsFees *dmntypes.QueryThoughtsFeesRequest  `json:"thoughts_fees,omitempty"`

	// grid queries
	SourceRoutes            *gridtypes.QuerySourceRequest      `json:"source_routes,omitempty"`
	SourceRoutedEnergy      *gridtypes.QuerySourceRequest      `json:"source_routed_energy,omitempty"`
	DestinationRoutedEnergy *gridtypes.QueryDestinationRequest `json:"destination_routed_energy,omitempty"`
	Route                   *gridtypes.QueryRouteRequest       `json:"route,omitempty"`

	// bandwidth queries
	BandwidthLoad   *bandwidthtypes.QueryLoadRequest            `json:"bandwidth_load,omitempty"`
	BandwidthPrice  *bandwidthtypes.QueryPriceRequest           `json:"bandwidth_price,omitempty"`
	TotalBandwidth  *bandwidthtypes.QueryTotalBandwidthRequest  `json:"total_bandwidth,omitempty"`
	NeuronBandwidth *bandwidthtypes.QueryNeuronBandwidthRequest `json:"neuron_bandwidth,omitempty"`
}

type CustomMessenger struct {
	wrapped          wasmkeeper.Messenger
	moduleMessengers []ModuleMessenger
	graphKeeper      *graphkeeper.GraphKeeper
	dmnKeeper        *dmnkeeper.Keeper
	gridKeeper       *gridkeeper.Keeper
	resourcesKeeper  *resourceskeeper.Keeper
}

func CustomMessageDecorator(
	moduleMessengers []ModuleMessenger,
	graph *graphkeeper.GraphKeeper,
	dmn *dmnkeeper.Keeper,
	grid *gridkeeper.Keeper,
	resources *resourceskeeper.Keeper,
) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:          old,
			moduleMessengers: moduleMessengers,
			graphKeeper:      graph,
			dmnKeeper:        dmn,
			gridKeeper:       grid,
			resourcesKeeper:  resources,
		}
	}
}

type CyberMsg struct {
	// graph messages
	Cyberlink *graphtypes.MsgCyberlink `json:"cyberlink,omitempty"`

	// dmn messages
	CreateThought         *dmntypes.MsgCreateThought         `json:"create_thought,omitempty"`
	ForgetThought         *dmntypes.MsgForgetThought         `json:"forget_thought,omitempty"`
	ChangeThoughtInput    *dmntypes.MsgChangeThoughtInput    `json:"change_thought_input,omitempty"`
	ChangeThoughtPeriod   *dmntypes.MsgChangeThoughtPeriod   `json:"change_thought_period,omitempty"`
	ChangeThoughtBlock    *dmntypes.MsgChangeThoughtBlock    `json:"change_thought_block,omitempty"`
	ChangeThoughtGasPrice *dmntypes.MsgChangeThoughtGasPrice `json:"change_thought_gas_price,omitempty"`
	ChangeThoughtParticle *dmntypes.MsgChangeThoughtParticle `json:"change_thought_particle,omitempty"`
	ChangeThoughtName     *dmntypes.MsgChangeThoughtName     `json:"change_thought_name,omitempty"`

	// grid messages
	CreateEnergyRoute   *gridtypes.MsgCreateRoute   `json:"create_energy_route,omitempty"`
	EditEnergyRoute     *gridtypes.MsgEditRoute     `json:"edit_energy_route,omitempty"`
	EditEnergyRouteName *gridtypes.MsgEditRouteName `json:"edit_energy_route_name,omitempty"`
	DeleteEnergyRoute   *gridtypes.MsgDeleteRoute   `json:"delete_energy_route,omitempty"`

	// resources messages
	Investmint *resourcestypes.MsgInvestmint `json:"investmint,omitempty"`
}
