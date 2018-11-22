package common

import (
	"github.com/tendermint/tendermint/lite"
	"github.com/tendermint/tendermint/types"
)

var _ lite.Verifier = (*NoopVerifier)(nil)

type NoopVerifier struct {
	ChainId string
}

func (NoopVerifier) Verify(sheader types.SignedHeader) error {
	return nil
}

func (v NoopVerifier) ChainID() string {
	return v.ChainId
}
