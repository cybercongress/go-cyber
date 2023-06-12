package merkle

import (
	"hash"
)

type Proof struct {
	LeftSide bool   `json:"leftSide"` // where proof should be placed to sum with hash (left or right side)
	Hash     []byte `json:"hash"`
}

// calculate sum hash.
func (p *Proof) SumWith(hashF hash.Hash, hash []byte) []byte {
	if p.LeftSide {
		return sum(hashF, p.Hash, hash)
	} else {
		return sum(hashF, hash, p.Hash)
	}
}
