package rpc

import (
	"github.com/cybercongress/go-cyber/app"
	"github.com/tendermint/tendermint/rpc/core"
	rpcserver "github.com/tendermint/tendermint/rpc/jsonrpc/server"
)

var cyberdApp *app.CyberdApp
var codec = app.MakeCodec()

func SetCyberdApp(cApp *app.CyberdApp) {
	cyberdApp = cApp
}

var Routes = map[string]*rpcserver.RPCFunc{
	"search":                  rpcserver.NewRPCFunc(Search, "cid,page,perPage"),
	"backlinks":               rpcserver.NewRPCFunc(Backlinks, "cid,page,perPage"),
	"top":                     rpcserver.NewRPCFunc(Top, "page,perPage"),
	"rank":                    rpcserver.NewRPCFunc(Rank, "cid,proof"),
	"account":                 rpcserver.NewRPCFunc(Account, "address"),
	"account_bandwidth":       rpcserver.NewRPCFunc(AccountBandwidth, "address"),
	"account_links":		   rpcserver.NewRPCFunc(AccountLinks, "address,page,perPage"),
	"is_link_exist":           rpcserver.NewRPCFunc(IsLinkExist, "from,to,address"),
	"current_bandwidth_price": rpcserver.NewRPCFunc(CurrentBandwidthPrice, ""),
	"current_network_load":    rpcserver.NewRPCFunc(CurrentNetworkLoad, ""),
	"index_stats":             rpcserver.NewRPCFunc(IndexStats, ""),
}

func init() {
	for path, rpcFunc := range Routes {
		core.Routes[path] = rpcFunc
	}
}
