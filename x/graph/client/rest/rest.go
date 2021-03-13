package rest

import (
	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/cybercongress/go-cyber/x/graph/types"
)

func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type CreateCyberlinkReq struct {
	BaseReq        rest.BaseReq       `json:"base_req"`
	Links		   []types.Link 	  `json:"links"`
}