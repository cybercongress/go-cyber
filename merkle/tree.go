package merkle

import (
	"bytes"
	"hash"
	"math"
)

// Merkle tree data structure based on RFC-6962 standard (https://tools.ietf.org/html/rfc6962#section-2.1)
// we separate whole tree to subtrees where nodes count equal power of 2
// root hash calculates from right to left by summing subtree roots hashes.
type Tree struct {
	// this tree subtrees start from lowest height (extreme right subtree)
	subTree *Subtree

	// last index of elements in this tree (exclusive)
	lastIndex int

	hashF hash.Hash // DON'T USE IT FOR PARALLEL CALCULATION (results in errors)
}

func NewTree(hashF hash.Hash) Tree {
	return Tree{hashF: hashF}
}

func (t *Tree) joinAllSubtrees() {

	for t.subTree.left != nil && t.subTree.height == t.subTree.left.height {

		newSubtreeRoot := &Node{
			hash:       sum(t.hashF, t.subTree.left.root.hash, t.subTree.root.hash),
			parent:     nil,
			left:       t.subTree.left.root,
			right:      t.subTree.root,
			firstIndex: t.subTree.left.root.firstIndex,
			lastIndex:  t.subTree.root.lastIndex,
		}

		newSubtreeRoot.left.parent = newSubtreeRoot
		newSubtreeRoot.right.parent = newSubtreeRoot

		t.subTree = &Subtree{
			root:   newSubtreeRoot,
			right:  nil,
			left:   t.subTree.left.left,
			height: t.subTree.height + 1,
			hashF:  t.hashF,
		}

		if t.subTree.left != nil {
			t.subTree.left.right = t.subTree
		}

	}
}

func (t *Tree) Reset() {
	t.lastIndex = 0
	t.subTree = nil
}

// build completely new tree with data
// works the same (by time) as using Push method one by one
func (t *Tree) BuildNew(data [][]byte) {
	t.Reset()
	itemsLeft := int64(len(data))

	nextSubtreeLen := int64(math.Pow(2, float64(int64(math.Log2(float64(itemsLeft))))))
	startIndex := int64(0)
	endIndex := startIndex + nextSubtreeLen

	for nextSubtreeLen != 0 {

		nextSubtree := buildSubTree(t.hashF, int(startIndex), data[startIndex:endIndex])

		if t.subTree != nil {
			t.subTree.right = nextSubtree
			nextSubtree.left = t.subTree
			t.subTree = nextSubtree
		} else {
			t.subTree = nextSubtree
		}

		itemsLeft = itemsLeft - nextSubtreeLen
		nextSubtreeLen = int64(math.Pow(2, float64(int64(math.Log2(float64(itemsLeft))))))
		startIndex = endIndex
		endIndex = startIndex + nextSubtreeLen
	}

	t.lastIndex = int(endIndex)

}

// n*log(n)
func (t *Tree) Push(data []byte) {

	newSubtreeRoot := &Node{
		hash:       sum(t.hashF, data),
		parent:     nil,
		left:       nil,
		right:      nil,
		firstIndex: t.lastIndex,
		lastIndex:  t.lastIndex,
	}

	t.lastIndex++

	t.subTree = &Subtree{
		root:   newSubtreeRoot,
		right:  nil,
		left:   t.subTree,
		height: 0,
		hashF:  t.hashF,
	}

	if t.subTree.left != nil {
		t.subTree.left.right = t.subTree
	}

	t.joinAllSubtrees()
}

// going from right trees to left
func (t *Tree) GetIndexProofs(i int) []Proof {

	proofs := make([]Proof, 0, int64(math.Log2(float64(t.lastIndex))))

	for current := t.subTree; current != nil; {

		if i >= current.root.firstIndex && i <= current.root.lastIndex {
			proofs = append(proofs, current.root.GetIndexProofs(i)...)
			proofs = append(proofs, current.GetRootProofs()...)
			return proofs
		}

		current = current.left
	}

	return proofs
}

func (t *Tree) ValidateIndex(i int, data []byte) bool {
	return t.ValidateIndexByProofs(i, data, t.GetIndexProofs(i))
}

func (t *Tree) ValidateIndexByProofs(i int, data []byte, proofs []Proof) bool {

	rootHash := sum(t.hashF, data)
	for _, proof := range proofs {
		rootHash = proof.SumWith(t.hashF, rootHash)
	}

	return bytes.Equal(rootHash, t.RootHash())
}

// root hash calculates from right to left by summing subtree roots hashes.
func (t *Tree) RootHash() []byte {

	if t.subTree == nil {
		return nil
	}

	rootHash := t.subTree.root.hash
	current := t.subTree.left

	for current != nil {
		rootHash = sum(t.hashF, rootHash, current.root.hash)
		current = current.left
	}

	return rootHash
}
