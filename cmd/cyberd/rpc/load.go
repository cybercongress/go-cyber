package rpc

import "github.com/tendermint/tendermint/rpc/lib/types"

type ResultNetworkLoad struct {
	Load float64 `amino:"unsafe" json:"load"`
}

func CurrentNetworkLoad(ctx *rpctypes.Context) (*ResultNetworkLoad, error) {
	return &ResultNetworkLoad{cyberdApp.CurrentNetworkLoad()}, nil
}
