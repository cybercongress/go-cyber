package core

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/tendermint/tendermint/rpc/client"
)

type ProxyContext struct {
	Codec *wire.Codec
	Node  client.Client
}

func NewProxyContext(endpoint string) ProxyContext {
	return ProxyContext{
		Codec: MakeCodec(),
		Node:  client.NewHTTP(endpoint, "/websocket"),
	}
}
