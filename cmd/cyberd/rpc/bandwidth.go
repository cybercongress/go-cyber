package rpc

import "github.com/tendermint/tendermint/rpc/lib/types"

type ResultBandwidthPrice struct {
	Price float64 `amino:"unsafe" json:"price"`
}

func CurrentBandwidthPrice(ctx *rpctypes.Context) (*ResultBandwidthPrice, error) {
	return &ResultBandwidthPrice{cyberdApp.CurrentBandwidthPrice()}, nil
}
