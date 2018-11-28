package core

import "github.com/tendermint/tendermint/crypto/secp256k1"

type Signature struct {
	PubKey        secp256k1.PubKeySecp256k1 `json:"pub_key"` // optional
	Signature     []byte                    `json:"signature"`
	AccountNumber int64                     `json:"account_number"`
	Sequence      int64                     `json:"sequence"`
}
