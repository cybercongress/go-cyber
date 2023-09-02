package merkle

import (
	"hash"
	"math"
)

func sum(h hash.Hash, data ...[]byte) []byte {
	h.Reset()
	for _, d := range data {
		// the Hash interface specifies that Write never returns an error
		_, _ = h.Write(d)
	}
	return h.Sum(nil)
}

// number of data elements should be power of 2
// not suitable for parallel calculations cause using same hash.Hash
func buildSubTree(h hash.Hash, full bool, startIndex int, data [][]byte) *Subtree {
	nodes := make([]*Node, len(data))
	for i := 0; i < len(data); i++ {
		nodes[i] = &Node{
			hash:       sum(h, data[i]),
			firstIndex: startIndex + i,
			lastIndex:  startIndex + i,
		}
	}

	root := sumNodes(h, full, nodes)[0]

	return &Subtree{
		root:   root,
		left:   nil,
		right:  nil,
		height: int(math.Log2(float64(len(data)))),
		hashF:  h,
	}
}

func sumNodes(h hash.Hash, full bool, nodes []*Node) []*Node {
	if len(nodes) == 1 {
		return nodes
	}

	newNodes := make([]*Node, len(nodes)/2)
	for i := 0; i < len(nodes); i += 2 {
		newNodes[i/2] = joinNodes(h, full, nodes[i], nodes[i+1])
	}

	return sumNodes(h, full, newNodes)
}

func joinNodes(h hash.Hash, full bool, left *Node, right *Node) *Node {
	newNode := &Node{
		firstIndex: left.firstIndex,
		lastIndex:  right.lastIndex,
		hash:       sum(h, left.hash, right.hash),
	}

	if full {
		newNode.left = left
		newNode.right = right
		left.parent = newNode
		right.parent = newNode
	}
	return newNode
}
