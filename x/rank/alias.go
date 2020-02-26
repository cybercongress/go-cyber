package rank

import (
	"github.com/cybercongress/go-cyber/x/rank/internal/keeper"
	"github.com/cybercongress/go-cyber/x/rank/internal/types"
)

const (
	ModuleName 			   = types.ModuleName
	DefaultParamspace 	   = types.DefaultParamspace
	StoreKey   			   = types.StoreKey
	QuerierRoute           = types.QuerierRoute
	QueryParameters        = types.QueryParameters
	QueryCalculationWindow = types.QueryCalculationWindow
	QueryDampingFactor     = types.QueryDampingFactor
	QueryTolerance         = types.QueryTolerance
	CPU        			   = types.CPU
	GPU        			   = types.GPU
)

var (
	// keeper
	NewStateKeeper 		= keeper.NewStateKeeper
	NewQuerier          = keeper.NewQuerier
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	ParamKeyTable  		= types.ParamKeyTable
	NewParams	        = types.NewParams
	DefaultParams       = types.DefaultParams

	ModuleCdc           = types.ModuleCdc
)

type (
	StateKeeper  = keeper.StateKeeper

	GenesisState = types.GenesisState
	Params       = types.Params
	ComputeUnit  = types.ComputeUnit
)
