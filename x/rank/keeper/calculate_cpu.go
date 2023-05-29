package keeper

import (
	//"fmt"
	//"encoding/binary"
	"math"
	//"math/big"
	graphtypes "github.com/cybercongress/go-cyber/v2/x/graph/types"
	"github.com/cybercongress/go-cyber/v2/x/rank/types"
)

func calculateRankCPU(ctx *types.CalculationContext) types.EMState {
	inLinks := ctx.GetInLinks()
	tolerance := ctx.GetTolerance()
	dampingFactor := ctx.GetDampingFactor()

	size := ctx.GetCidsCount()
	if size == 0 || len(ctx.GetStakes()) == 0 {
		return types.EMState{
			[]float64{},
			[]float64{},
			[]float64{},
		}
	}

	rank := make([]float64, size)
	entropy := make([]float64, size)
	karma := make([]float64, len(ctx.GetStakes()))
	defaultRank := (1.0 - dampingFactor) / float64(size)
	danglingNodesSize := uint64(0)

	for i := range rank {
		rank[i] = defaultRank
		if len(inLinks[graphtypes.CidNumber(i)]) == 0 {
			danglingNodesSize++
		}
	}

	innerProductOverSize := defaultRank * (float64(danglingNodesSize) / float64(size))
	defaultRankWithCorrection := float64(dampingFactor*innerProductOverSize) + defaultRank

	change := tolerance + 1

	steps := 0
	prevrank := make([]float64, 0)
	prevrank = append(prevrank, rank...)
	for change > tolerance {
		rank = step(ctx, defaultRankWithCorrection, dampingFactor, prevrank)
		change = calculateChange(prevrank, rank)
		prevrank = rank
		steps++
	}

	// experimental features, out of consensus, available with API
	entropyCalc(ctx, entropy, size, dampingFactor)
	karmaCalc(ctx, rank, entropy, karma)

	return types.EMState{
		rank,
		entropy,
		karma,
	}
}

func step(ctx *types.CalculationContext, defaultRankWithCorrection float64, dampingFactor float64, prevrank []float64) []float64 {
	rank := append(make([]float64, 0, len(prevrank)), prevrank...)

	for cid := range ctx.GetInLinks() {
		_, sortedCids, ok := ctx.GetSortedInLinks(cid)

		if !ok {
			continue
		} else {
			ksum := float64(0)
			for _, j := range sortedCids {
				linkStake := getOverallLinkStake(ctx, j, cid)
				jCidOutStake := getOverallOutLinksStake(ctx, j)
				if linkStake == 0 || jCidOutStake == 0 {
					continue
				}
				weight := float64(linkStake) / float64(jCidOutStake)
				// if math.IsNaN(weight) { weight = float64(0) }
				ksum = prevrank[j]*weight + ksum // force no-fma here by explicit conversion
			}
			rank[cid] = ksum*dampingFactor + defaultRankWithCorrection // force no-fma here by explicit conversion
		}
	}

	return rank
}

func getOverallLinkStake(ctx *types.CalculationContext, from graphtypes.CidNumber, to graphtypes.CidNumber) uint64 {
	stake := uint64(0)
	users := ctx.GetOutLinks()[from][to]
	for user := range users {
		// stake += ctx.GetStakes()[uint64(user)]
		stake += getNormalizedStake(ctx, uint64(user))
	}
	return stake
}

func getOverallOutLinksStake(ctx *types.CalculationContext, from graphtypes.CidNumber) uint64 {
	stake := uint64(0)
	for to := range ctx.GetOutLinks()[from] {
		stake += getOverallLinkStake(ctx, from, to)
	}
	return stake
}

func getOverallInLinksStake(ctx *types.CalculationContext, from graphtypes.CidNumber) uint64 {
	stake := uint64(0)
	for to := range ctx.GetInLinks()[from] {
		stake += getOverallLinkStake(ctx, to, from) // reverse order here
	}
	return stake
}

func getNormalizedStake(ctx *types.CalculationContext, agent uint64) uint64 {
	return ctx.GetStakes()[agent] / ctx.GetNeudegs()[agent]
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

func entropyCalc(ctx *types.CalculationContext, entropy []float64, cidsCount int64, dampingFactor float64) {
	swd := make([]float64, cidsCount)
	sumswd := make([]float64, cidsCount)
	for i := range swd {
		swd[i] = dampingFactor*float64(
			getOverallInLinksStake(ctx, graphtypes.CidNumber(i))) + (1-dampingFactor)*float64(
			getOverallOutLinksStake(ctx, graphtypes.CidNumber(i)))
	}

	for i := range sumswd {
		for to := range ctx.GetInLinks()[graphtypes.CidNumber(i)] {
			sumswd[i] += dampingFactor * swd[to]
		}
		for to := range ctx.GetOutLinks()[graphtypes.CidNumber(i)] {
			sumswd[i] += (1 - dampingFactor) * swd[to]
		}
	}

	for i := range entropy {
		if swd[i] == 0 {
			continue
		}
		for to := range ctx.GetInLinks()[graphtypes.CidNumber(i)] {
			if sumswd[to] == 0 {
				continue
			}
			entropy[i] += math.Abs(-swd[i] / sumswd[to] * math.Log2(swd[i]/sumswd[to]))
		}
		for to := range ctx.GetOutLinks()[graphtypes.CidNumber(i)] {
			if sumswd[to] == 0 {
				continue
			}
			entropy[i] += math.Abs(-swd[i] / sumswd[to] * math.Log2(swd[i]/sumswd[to]))
		}
	}
}

func karmaCalc(ctx *types.CalculationContext, rank []float64, entropy []float64, karma []float64) {
	for from := range ctx.GetOutLinks() {
		stake := getOverallOutLinksStake(ctx, from)
		for to := range ctx.GetOutLinks()[from] {
			if stake == 0 {
				continue
			}
			users := ctx.GetOutLinks()[from][to]
			for user := range users {
				// if (ctx.GetStakes()[uint64(user)] == 0) { continue }
				if getNormalizedStake(ctx, uint64(user)) == 0 {
					continue
				}
				// w := float64(ctx.GetStakes()[uint64(user)]) / float64(stake)
				w := float64(getNormalizedStake(ctx, uint64(user))) / float64(stake)
				if math.IsNaN(w) {
					w = float64(0)
				}
				luminosity := rank[from] * entropy[from]
				karma[user] += w * float64(luminosity)
			}
		}
	}
}
