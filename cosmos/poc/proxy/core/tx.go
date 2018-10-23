package core

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/cyberd/cosmos/poc/http/util"
	"io/ioutil"
	"net/http"
)

func TxHandlerFn(ctx ProxyContext, unmarshal UnmarshalTxRequest) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		txRequestBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			util.HandleError(err, w)
			return
		}

		txReq, err := unmarshal(txRequestBytes)
		if err != nil {
			util.HandleError(err, w)
			return
		}

		// BUILDING COSMOS SDK TX
		signatures := make([]auth.StdSignature, 0, len(txReq.GetSignatures()))
		for _, sig := range txReq.GetSignatures() {
			stdSig := auth.StdSignature{
				PubKey: sig.PubKey, Signature: sig.Signature, AccountNumber: sig.AccountNumber, Sequence: sig.Sequence,
			}
			signatures = append(signatures, stdSig)
		}

		stdTx := auth.StdTx{Msgs: txReq.GetMsgs(), Fee: txReq.GetFee(), Signatures: signatures, Memo: txReq.GetMemo()}

		stdTxBytes, err := ctx.Codec.MarshalBinary(stdTx)
		if err != nil {
			util.HandleError(err, w)
			return
		}

		resp, err := ctx.Node.BroadcastTxSync(stdTxBytes)
		if err != nil {
			util.HandleError(err, w)
			return
		}

		respBytes, err := json.Marshal(resp)
		if err != nil {
			util.HandleError(err, w)
			return
		}
		w.Write(respBytes)
	}
}
