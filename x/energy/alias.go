package energy

import (
	"github.com/cybercongress/go-cyber/x/energy/exported"
	"github.com/cybercongress/go-cyber/x/energy/keeper"
	"github.com/cybercongress/go-cyber/x/energy/types"
)

const (
	ModuleName                = types.ModuleName
	StoreKey                = types.StoreKey
	DefaultParamspace 	      = types.DefaultParamspace
	QueryParams  			  = types.QueryParams
	EnergyPoolName 			  = types.EnergyPoolName
)

var (
	NewQuerier                = keeper.NewQuerier
	NewKeeper				  = keeper.NewKeeper
	DefaultGenesisState       = types.DefaultGenesisState
)

type (
	EnergyKeeper 		      = exported.EnergyKeeper
	Keeper                    = keeper.Keeper
	Route				  	  = types.Route
	Routes				  	  = types.Routes
	Params 					  = types.Params
	GenesisState			  = types.GenesisState
)