package rpc

import (
	"github.com/cybercongress/cyberd/app"
	"github.com/tendermint/tendermint/rpc/core"
	"github.com/tendermint/tendermint/rpc/lib/server"
)

var cyberdApp *app.CyberdApp

func SetCyberdApp(cApp *app.CyberdApp) {
	cyberdApp = cApp
}

var Routes = map[string]*rpcserver.RPCFunc{
	"search":                  rpcserver.NewRPCFunc(Search, "cid,page,perPage"),
	"account":                 rpcserver.NewRPCFunc(Account, "address"),
	"account_bandwidth":       rpcserver.NewRPCFunc(AccountBandwidth, "address"),
	"is_link_exist":           rpcserver.NewRPCFunc(IsLinkExist, "from,to,address"),
	"current_bandwidth_price": rpcserver.NewRPCFunc(CurrentBandwidthPrice, ""),
	"index_stats":             rpcserver.NewRPCFunc(IndexStats, ""),
}

func init() {
	for path, rpcFunc := range Routes {
		core.Routes[path] = rpcFunc
	}
}
