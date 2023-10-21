package types

import (
	"sort"

	graphtypes "github.com/cybercongress/go-cyber/v2/x/graph/types"
)

type CalculationContext struct {
	CidsCount    int64
	LinksCount   int64
	NeuronsCount int64

	inLinks  map[graphtypes.CidNumber]graphtypes.CidLinks
	outLinks map[graphtypes.CidNumber]graphtypes.CidLinks

	stakes  map[uint64]uint64
	neudegs map[uint64]uint64

	FullTree bool

	DampingFactor float64
	Tolerance     float64
}

func NewCalcContext(
	linkIndex GraphIndexedKeeper, graphKeeper GraphKeeper,
	stakeKeeper StakeKeeper, fullTree bool, dampingFactor float64, tolerance float64,
	cidsCount, linksCount, neuronsCount uint64,
) *CalculationContext {
	return &CalculationContext{
		CidsCount:    int64(cidsCount),
		LinksCount:   int64(linksCount),
		NeuronsCount: int64(neuronsCount),

		inLinks:  linkIndex.GetInLinks(),
		outLinks: linkIndex.GetOutLinks(),

		stakes:  stakeKeeper.GetTotalStakesAmpere(),
		neudegs: graphKeeper.GetNeudegs(),

		FullTree: fullTree,

		DampingFactor: dampingFactor,
		Tolerance:     tolerance,
	}
}

func (calctx *CalculationContext) GetInLinks() map[graphtypes.CidNumber]graphtypes.CidLinks {
	return calctx.inLinks
}

func (calctx *CalculationContext) GetOutLinks() map[graphtypes.CidNumber]graphtypes.CidLinks {
	return calctx.outLinks
}

func (calctx *CalculationContext) GetCidsCount() int64 {
	return calctx.CidsCount
}

func (calctx *CalculationContext) GetNeuronsCount() int64 {
	return calctx.NeuronsCount
}

func (calctx *CalculationContext) GetStakes() map[uint64]uint64 {
	return calctx.stakes
}

func (calctx *CalculationContext) GetNeudegs() map[uint64]uint64 {
	return calctx.neudegs
}

func (calctx *CalculationContext) GetTolerance() float64 {
	return calctx.Tolerance
}

func (calctx *CalculationContext) GetDampingFactor() float64 {
	return calctx.DampingFactor
}

func (calctx *CalculationContext) GetSortedInLinks(cid graphtypes.CidNumber) (graphtypes.CidLinks, []graphtypes.CidNumber, bool) {
	links := calctx.inLinks[cid]

	if len(links) == 0 {
		return nil, nil, false
	}

	numbers := make([]graphtypes.CidNumber, 0, len(links))
	for num := range links {
		numbers = append(numbers, num)
	}

	sort.Slice(numbers, func(i, j int) bool { return numbers[i] < numbers[j] })

	return links, numbers, true
}

func (calctx *CalculationContext) GetSortedOutLinks(cid graphtypes.CidNumber) (graphtypes.CidLinks, []graphtypes.CidNumber, bool) {
	links := calctx.outLinks[cid]

	if len(links) == 0 {
		return nil, nil, false
	}

	numbers := make([]graphtypes.CidNumber, 0, len(links))
	for num := range links {
		numbers = append(numbers, num)
	}

	sort.Slice(numbers, func(i, j int) bool { return numbers[i] < numbers[j] })

	return links, numbers, true
}
