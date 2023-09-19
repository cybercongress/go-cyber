package merkle

import (
	"crypto/sha256"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPushAndProofs(t *testing.T) {
	tree := NewTree(sha256.New(), true)

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

func TestBuildNewAndProofs(t *testing.T) {
	tree := NewTree(sha256.New(), true)

	allData := make([][]byte, 0, 31)

	for i := 0; i < 31; i++ {
		data := make([]byte, 8)
		binary.LittleEndian.PutUint64(data, uint64(i))
		allData = append(allData, data)
	}

	tree.BuildNew(allData)

	// Check all proofs
	for i := 0; i < 31; i++ {
		proofs := tree.GetIndexProofs(i)
		data := make([]byte, 8)
		binary.LittleEndian.PutUint64(data, uint64(i))
		require.Equal(t, true, tree.ValidateIndexByProofs(i, data, proofs))
	}
}

func TestEqualityOfBuildNewAndPush(t *testing.T) {
	tree1 := NewTree(sha256.New(), true)

	data := make([]byte, 8)

	for i := 0; i < 31; i++ {
		binary.LittleEndian.PutUint64(data, uint64(i))
		tree1.Push(data)
	}

	tree2 := NewTree(sha256.New(), true)

	allData := make([][]byte, 0, 31)

	for i := 0; i < 31; i++ {
		data := make([]byte, 8)
		binary.LittleEndian.PutUint64(data, uint64(i))
		allData = append(allData, data)
	}

	tree2.BuildNew(allData)

	require.Equal(t, tree1.RootHash(), tree2.RootHash())
}

func TestNotFull(t *testing.T) {
	tree1 := NewTree(sha256.New(), true)

	data := make([]byte, 8)

	for i := 0; i < 31; i++ {
		binary.LittleEndian.PutUint64(data, uint64(i))
		tree1.Push(data)
	}

	tree2 := NewTree(sha256.New(), false)

	for i := 0; i < 31; i++ {
		binary.LittleEndian.PutUint64(data, uint64(i))
		tree2.Push(data)
	}

	tree3 := NewTree(sha256.New(), false)

	allData := make([][]byte, 0, 31)

	for i := 0; i < 31; i++ {
		data := make([]byte, 8)
		binary.LittleEndian.PutUint64(data, uint64(i))
		allData = append(allData, data)
	}

	tree3.BuildNew(allData)

	require.Equal(t, tree1.RootHash(), tree2.RootHash())
	require.Equal(t, tree1.RootHash(), tree3.RootHash())
}

func TestExportImport(t *testing.T) {
	tree1 := NewTree(sha256.New(), true)

	data := make([]byte, 8)

	for i := 0; i < 31; i++ {
		binary.LittleEndian.PutUint64(data, uint64(i))
		tree1.Push(data)
	}

	subtreeRoots := tree1.ExportSubtreesRoots()

	tree2 := NewTree(sha256.New(), false)
	tree2.ImportSubtreesRoots(subtreeRoots)

	require.Equal(t, tree1.RootHash(), tree2.RootHash())

	binary.LittleEndian.PutUint64(data, uint64(31))
	tree1.Push(data)
	tree2.Push(data)
	require.Equal(t, tree1.RootHash(), tree2.RootHash())
}
