package merkle

import "crypto/sha256"

type Proof struct {
	left bool // where proof should be placed for concat with hash (left or right)
	hash []byte
}

// calculate sum hash
func (p *Proof) ConcatWith(hash []byte) []byte {
	h := sha256.New()

	if p.left {
		h.Write(p.hash)
		h.Write(hash)
	} else {
		h.Write(hash)
		h.Write(p.hash)
	}

	return h.Sum(nil)
}
