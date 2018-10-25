package rank

import (
	. "github.com/cybercongress/cyberd/cosmos/poc/app/storage"
	"sync"
)

const (
	one       uint64 = 1000000000 // represents 1.000000000
	d         uint64 = 850000000  // represents 0.850000000
	tolerance uint64 = 10000000   // represents 0.010000000
)

func CalculateRank(data *InMemoryStorage) ([]uint64, int) {

	size := data.GetCidsCount()

	if size == 0 {
		return []uint64{}, 0
	}

	prevrank := make([]uint64, size)

	tOverSize := (one - d) / size
	danglingNodes := calculateDanglingNodes(data)

	for _, i := range danglingNodes {
		prevrank[i] = tOverSize
	}

	change := 2 * one

	steps := 0
	var rank []uint64
	for change > tolerance {
		rank = step(tOverSize, prevrank, danglingNodes, data)
		change = calculateChange(prevrank, rank)
		prevrank = rank
		steps++
	}

	return rank, steps
}

func calculateDanglingNodes(data *InMemoryStorage) []int64 {

	cidsCount := data.GetCidsCount()
	outLinks := data.GetInLinks()
	danglingNodes := make([]int64, 0)

	i := uint64(0)
	for i < cidsCount {
		if len(outLinks[CidNumber(i)]) == 0 {
			danglingNodes = append(danglingNodes, int64(i))
		}
		i++
	}

	return danglingNodes
}

func step(tOverSize uint64, prevrank []uint64, danglingNodes []int64, data *InMemoryStorage) []uint64 {

	innerProduct := uint64(0)
	for _, danglingNode := range danglingNodes {
		innerProduct += prevrank[danglingNode]
	}

	innerProductOverSize := innerProduct / uint64(len(prevrank))
	rank := append(make([]uint64, 0, len(prevrank)), prevrank...)

	var wg sync.WaitGroup
	wg.Add(len(data.GetInLinks()))

	for i, inLinksForI := range data.GetInLinks() {

		go func(cid CidNumber, inLinks CidLinks) {
			defer wg.Done()
			ksum := uint64(0)

			for j := range inLinks {
				linkStake := data.GetOverallLinkStake(CidNumber(j), CidNumber(cid))
				jCidOutStake := data.GetOverallOutLinksStake(CidNumber(j))
				ksum += prevrank[j] / (jCidOutStake / linkStake)
			}

			// 17/20 = 0.85 = d
			rank[cid] = (ksum+innerProductOverSize)/20*17 + tOverSize
		}(i, inLinksForI)
	}
	wg.Wait()
	return rank
}

func calculateChange(prevrank, rank []uint64) uint64 {

	acc := uint64(0)
	for i, pForI := range prevrank {
		if pForI > rank[i] {
			acc += pForI - rank[i]
		} else {
			acc += rank[i] - pForI
		}
	}

	return acc
}
