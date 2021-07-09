package keeper

import (
	//"fmt"
	//"encoding/binary"
	"math"
	//"math/big"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
	"github.com/cybercongress/go-cyber/x/rank/types"
)


func calculateRankCPU(ctx *types.CalculationContext) types.EMState {

	inLinks := ctx.GetInLinks()
	tolerance := ctx.GetTolerance()
	dampingFactor := ctx.GetDampingFactor()

	// for cross debugging with GPU, remove before release
	// will panic if stakes don't have all accounts
	stakesCount := len(ctx.GetStakes())
	stakesTest := make([]uint64, stakesCount)
	for acc, stake := range ctx.GetStakes() {
		stakesTest[uint64(acc)] = stake
	}

	size := ctx.GetCidsCount()
	if size == 0  || len(ctx.GetStakes()) == 0 {
		return types.EMState{
			[]float64{},
			[]float64{},
			[]float64{},
			[]float64{},
		}
	}

	rank := make([]float64, size)
	entropy := make([]float64, size)
	luminosity := make([]float64, size)
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

	// TODO return sum to API after implementation in GPU
	_ = entropyCalc(ctx, entropy)
	_ = luminosityCalc(rank, entropy, luminosity)
	_ = karmaCalc(ctx, luminosity, karma)

	return types.EMState{
		rank,
		entropy,
		luminosity,
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
				weight := float64(linkStake) / float64(jCidOutStake)
				if math.IsNaN(weight) { weight = float64(0) }
				ksum = prevrank[j]*weight + ksum //force no-fma here by explicit conversion
			}
			rank[cid] = ksum*dampingFactor + defaultRankWithCorrection //force no-fma here by explicit conversion
		}
	}

	return rank
}

func getOverallLinkStake(ctx *types.CalculationContext, from graphtypes.CidNumber, to graphtypes.CidNumber) uint64 {

	stake := uint64(0)
	users := ctx.GetOutLinks()[from][to]
	for user := range users {
		stake += ctx.GetStakes()[uint64(user)]
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

func entropyCalc(ctx *types.CalculationContext, entropy []float64) (float64) {
	e := float64(0)

	for from := range ctx.GetOutLinks() {
		outStake := getOverallOutLinksStake(ctx, from)
		inStake := getOverallInLinksStake(ctx, from)
		ois := outStake + inStake
		for to := range ctx.GetOutLinks()[from] {
			users := ctx.GetOutLinks()[from][to]
			for user := range users {
				w := float64(ctx.GetStakes()[uint64(user)]) / float64(ois)
				if math.IsNaN(w) { w = float64(0) }
				e -= w*math.Log2(w)
				entropy[from] -= w*math.Log2(w)
			}
		}
	}

	return e
}

func luminosityCalc(rank, entropy, luminosity []float64) (float64) {
	l := float64(0)

	for i, _ := range rank {
		luminosity[i] = rank[i] * entropy[i]
		l += luminosity[i]
	}

	return l
}

func karmaCalc(ctx *types.CalculationContext, light []float64, karma []float64) (float64) {
	k := float64(0)

	for from := range ctx.GetOutLinks() {
		outStake := getOverallOutLinksStake(ctx, from)
		inStake := getOverallInLinksStake(ctx, from)
		ois := outStake + inStake
		for to := range ctx.GetOutLinks()[from] {
			users := ctx.GetOutLinks()[from][to]
			for user := range users {
				w := float64(ctx.GetStakes()[uint64(user)]) / float64(ois)
				if math.IsNaN(w) { w = float64(0) }
				karma[user] += w*float64(light[from])
				k += w*float64(light[from])
			}
		}
	}

	return k
}
