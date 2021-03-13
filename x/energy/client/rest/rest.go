package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

const (
	Source           = "src"
	Destination      = "dst"
)

// RegisterRoutes registers power-related REST handlers to a router
func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type CreateEnergyRouteReq struct {
	BaseReq rest.BaseReq 		`json:"base_req" yaml:"base_req"`
	//Source 		 string 		`json:"source" yaml:"source"`
	Destination  string 		`json:"destination" yaml:"destination"`
	Alias 		 string 		`json:"alias" yaml:"alias"`
}

type EditEnergyRouteReq struct {
	BaseReq rest.BaseReq 		`json:"base_req" yaml:"base_req"`
	//Source 		 string 		`json:"source" yaml:"source"`
	Destination  string 		`json:"destination" yaml:"destination"`
	Value 		 string 		`json:"value" yaml:"value"`
}

type DeleteEnergyRouteReq struct {
	BaseReq rest.BaseReq 		`json:"base_req" yaml:"base_req"`
	//Source 		 string 		`json:"source" yaml:"source"`
	Destination  string 		`json:"destination" yaml:"destination"`
}

type EditEnergyRouteAliasReq struct {
	BaseReq rest.BaseReq 		`json:"base_req" yaml:"base_req"`
	//Source 		 string		    `json:"source" yaml:"source"`
	Destination  string		    `json:"destination" yaml:"destination"`
	Alias 		 string 		`json:"alias" yaml:"alias"`
}
