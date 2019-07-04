package merkle

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

func (n *Node) GetIndexProofs(i int) []Proof {
	proofs := make([]Proof, 0)

	if n.left != nil && i >= n.left.firstIndex && i <= n.left.lastIndex {
		proofs = n.left.GetIndexProofs(i)
		proofs = append(proofs, Proof{Hash: n.right.hash, LeftSide: false})
	}

	if n.right != nil && i >= n.right.firstIndex && i <= n.right.lastIndex {
		proofs = n.right.GetIndexProofs(i)
		proofs = append(proofs, Proof{Hash: n.left.hash, LeftSide: true})
	}

	return proofs
}
