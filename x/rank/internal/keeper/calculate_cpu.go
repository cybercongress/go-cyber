package keeper

import (
	"github.com/cybercongress/cyberd/x/link"
	"github.com/cybercongress/cyberd/x/rank/internal/types"

	"sync"
)

const (
	d         float64 = 0.85
	tolerance float64 = 1e-3
)

func calculateRankCPU(ctx *types.CalculationContext) []float64 {

	inLinks := ctx.GetInLinks()

	size := ctx.GetCidsCount()
	if size == 0 {
		return []float64{}
	}

	rank := make([]float64, size)
	defaultRank := (1.0 - d) / float64(size)
	danglingNodesSize := uint64(0)

	for i := range rank {
		rank[i] = defaultRank
		if len(inLinks[link.CidNumber(i)]) == 0 {
			danglingNodesSize++
		}
	}

	innerProductOverSize := defaultRank * (float64(danglingNodesSize) / float64(size))
	defaultRankWithCorrection := float64(d*innerProductOverSize) + defaultRank

	change := tolerance + 1

	steps := 0
	prevrank := make([]float64, 0)
	prevrank = append(prevrank, rank...)
	for change > tolerance {
		rank = step(ctx, defaultRankWithCorrection, prevrank)
		change = calculateChange(prevrank, rank)
		prevrank = rank
		steps++
	}

	return rank
}

func step(ctx *types.CalculationContext, defaultRankWithCorrection float64, prevrank []float64) []float64 {

	rank := append(make([]float64, 0, len(prevrank)), prevrank...)

	var wg sync.WaitGroup
	wg.Add(len(ctx.GetInLinks()))

	for cid := range ctx.GetInLinks() {

		go func(i link.CidNumber) {
			defer wg.Done()
			_, sortedCids, ok := ctx.GetSortedInLinks(i)

			if !ok {
				return
			} else {
				ksum := float64(0)
				for _, j := range sortedCids {
					linkStake := getOverallLinkStake(ctx, j, i)
					jCidOutStake := getOverallOutLinksStake(ctx, j)
					weight := float64(linkStake) / float64(jCidOutStake)
					ksum = float64(prevrank[j]*weight) + ksum //force no-fma here by explicit conversion
				}

				rank[i] = float64(ksum*d) + defaultRankWithCorrection //force no-fma here by explicit conversion
			}

		}(link.CidNumber(cid))
	}
	wg.Wait()
	return rank
}

func getOverallLinkStake(ctx *types.CalculationContext, from link.CidNumber, to link.CidNumber) uint64 {

	stake := uint64(0)
	users := ctx.GetOutLinks()[from][to]
	for user := range users {
		stake += ctx.GetStakes()[user]
	}
	return stake
}

func getOverallOutLinksStake(ctx *types.CalculationContext, from link.CidNumber) uint64 {

	stake := uint64(0)
	for to := range ctx.GetOutLinks()[from] {
		stake += getOverallLinkStake(ctx, from, to)
	}
	return stake
}

func calculateChange(prevrank, rank []float64) float64 {

	maxDiff := 0.0
	diff := 0.0
	for i, pForI := range prevrank {
		if pForI > rank[i] {
			diff = pForI - rank[i]
		} else {
			diff = rank[i] - pForI
		}
		if diff > maxDiff {
			maxDiff = diff
		}
	}

	return maxDiff
}
