package rpc

import (
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"

	"github.com/cybercongress/go-cyber/merkle"
)

type RankAndProofResult struct {
	Proofs []merkle.Proof `json:"proofs"`
	Rank   float64        `amino:"unsafe" json:"rank"`
}

func Rank(ctx *rpctypes.Context, cid string, proof bool) (*RankAndProofResult, error) {
	rankValue, proofs, err := cyberdApp.Rank(cid, proof)
	return &RankAndProofResult{proofs, rankValue}, err
}
