package proxy

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/tendermint/tendermint/rpc/client"
	"io/ioutil"
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

func (ctx ProxyContext) Get(endpoint string) (response []byte, err error) {

	resp, err := ctx.HttpClient.Get(ctx.NodeUrl + endpoint)
	if err != nil {
		return
	}

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}
