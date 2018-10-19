package proxy

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/cyberd/cosmos/poc/app"
	"io/ioutil"
	"net/http"
)

type LinkRequest struct {
	Fee        auth.StdFee   `json:"fee"`
	Msgs       []app.MsgLink `json:"msgs"`
	Signatures []Signature   `json:"signatures"`
	Memo       string        `json:"memo"`
}

func LinkHandlerFn(ctx ProxyContext) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		requestBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			HandleError(err, w)
			return
		}

		var request LinkRequest
		err = json.Unmarshal(requestBytes, &request)
		if err != nil {
			HandleError(err, w)
			return
		}

		// BUILDING COSMOS SDK TX
		signatures := make([]auth.StdSignature, 0, len(request.Signatures))
		for _, sig := range request.Signatures {
			stdSig := auth.StdSignature{
				PubKey: sig.PubKey, Signature: sig.Signature, AccountNumber: sig.AccountNumber, Sequence: sig.Sequence,
			}
			signatures = append(signatures, stdSig)
		}

		msgs := make([]sdk.Msg, 0, len(request.Msgs))
		for _, msg := range request.Msgs {
			msgs = append(msgs, msg)
		}

		stdTx := auth.StdTx{Msgs: msgs, Fee: request.Fee, Signatures: signatures, Memo: request.Memo}

		stdTxBytes, err := ctx.Codec.MarshalBinary(stdTx)
		if err != nil {
			HandleError(err, w)
			return
		}

		resp, err := ctx.Node.BroadcastTxCommit(stdTxBytes)
		if err != nil {
			HandleError(err, w)
			return
		}

		respBytes, err := json.Marshal(resp)
		if err != nil {
			HandleError(err, w)
			return
		}
		w.Write(respBytes)
	}
}
