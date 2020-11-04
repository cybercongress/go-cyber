package rpc

import (
	"github.com/cybercongress/go-cyber/merkle"
)

type RankAndProofResult struct {
	Proofs []merkle.Proof `json:"proofs"`
	Rank   uint64          `json:"rank"`
}

func Rank(cid string, proof bool) (*RankAndProofResult, error) {
	rankValue, proofs, err := cyberdApp.Rank(cid, proof)
	return &RankAndProofResult{proofs, rankValue}, err
}
