package rpc

import (
	"github.com/cybercongress/cyberd/merkle"
)

type RankAndProofResult struct {
	Proofs []merkle.Proof `json:"proofs"`
	rankval float64  `json:"rankvalue"`
}

func Rank(cid string, proof bool) (*RankAndProofResult, error) {
	rankValue, proofs, err := cyberdApp.Rank(cid, proof)
	return &RankAndProofResult{proofs, rankValue}, err
}
