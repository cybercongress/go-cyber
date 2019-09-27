package rank

import (
	"github.com/cybercongress/cyberd/x/rank/types"
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
	RouterKey  = types.RouterKey
)

type (
	Params       = types.Params
	Keeper       = types.Keeper
	GenesisState = types.GenesisState
)

var (
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
)
