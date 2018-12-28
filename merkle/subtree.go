package merkle

import "hash"

// subtrees always should have power of 2 number of elements.
// tree could contain few of subtrees.
// root hash calculates from right to left by summing subtree roots hashes.
type Subtree struct {
	root *Node // root node

	left  *Subtree // left subtree from this one
	right *Subtree  // right  subtree from this one

	height int // height of subtree
	// hash function to hash sum nodes and hash data
	hashF hash.Hash
}

// get proofs for root of this subtree
func (t *Subtree) GetRootProofs() []Proof {
	proofs := make([]Proof, 0)

	proofs = append(proofs, t.getRightSubtreesProof()...)
	proofs = append(proofs, t.getLeftSubtreesProof()...)

	return proofs
}

func (t *Subtree) getLeftSubtreesProof() []Proof {
	proofs := make([]Proof, 0)
	current := t.left
	for current != nil {
		proofs = append(proofs, Proof{hash: current.root.hash, leftSide: false})
		current = current.left
	}
	return proofs
}

// right proof is only one cause we have to sum all right subtrees
// we have to sum hashes from right to left
func (t *Subtree) getRightSubtreesProof() []Proof {

	if t.right == nil {
		return make([]Proof, 0)
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
		proofHash = sum(t.hashF, proofHash, hashesToSum[i])
	}

	return []Proof{{hash: proofHash, leftSide: true}}
}
