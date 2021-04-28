package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers power-related REST handlers to a router
func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type InvestmintReq struct {
	BaseReq  rest.BaseReq 	`json:"base_req" yaml:"base_req"`
	Amount   sdk.Coin 		`json:"amount" yaml:"amount"`
	Resource string 		`json:"resource" yaml:"resource"`
	Length   uint64 		`json:"length" yaml:"length"`
}