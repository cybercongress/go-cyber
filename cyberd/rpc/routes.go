package rpc

import (
	"github.com/cybercongress/cyberd/cosmos/poc/app"
	"github.com/tendermint/tendermint/rpc/core"
	"github.com/tendermint/tendermint/rpc/lib/server"
)

var cyberdApp *app.CyberdApp

func SetCyberdApp(cApp *app.CyberdApp) {
	cyberdApp = cApp
}

var Routes = map[string]*rpcserver.RPCFunc{
	"search":  rpcserver.NewRPCFunc(Search, "cid,page,perPage"),
	"account": rpcserver.NewRPCFunc(Account, "address"),
}

func init() {
	for path, rpcFunc := range Routes {
		core.Routes[path] = rpcFunc
	}
}
