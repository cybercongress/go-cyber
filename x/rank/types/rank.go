package types

import (
	"crypto/sha256"
	"encoding/binary"
	//"fmt"
	"sort"
	"time"

	"github.com/cybercongress/go-cyber/merkle"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"

	"github.com/tendermint/tendermint/libs/log"
)

type EMState struct {
	RankValues       []float64
	EntropyValues  	 []float64
	LuminosityValues []float64
	KarmaValues		 []float64
}

type Rank struct {
	RankValues 		 []uint64
	EntropyValues 	 []uint64
	LuminosityValues []uint64
	KarmaValues 	 []uint64
	MerkleTree 		 *merkle.Tree // ranks merkle
	CidCount   		 uint64
	TopCIDs	   		 []RankedCidNumber
}

func NewRank(state EMState, logger log.Logger, fullTree bool) Rank {

	// TODO separate mul impl for cpu and gpu
	start := time.Now()
	cidsCount := uint64(len(state.RankValues))

	rankValues := make([]uint64, cidsCount)
	for i, f64 := range state.RankValues {
		rankValues[i] = uint64(f64*1e20)
	}
	state.RankValues = make([]float64, 0)

	entropyValues := make([]uint64, cidsCount)
	for i, f64 := range state.EntropyValues {
		entropyValues[i] = uint64(f64*1e20)
	}
	state.EntropyValues = make([]float64, 0)

	luminosityValues := make([]uint64, cidsCount)
	for i, f64 := range state.LuminosityValues {
		luminosityValues[i] = uint64(f64*1e20)
	}
	state.LuminosityValues = make([]float64, 0)

	karmaValues := make([]uint64, len(state.KarmaValues))
	for i, f64 := range state.KarmaValues {
		karmaValues[i] = uint64(f64*1e20)
	}
	state.KarmaValues = make([]float64, 0)

	logger.Info("EMState constructing to uint", "time", time.Since(start))

	start = time.Now()
	merkleTree := merkle.NewTree(sha256.New(), fullTree)
	for _, u64 := range rankValues {
		rankBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(rankBytes, u64)
		merkleTree.Push(rankBytes)
	}
	logger.Info("Rank constructing tree", "time", time.Since(start))

	// NOTE fulltree true if search index enabled
	var newSortedCIDs []RankedCidNumber
	if (fullTree == true) {
		newSortedCIDs = BuildTop(rankValues, 1000)
	}

	return Rank{
		RankValues: 	rankValues,
		EntropyValues:  entropyValues,
		LuminosityValues: luminosityValues,
		KarmaValues: 	karmaValues,
		MerkleTree: 	merkleTree,
		CidCount: 		cidsCount,
		TopCIDs: 		newSortedCIDs,
	}
}

func NewFromMerkle(cidCount uint64, treeBytes []byte) Rank {
	rank := Rank{
		RankValues:       nil,
		EntropyValues:    nil,
		LuminosityValues: nil,
		KarmaValues:      nil,
		MerkleTree: 	  merkle.NewTree(sha256.New(), false),
		CidCount:   	  cidCount,
		TopCIDs:    	  nil,
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
	r.LuminosityValues = nil
	r.KarmaValues = nil
	r.MerkleTree = nil
	r.CidCount = 0
	r.TopCIDs = nil
}

func (r *Rank) CopyWithoutTree() Rank {

	if r.RankValues == nil {
		return Rank{
			RankValues: nil,
			EntropyValues: nil,
			LuminosityValues: nil,
			KarmaValues: nil,
			MerkleTree: nil,
			CidCount: 0,
			TopCIDs: nil,
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
	// SHOULD NOT HAPPEN
	if n != len(r.EntropyValues) {
		panic("Not all entropy values have been copied")
	}

	copiedLuminosityValues := make([]uint64, r.CidCount)
	n = copy(copiedLuminosityValues, r.LuminosityValues)
	// SHOULD NOT HAPPEN
	if n != len(r.LuminosityValues) {
		panic("Not all luminosity values have been copied")
	}

	copiedKarmaValues := make([]uint64, len(r.KarmaValues))
	n = copy(copiedKarmaValues, r.KarmaValues)
	// SHOULD NOT HAPPEN
	if n != len(r.KarmaValues) {
		panic("Not all karma values have been copied")
	}

	return Rank{
		RankValues: copiedRankValues,
		EntropyValues: copiedEntropyValues,
		LuminosityValues: copiedLuminosityValues,
		KarmaValues: copiedKarmaValues,
		MerkleTree: nil,
		CidCount: r.CidCount,
		TopCIDs: r.TopCIDs,
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

	if r.LuminosityValues != nil {
		r.LuminosityValues = append(r.LuminosityValues, make([]uint64, newCidsCount)...)
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
		if (rank == 0) { continue }
		newRankedCid := RankedCidNumber{graphtypes.CidNumber(cid), rank}
		newSortedCIDs = append(newSortedCIDs, newRankedCid)
	}
	sort.Stable(sort.Reverse(newSortedCIDs))
	if (len(values) > size) {
		newSortedCIDs = newSortedCIDs[0:(size-1)]
	}
	return newSortedCIDs
}
