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
	CidCount   uint64
}

func NewRank(values []float64, fullTree bool) Rank {
	merkleTree := merkle.NewTree(sha256.New(), fullTree)
	for _, f64 := range values {
		rankBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(rankBytes, math.Float64bits(f64))
		merkleTree.Push(rankBytes)
	}
	return Rank{Values: values, MerkleTree: merkleTree, CidCount: uint64(len(values))}
}

func NewFromMerkle(cidCount uint64, treeBytes []byte) Rank {
	rank := Rank{
		Values:     nil,
		MerkleTree: merkle.NewTree(sha256.New(), false),
		CidCount:   cidCount,
	}

	rank.MerkleTree.ImportSubtreesRoots(treeBytes)
	return rank
}

func (r Rank) IsEmpty() bool {
	return (r.Values == nil || len(r.Values) == 0) && r.MerkleTree == nil
}

func (r *Rank) Clear() {
	r.Values = nil
	r.MerkleTree = nil
	r.CidCount = 0
}

func (r *Rank) CopyWithoutTree() Rank {

	if r.Values == nil {
		return Rank{Values: nil, MerkleTree: nil, CidCount: 0}
	}

	copiedValues := make([]float64, len(r.Values))
	n := copy(copiedValues, r.Values)
	// SHOULD NOT HAPPEN
	if n != len(r.Values) {
		panic("Not all rank values have been copied")
	}
	return Rank{Values: copiedValues, MerkleTree: nil, CidCount: r.CidCount}
}

// TODO: optimize. Possible solution: adjust capacity of rank values slice after rank calculation
func (r *Rank) AddNewCids(currentCidCount uint64) {

	newCidsCount := currentCidCount - r.CidCount

	// add new cids with rank = 0
	if r.Values != nil {
		r.Values = append(r.Values, make([]float64, newCidsCount)...)
	}

	// extend merkle tree
	zeroRankBytes := make([]byte, 8)
	for i := uint64(0); i < newCidsCount; i++ {
		r.MerkleTree.Push(zeroRankBytes)
	}

	r.CidCount = currentCidCount
}
