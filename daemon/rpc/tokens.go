package rpc

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type SendRequest struct {
	Fee        auth.StdFee    `json:"fee"`
	Msgs       []bank.MsgSend `json:"msgs"`
	Signatures []Signature    `json:"signatures"`
	Memo       string         `json:"memo"`
}

func (r SendRequest) GetFee() auth.StdFee        { return r.Fee }
func (r SendRequest) GetSignatures() []Signature { return r.Signatures }
func (r SendRequest) GetMemo() string            { return r.Memo }
func (r SendRequest) GetMsgs() []types.Msg {
	msgs := make([]types.Msg, 0, len(r.Msgs))
	for _, msg := range r.Msgs {
		msgs = append(msgs, msg)
	}
	return msgs
}

func UnmarshalSendRequestFn(reqBytes []byte) (TxRequest, error) {
	var req SendRequest
	err := json.Unmarshal(reqBytes, &req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
