package rest

import (
	"github.com/gorilla/mux"
	"github.com/cosmos/cosmos-sdk/client"
)

func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
}
