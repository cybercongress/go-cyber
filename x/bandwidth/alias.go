package bandwidth

import (
	"github.com/cybercongress/cyberd/x/bandwidth/exported"
	"github.com/cybercongress/cyberd/x/bandwidth/internal/keeper"
	"github.com/cybercongress/cyberd/x/bandwidth/internal/types"
)

const (
	DefaultParamspace = types.DefaultParamspace
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
)

type (
	Keeper                    = exported.Keeper
	BlockSpentBandwidthKeeper = exported.BlockSpentBandwidthKeeper

	Meter        = types.BandwidthMeter
	AcсBandwidth = types.AcсBandwidth
	Params       = types.Params
	GenesisState = types.GenesisState
)

var (
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	NewDefaultParams    = types.NewDefaultParams

	NewAccBandwidthKeeper        = keeper.NewAccBandwidthKeeper
	NewBlockSpentBandwidthKeeper = keeper.NewBlockSpentBandwidthKeeper
	ParamKeyTable                = keeper.ParamKeyTable
)
