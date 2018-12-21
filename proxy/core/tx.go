package core

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/cyberd/util"
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
			stdSig := auth.StdSignature{PubKey: sig.PubKey, Signature: sig.Signature}
			signatures = append(signatures, stdSig)
		}

		stdTx := auth.StdTx{Msgs: txReq.GetMsgs(), Fee: txReq.GetFee(), Signatures: signatures, Memo: txReq.GetMemo()}

		stdTxBytes, err := ctx.Codec.MarshalBinaryLengthPrefixed(stdTx)
		if err != nil {
			util.HandleError(err, w)
			return
		}

		commit, err := util.GetBooleanParamValue("commit", false, r)
		if err != nil {
			util.HandleError(err, w)
			return
		}

		var result []byte

		if commit {
			result, err = ctx.BroadcastTxCommit(stdTxBytes)
		} else {
			result, err = ctx.BroadcastTxSync(stdTxBytes)
		}

		if err != nil {
			util.HandleError(err, w)
			return
		}

		w.Write(result)
	}
}
