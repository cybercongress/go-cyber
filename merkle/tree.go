package merkle

import (
	"crypto/sha256"
	"math"
)

type Subtree struct {
	root *Node

	left  *Subtree
	right *Subtree

	height int
}

func (t *Subtree) GetTreeProofs() [][]byte {
	proofs := make([][]byte, 0)

	proofs = append(proofs, t.getRightProofs()...)
	proofs = append(proofs, t.getLeftProofs()...)

	return proofs
}

func (t *Subtree) getLeftProofs() [][]byte {
	proofs := make([][]byte, 0, 1)
	current := t.left
	for current != nil {
		proofs = append(proofs, current.root.hash)
		current = current.left
	}
	return proofs
}

// right proof is only one cause we have to merge all right trees
func (t *Subtree) getRightProofs() [][]byte {

	if t.right == nil {
		return make([][]byte, 0)
	}

	hashesToSum := make([][]byte, 0)

	rightTree := t.right
	for rightTree != nil {
		hashesToSum = append(hashesToSum, rightTree.root.hash)
		rightTree = rightTree.right
	}

	n := len(hashesToSum) - 1
	proofHash := hashesToSum[n]
	for i := n - 1; i >= 0; i-- {
		h := sha256.New()
		h.Write(proofHash)
		h.Write(hashesToSum[i])
		proofHash = h.Sum(nil)
	}

	return [][]byte{proofHash}
}

type Node struct {
	hash []byte

	parent *Node
	left   *Node
	right  *Node

	// first index of elements in this node subnodes (inclusive)
	firstIndex int
	// last index of elements in this node subnodes (inclusive)
	lastIndex int
}

func (n *Node) GetIndexProofs(i int) [][]byte {
	proofs := make([][]byte, 0)

	if n.left != nil && i >= n.left.firstIndex && i <= n.left.lastIndex {
		proofs = n.left.GetIndexProofs(i)
		proofs = append(proofs, n.right.hash)
	}

	if n.right != nil && i >= n.right.firstIndex && i <= n.right.lastIndex {
		proofs = n.right.GetIndexProofs(i)
		proofs = append(proofs, n.left.hash)
	}

	return proofs
}

// we separate whole tree to sub trees where nodes count equal power of 2
type Tree struct {
	// this tree subtrees start from lowest height (from last right subtree)
	subTree *Subtree

	// first index of elements in this tree (inclusive)
	firstIndex int
	// last index of elements in this tree (exclusive)
	lastIndex int

	hash []byte
}

func NewTree() Tree {
	return Tree{}
}

// only for root tree
func (t *Tree) joinAllSubtrees() {

	for t.subTree.left != nil && t.subTree.height == t.subTree.left.height {

		newRootHash := sha256.New()
		newRootHash.Write(t.subTree.left.root.hash)
		newRootHash.Write(t.subTree.root.hash)

		newSubtreeRoot := &Node{
			hash:       newRootHash.Sum(nil),
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
		}

		if t.subTree.left != nil {
			t.subTree.left.right = t.subTree
		}

	}
}

// n*log(n)
func (t *Tree) Push(data []byte) {

	hash := sha256.New()
	hash.Write(data)

	newSubtreeRoot := &Node{
		hash:       hash.Sum(nil),
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
	}

	if t.subTree.left != nil {
		t.subTree.left.right = t.subTree
	}

	t.joinAllSubtrees()

	t.hash = t.GetRootHash()
}

// going from right trees to left
func (t *Tree) GetIndexProofs(i int) [][]byte {

	proofs := make([][]byte, 0, int64(math.Log2(float64(t.lastIndex))))

	for current := t.subTree; current != nil; {

		if i >= current.root.firstIndex && i <= current.root.lastIndex {
			proofs = append(proofs, current.root.GetIndexProofs(i)...)
			proofs = append(proofs, current.GetTreeProofs()...)
			return proofs
		}

		current = current.left
	}

	return proofs
}

func (t *Tree) GetRootHash() []byte {

	if t.subTree == nil {
		return nil
	}

	rootHash := t.subTree.root.hash
	current := t.subTree.left

	for current != nil {
		h := sha256.New()
		h.Write(rootHash)
		h.Write(current.root.hash)

		rootHash = h.Sum(nil)
		current = current.left
	}

	return rootHash
}
