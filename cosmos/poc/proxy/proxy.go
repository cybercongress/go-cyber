package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cybercongress/cyberd/cosmos/poc/app"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/rpc/lib/types"
	"io/ioutil"
	"net/http"

	"github.com/rs/cors"
)

type Signature struct {
	PubKey        secp256k1.PubKeySecp256k1 `json:"pub_key"` // optional
	Signature     []byte           `json:"signature"`
	AccountNumber int64            `json:"account_number"`
	Sequence      int64            `json:"sequence"`
}

type LinkRequest struct {
	Fee        auth.StdFee         `json:"fee"`
	Msg        []app.MsgLink       `json:"msg"`
	Signatures []Signature `json:"signatures"`
	Memo       string              `json:"memo"`
}

func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	sdk.RegisterWire(cdc)
	bank.RegisterWire(cdc)
	ibc.RegisterWire(cdc)
	auth.RegisterWire(cdc)
	cdc.RegisterConcrete(app.MsgLink{}, "cyberd/Link", nil)

	cdc.Seal()
	return cdc
}

var codec = MakeCodec()

func handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	signatureBytes, _ := ioutil.ReadAll(r.Body)

	//addr, _ := sdk.AccAddressFromHex("cosmosaccaddr188uy9swfs5qvenvxjk209fwrswe3rhvhuyyu0y")
	//
	//linkEx := LinkRequest{Fee: auth.StdFee{}, Msg: []app.MsgLink{{addr, "2", "3"}}, Signatures: []auth.StdSignature{}, Memo: ""}
	//
	//jsonBytes, _ := json.Marshal(linkEx)
	//fmt.Println(string(jsonBytes))


	var request LinkRequest

	err := json.Unmarshal(signatureBytes, &request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(request)


	// BUILDING COSMOS SDK TX

	signatures := make([]auth.StdSignature, 0, len(request.Signatures))
	for _, sig := range request.Signatures {
		stdSig := auth.StdSignature{PubKey: sig.PubKey, Signature: sig.Signature, AccountNumber: sig.AccountNumber, Sequence: sig.Sequence}
		signatures = append(signatures, stdSig)
	}

	msgs := make([]sdk.Msg, 0, len(request.Msg))
	for _, msg := range request.Msg {
		msgs = append(msgs, msg)
	}

	stdTx := auth.StdTx{Msgs: msgs, Fee: request.Fee, Signatures: signatures, Memo: request.Memo}
	stdTxBytes, _ := codec.MarshalBinary(stdTx)


	// BUILDING RPC REQUEST

	params := map[string]interface{}{ "tx": stdTxBytes }
	rpcReq, err := rpctypes.MapToRequest(codec, "jsonrpc-client", "broadcast_tx_commit", params)
	rpcReqBytes, err := json.Marshal(rpcReq)

	requestBuf := bytes.NewBuffer(rpcReqBytes)

	resp, _ := http.Post("http://localhost:26657/", "text/json", requestBuf)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("post:\n", string(body))

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/link", handler)

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	handler := c.Handler(mux)
	http.ListenAndServe(":8081", handler)
}
