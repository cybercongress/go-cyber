package rank

import (
	"crypto/sha256"
	"encoding/binary"
	"github.com/cybercongress/cyberd/merkle"
	"math"
)

type Rank struct {
	Values     []float64
	MerkleTree *merkle.Tree
}

func NewRank(values []float64, fullTree bool) Rank {
	merkleTree := merkle.NewTree(sha256.New(), fullTree)
	for _, f64 := range values {
		rankBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(rankBytes, math.Float64bits(f64))
		merkleTree.Push(rankBytes)
	}
	return Rank{Values: values, MerkleTree: merkleTree}
}

func (r *Rank) Reset() {
	r.Values = nil
	r.MerkleTree = nil
}

// TODO: optimize. Possible solution: adjust capacity of rank values slice after rank calculation
func (r *Rank) AddNewCids(currentCidCount uint64) {
	if r.Values == nil {
		return
	}
	newCidsCount := currentCidCount - uint64(len(r.Values))
	// add new cids with rank = 0
	r.Values = append(r.Values, make([]float64, newCidsCount)...)
	// extend merkle tree
	zeroRankBytes := make([]byte, 8)
	for i := uint64(0); i < newCidsCount; i++ {
		r.MerkleTree.Push(zeroRankBytes)
	}
}
