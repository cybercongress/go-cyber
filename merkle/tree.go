package merkle

import (
	"bytes"
	"encoding/binary"
	"hash"
	"math"
)

// Merkle tree data structure based on RFC-6962 standard (https://tools.ietf.org/html/rfc6962#section-2.1)
// we separate whole tree to subtrees where nodes count equal power of 2
// root hash calculates from right to left by summing subtree roots hashes.
type Tree struct {
	// this tree subtrees start from lowest height (extreme right subtree)
	subTree *Subtree
	// subtrees total count (for convenience)
	subTreesCount int

	// last index of elements in this tree (exclusive)
	lastIndex int

	// DON'T USE IT FOR PARALLEL CALCULATION (results in errors)
	hashF hash.Hash

	// if false then store only roots of subtrees (no proofs available => suitable for consensus only)
	// for 1,099,511,627,775 links tree would contain only 40 root hashes.
	full bool
}

func NewTree(hashF hash.Hash, full bool) *Tree {
	return &Tree{hashF: hashF, full: full}
}

func (t *Tree) joinAllSubtrees() {

	for t.subTree.left != nil && t.subTree.height == t.subTree.left.height {

		newSubtreeRoot := &Node{
			hash:       sum(t.hashF, t.subTree.left.root.hash, t.subTree.root.hash),
			firstIndex: t.subTree.left.root.firstIndex,
			lastIndex:  t.subTree.root.lastIndex,
		}

		// for full tree we should keep all nodes
		if t.full {
			newSubtreeRoot.left = t.subTree.left.root
			newSubtreeRoot.right = t.subTree.root
			newSubtreeRoot.left.parent = newSubtreeRoot
			newSubtreeRoot.right.parent = newSubtreeRoot
		}

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

		t.subTreesCount--
	}
}

func (t *Tree) Reset() {
	t.lastIndex = 0
	t.subTree = nil
	t.subTreesCount = 0
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

		nextSubtree := buildSubTree(t.hashF, t.full, int(startIndex), data[startIndex:endIndex])

		if t.subTree != nil {
			t.subTree.right = nextSubtree
			nextSubtree.left = t.subTree
			t.subTree = nextSubtree
		} else {
			t.subTree = nextSubtree
		}

		t.subTreesCount++

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

	t.subTreesCount++
	t.joinAllSubtrees()
}

// going from right trees to left
func (t *Tree) GetIndexProofs(i int) []Proof {

	// we cannot build proofs with not full tree
	if !t.full {
		return nil
	}

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
	if !t.full {
		return false
	}
	return t.ValidateIndexByProofs(i, data, t.GetIndexProofs(i))
}

func (t *Tree) ValidateIndexByProofs(i int, data []byte, proofs []Proof) bool {

	rootHash := sum(t.hashF, data)
	for _, proof := range proofs {
		rootHash = proof.SumWith(t.hashF, rootHash)
	}

	return bytes.Equal(rootHash, t.RootHash())
}

// root hash calculates from right to left by summing subtrees root hashes.
func (t *Tree) RootHash() []byte {

	if t.subTree == nil {
		return sum(t.hashF) // zero hash
	}

	rootHash := t.subTree.root.hash
	current := t.subTree.left

	for current != nil {
		rootHash = sum(t.hashF, rootHash, current.root.hash)
		current = current.left
	}

	return rootHash
}

// from right to left
// we need to export root hash and height of tree
// from those bytes we could restore it later and use for consensus
func (t *Tree) ExportSubtreesRoots() []byte {
	if t.subTree == nil {
		return make([]byte, 0)
	}

	hashSize := t.hashF.Size()
	heightSize := 8 // using 8 byte integer
	result := make([]byte, 0, (hashSize+heightSize)*t.subTreesCount)
	current := t.subTree

	for current != nil {
		heightBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(heightBytes, uint64(current.height))

		result = append(result, current.root.hash...)
		result = append(result, heightBytes...)
		current = current.left
	}

	return result
}

// from right to left
// after import we loosing indices (actually they don't need for consensus and pushing)
func (t *Tree) ImportSubtreesRoots(subTreesRoots []byte) {
	t.Reset()
	t.full = false

	hashSize := t.hashF.Size()
	heightSize := 8

	t.subTreesCount = len(subTreesRoots) / (hashSize + heightSize)

	start := 0
	var first *Subtree
	var current *Subtree
	for i := 0; i < t.subTreesCount; i++ {
		end := start + hashSize
		rootHash := subTreesRoots[start:end]
		height := binary.LittleEndian.Uint64(subTreesRoots[end : end+heightSize])
		nextSubtree := &Subtree{
			root: &Node{
				hash: rootHash,
			},
			height: int(height),
		}

		if current != nil {
			nextSubtree.right = current
			current.left = nextSubtree
			current = nextSubtree
		} else {
			current = nextSubtree
			first = nextSubtree
		}

		start = end + heightSize
	}

	t.subTree = first
}
