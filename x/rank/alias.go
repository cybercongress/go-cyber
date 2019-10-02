package rank

import (
	"github.com/cybercongress/cyberd/x/rank/exported"
	"github.com/cybercongress/cyberd/x/rank/internal/keeper"
	"github.com/cybercongress/cyberd/x/rank/internal/types"
)

const (
	// keeper
	DefaultParamspace = keeper.DefaultParamspace

	// types
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
	RouterKey  = types.RouterKey
	CPU        = types.CPU
	GPU        = types.GPU
)

type (
	// exported
	Keeper      = exported.Keeper
	StateKeeper = exported.StateKeeper

	// types
	Params       = types.Params
	GenesisState = types.GenesisState
	ComputeUnit  = types.ComputeUnit
)

var (
	// keeper
	NewStateKeeper = keeper.NewStateKeeper
	ParamKeyTable  = keeper.ParamKeyTable

	// types
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
	NewDefaultParams    = types.NewDefaultParams
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
)
