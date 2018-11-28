package core

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cybercongress/cyberd/app"
)

type TxRequest interface {
	GetFee() auth.StdFee
	GetSignatures() []Signature
	GetMemo() string
	GetMsgs() []types.Msg
}

type UnmarshalTxRequest func([]byte) (TxRequest, error)

type LinkRequest struct {
	Fee        auth.StdFee   `json:"fee"`
	Msgs       []app.MsgLink `json:"msgs"`
	Signatures []Signature   `json:"signatures"`
	Memo       string        `json:"memo"`
}

func (r LinkRequest) GetFee() auth.StdFee        { return r.Fee }
func (r LinkRequest) GetSignatures() []Signature { return r.Signatures }
func (r LinkRequest) GetMemo() string            { return r.Memo }
func (r LinkRequest) GetMsgs() []types.Msg {
	msgs := make([]types.Msg, 0, len(r.Msgs))
	for _, msg := range r.Msgs {
		msgs = append(msgs, msg)
	}
	return msgs
}

func UnmarshalLinkRequest(reqBytes []byte) (TxRequest, error) {
	var req LinkRequest
	err := json.Unmarshal(reqBytes, &req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

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

func UnmarshalSendRequest(reqBytes []byte) (TxRequest, error) {
	var req SendRequest
	err := json.Unmarshal(reqBytes, &req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
