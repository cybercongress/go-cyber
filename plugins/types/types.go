package types

import (
	bandwidthtypes "github.com/cybercongress/go-cyber/v6/x/bandwidth/types"
	dmntypes "github.com/cybercongress/go-cyber/v6/x/dmn/types"
	graphtypes "github.com/cybercongress/go-cyber/v6/x/graph/types"
	gridtypes "github.com/cybercongress/go-cyber/v6/x/grid/types"
	ranktypes "github.com/cybercongress/go-cyber/v6/x/rank/types"
	resourcestypes "github.com/cybercongress/go-cyber/v6/x/resources/types"
	tokenfactorytypes "github.com/cybercongress/go-cyber/v6/x/tokenfactory/wasm/types"
)

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

	// token factory queries
	//TokenFactory *types.TokenFactoryQuery `json:"token_factory,omitempty"`

	FullDenom       *tokenfactorytypes.FullDenom       `json:"full_denom,omitempty"`
	Admin           *tokenfactorytypes.DenomAdmin      `json:"admin,omitempty"`
	Metadata        *tokenfactorytypes.GetMetadata     `json:"metadata,omitempty"`
	DenomsByCreator *tokenfactorytypes.DenomsByCreator `json:"denoms_by_creator,omitempty"`
	Params          *tokenfactorytypes.GetParams       `json:"params,omitempty"`
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

	// token factory messages
	TokenFactory *tokenfactorytypes.TokenFactoryMsg `json:"token_factory,omitempty"`
}
