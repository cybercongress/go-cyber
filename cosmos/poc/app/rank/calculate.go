package rank

import (
	. "github.com/cybercongress/cyberd/cosmos/poc/app/storage"
	"math"
	"sync"
)

const (
	d         = 0.85
	tolerance = 0.1
)

func CalculateRank(data *InMemoryStorage) ([]float64, int) {

	size := data.GetCidsCount()

	if size == 0 {
		return []float64{}, 0
	}
	inverseOfSize := 1.0 / float64(size)

	prevrank := make([]float64, size)

	tOverSize := (1.0 - d) / float64(size)
	danglingNodes := calculateDanglingNodes(data)

	for i := range danglingNodes {
		prevrank[i] = inverseOfSize
	}

	change := 2.0

	steps := 0
	var rank []float64
	for change > tolerance {
		rank = step(tOverSize, prevrank, danglingNodes, data)
		change = calculateChange(prevrank, rank)
		prevrank = rank
		steps++
	}

	return rank, steps
}

func calculateDanglingNodes(data *InMemoryStorage) []int64 {

	outLinks := data.GetOutLinks()
	danglingNodes := make([]int64, 0, len(outLinks))

	for i, cidLinks := range outLinks {
		if len(cidLinks) == 0 {
			danglingNodes = append(danglingNodes, int64(i))
		}
	}

	return danglingNodes
}

func step(tOverSize float64, prevrank []float64, danglingNodes []int64, data *InMemoryStorage) []float64 {

	innerProduct := 0.0
	for _, danglingNode := range danglingNodes {
		innerProduct += prevrank[danglingNode]
	}

	innerProductOverSize := innerProduct / float64(len(prevrank))
	rank := make([]float64, len(prevrank))

	var wg sync.WaitGroup
	wg.Add(len(data.GetInLinks()))

	for i, inLinksForI := range data.GetInLinks() {

		go func(cid CidNumber, inLinks CidLinks) {
			defer wg.Done()
			ksum := 0.0

			for j := range inLinks {
				linkStake := float64(data.GetOverallLinkStake(CidNumber(j), CidNumber(cid)))
				jCidOutStake := float64(data.GetOverallOutLinksStake(CidNumber(j)))
				ksum += prevrank[j] * (linkStake / jCidOutStake)
			}

			rank[cid] = d*(ksum+innerProductOverSize) + tOverSize
		}(i, inLinksForI)
	}
	wg.Wait()
	return rank
}

func calculateChange(prevrank, rank []float64) float64 {

	acc := 0.0
	for i, pForI := range prevrank {
		acc += math.Abs(pForI - rank[i])
	}

	return acc
}
