package merkle

import (
	"crypto/sha256"
	"math"
)

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

type Subtree struct {
	root *Node

	left  *Subtree
	right *Subtree

	height int
}

func (t *Subtree) GetTreeProofs() [][]byte {
	proofs := make([][]byte, 0)

	proofs = append(proofs, t.getRightProof())
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
func (t *Subtree) getRightProof() []byte {
	h := sha256.New()

	hashesToSum := make([][]byte, 0)

	rightTree := t.right
	for rightTree != nil {
		hashesToSum = append(hashesToSum, rightTree.root.hash)
		rightTree = rightTree.right
	}

	for i := len(hashesToSum) - 1; i >= 0; i-- {
		_, _ = h.Write(hashesToSum[i])
	}

	return h.Sum(nil)
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
		proofs = append(append(proofs, n.right.hash), n.left.GetIndexProofs(i)...)
	}

	if n.right != nil && i  >= n.right.firstIndex && i <= n.right.lastIndex {
		proofs = append(append(proofs, n.left.hash), n.right.GetIndexProofs(i)...)
	}

	return proofs
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
}

// going from right trees to left
func (t *Tree) GetIndexProofs(i int) [][]byte {

	proofs := make([][]byte, 0, int64(math.Log2(float64(t.lastIndex))))

	for current := t.subTree; current != nil; {

		if i >= current.root.firstIndex && i <= current.root.lastIndex {
			proofs = append(proofs, current.GetTreeProofs()...)
			proofs = append(proofs, current.root.GetIndexProofs(i)...)
			return proofs
		}

		current = current.left
	}

	return proofs
}


//[227 176 196 66 152 252 28 20 154 251 244 200 153 111 185 36 39 174 65 228 100 155 147 76 164 149 153 27 120 82 184 85]
//[75 217 49 194 89 99 166 172 221 108 189 32 135 207 23 40 30 96 254 174 103 138 69 33 215 139 153 126 162 232 247 2]
//[159 242 126 58 20 8 10 86 205 90 252 235 41 220 68 66 233 158 235 199 53 219 40 10 111 243 215 180 166 159 143 73]
//[136 251 153 33 90 8 180 55 66 184 146 46 64 60 223 40 135 200 252 202 215 213 222 106 30 140 104 137 177 227 196 23]
//[63 133 40 196 85 15 252 43 31 93 133 83 188 95 168 173 6 123 248 205 181 127 135 107 22 112 206 86 201 87 38 178]