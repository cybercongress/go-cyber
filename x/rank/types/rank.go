package types

import (
	"crypto/sha256"
	"encoding/binary"
	"sort"
	"time"

	"github.com/cybercongress/go-cyber/v4/merkle"
	graphtypes "github.com/cybercongress/go-cyber/v4/x/graph/types"

	"github.com/cometbft/cometbft/libs/log"
)

type EMState struct {
	RankValues    []float64
	EntropyValues []float64
	KarmaValues   []float64
}

type Rank struct {
	RankValues    []uint64
	EntropyValues []uint64
	KarmaValues   []uint64
	MerkleTree    *merkle.Tree // ranks merkle
	CidCount      uint64
	TopCIDs       []RankedCidNumber
	NegEntropy    uint64
}

func NewRank(state EMState, logger log.Logger, fullTree bool) Rank {
	// Entropy and karma are experimental features
	start := time.Now()
	particlesCount := uint64(len(state.RankValues))

	rankValues := make([]uint64, particlesCount)
	for i, f64 := range state.RankValues {
		rankValues[i] = uint64(f64 * 1e15)
	}
	state.RankValues = nil

	entropyValues := make([]uint64, particlesCount)
	for i, f64 := range state.EntropyValues {
		entropyValues[i] = uint64(f64 * 1e15)
	}
	negEntropy := float64(0)
	for _, f64 := range state.EntropyValues {
		negEntropy += f64
	}
	state.EntropyValues = nil

	karmaValues := make([]uint64, len(state.KarmaValues))
	for i, f64 := range state.KarmaValues {
		karmaValues[i] = uint64(f64 * 1e15)
	}
	state.KarmaValues = nil

	logger.Info("State processing to integers", "duration", time.Since(start).String())

	start = time.Now()
	merkleTree := merkle.NewTree(sha256.New(), fullTree)
	for _, u64 := range rankValues {
		rankBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(rankBytes, u64)
		merkleTree.Push(rankBytes)
	}
	logger.Info("Rank constructing tree", "duration", time.Since(start).String())

	// NOTE fulltree true if search index enabled
	start = time.Now()
	var newSortedCIDs []RankedCidNumber
	if fullTree {
		newSortedCIDs = BuildTop(rankValues, 1000)
		logger.Info("Build top", "duration", time.Since(start).String())
	}

	logger.Info("-Entropy", "bits", uint64(negEntropy))

	return Rank{
		RankValues:    rankValues,
		EntropyValues: entropyValues,
		KarmaValues:   karmaValues,
		MerkleTree:    merkleTree,
		CidCount:      particlesCount,
		TopCIDs:       newSortedCIDs,
		NegEntropy:    uint64(negEntropy),
	}
}

func NewFromMerkle(cidCount uint64, treeBytes []byte) Rank {
	rank := Rank{
		RankValues:    nil,
		EntropyValues: nil,
		KarmaValues:   nil,
		MerkleTree:    merkle.NewTree(sha256.New(), false),
		CidCount:      cidCount,
		TopCIDs:       nil,
		NegEntropy:    0,
	}

	rank.MerkleTree.ImportSubtreesRoots(treeBytes)
	return rank
}

func (r Rank) IsEmpty() bool {
	return (r.RankValues == nil || len(r.RankValues) == 0) && r.MerkleTree == nil
}

func (r *Rank) Clear() {
	r.RankValues = nil
	r.EntropyValues = nil
	r.KarmaValues = nil
	r.MerkleTree = nil
	r.CidCount = 0
	r.TopCIDs = nil
	r.NegEntropy = 0
}

func (r *Rank) CopyWithoutTree() Rank {
	if r.RankValues == nil {
		return Rank{
			RankValues:    nil,
			EntropyValues: nil,
			KarmaValues:   nil,
			MerkleTree:    nil,
			CidCount:      0,
			TopCIDs:       nil,
			NegEntropy:    0,
		}
	}

	copiedRankValues := make([]uint64, r.CidCount)
	n := copy(copiedRankValues, r.RankValues)
	// SHOULD NOT HAPPEN
	if n != len(r.RankValues) {
		panic("Not all rank values have been copied")
	}

	copiedEntropyValues := make([]uint64, r.CidCount)
	n = copy(copiedEntropyValues, r.EntropyValues)
	if n != len(r.EntropyValues) {
		panic("Not all entropy values have been copied")
	}

	copiedKarmaValues := make([]uint64, len(r.KarmaValues))
	n = copy(copiedKarmaValues, r.KarmaValues)
	if n != len(r.KarmaValues) {
		panic("Not all karma values have been copied")
	}

	copiedTopCIDs := make([]RankedCidNumber, len(r.TopCIDs))
	n = copy(copiedTopCIDs, r.TopCIDs)
	if n != len(r.TopCIDs) {
		panic("Not all sorted particles have been copied")
	}

	return Rank{
		RankValues:    copiedRankValues,
		EntropyValues: copiedEntropyValues,
		KarmaValues:   copiedKarmaValues,
		MerkleTree:    nil,
		CidCount:      r.CidCount,
		TopCIDs:       copiedTopCIDs,
		NegEntropy:    r.NegEntropy,
	}
}

// TODO: optimize. Possible solution: adjust capacity of rank values slice after rank calculation
func (r *Rank) AddNewCids(currentCidCount uint64) {
	newCidsCount := currentCidCount - r.CidCount

	// add new cids with rank = 0
	if r.RankValues != nil {
		r.RankValues = append(r.RankValues, make([]uint64, newCidsCount)...)
	}

	if r.EntropyValues != nil {
		r.EntropyValues = append(r.EntropyValues, make([]uint64, newCidsCount)...)
	}

	// extend merkle tree
	zeroRankBytes := make([]byte, 8)
	for i := uint64(0); i < newCidsCount; i++ {
		r.MerkleTree.Push(zeroRankBytes)
	}

	r.CidCount = currentCidCount
}

func BuildTop(values []uint64, size int) []RankedCidNumber {
	newSortedCIDs := make(sortableCidNumbers, 0, len(values))
	for cid, rank := range values {
		if rank == 0 {
			continue
		}
		newRankedCid := RankedCidNumber{graphtypes.CidNumber(cid), rank}
		newSortedCIDs = append(newSortedCIDs, newRankedCid)
	}
	sort.Stable(sort.Reverse(newSortedCIDs))
	if len(values) > size {
		newSortedCIDs = newSortedCIDs[0:(size - 1)]
	}
	return newSortedCIDs
}
