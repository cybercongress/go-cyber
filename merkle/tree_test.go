package merkle

import (
	"encoding/binary"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProofs(t *testing.T) {

	tree := NewTree()

	data := make([]byte, 8)

	for i := 0; i < 31; i++ {
		binary.LittleEndian.PutUint64(data, uint64(i))
		tree.Push(data)
	}

	// Check all proofs
	for i := 0; i < 31; i++ {
		proofs := tree.GetIndexProofs(i)
		binary.LittleEndian.PutUint64(data, uint64(i))
		require.Equal(t, true, tree.ValidateIndexByProofs(i, data, proofs))
	}

}
