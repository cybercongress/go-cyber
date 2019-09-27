package bandwidth

import (
	"github.com/cybercongress/cyberd/x/bandwidth/types"
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
	RouterKey  = types.RouterKey
)

type (
	Keeper                    = types.Keeper
	BlockSpentBandwidthKeeper = types.BlockSpentBandwidthKeeper
	AcсBandwidth              = types.AcсBandwidth
	Params                    = types.Params
	GenesisState              = types.GenesisState
)

var (
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
)
