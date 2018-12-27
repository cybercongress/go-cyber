package merkle

import (
	"crypto/sha256"
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

	indexToProof := 25
	proofs := tree.GetIndexProofs(indexToProof)

	require.Equal(t, 5, len(proofs))

	binary.LittleEndian.PutUint64(data, uint64(indexToProof))
	h := sha256.New()
	h.Write(data)

	rootHash :=  h.Sum(nil)

	h.Reset()
	h.Write(proofs[0])
	h.Write(rootHash)

	rootHash = h.Sum(nil)

	h.Reset()
	h.Write(rootHash)
	h.Write(proofs[1])

	rootHash = h.Sum(nil)

	h.Reset()
	h.Write(proofs[2])
	h.Write(rootHash)

	rootHash = h.Sum(nil)

	h.Reset()
	h.Write(rootHash)
	h.Write(proofs[3])

	rootHash = h.Sum(nil)

	h.Reset()
	h.Write(rootHash)
	h.Write(proofs[4])

	rootHash = h.Sum(nil)

	require.Equal(t, rootHash, tree.hash)
}
