package merkle

import (
	"hash"
)

type Proof struct {
	leftSide bool // where proof should be placed to sum with hash (left or right side)
	hash     []byte
}

// calculate sum hash
func (p *Proof) SumWith(hashF hash.Hash, hash []byte) []byte {

	if p.leftSide {
		return sum(hashF, p.hash, hash)
	} else {
		return sum(hashF, hash, p.hash)
	}
}
