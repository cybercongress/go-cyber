package core

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/tendermint/tendermint/rpc/client"
	"net/http"
	"time"
)

type ProxyContext struct {
	Codec      *wire.Codec
	Node       client.Client
	NodeUrl    string
	HttpClient *http.Client
}

func NewProxyContext(endpoint string) ProxyContext {
	return ProxyContext{
		Codec:   MakeCodec(),
		Node:    client.NewHTTP(endpoint, "/websocket"),
		NodeUrl: endpoint,
		HttpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}
