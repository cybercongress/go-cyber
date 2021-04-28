package types

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
)

type CalculationContext struct {
	CidsCount  int64
	LinksCount int64

	inLinks  map[graphtypes.CidNumber]graphtypes.CidLinks
	outLinks map[graphtypes.CidNumber]graphtypes.CidLinks

	stakes map[uint64]uint64

	FullTree bool

	DampingFactor float64
	Tolerance 	  float64
}

func NewCalcContext(
	ctx sdk.Context, linkIndex GraphIndexedKeeper, numberKeeper GraphKeeper,
	stakeKeeper StakeKeeper, fullTree bool, dampingFactor float64, tolerance float64,
	cidsCount, linksCount uint64,
) *CalculationContext {

	return &CalculationContext{

		CidsCount:  int64(cidsCount),
		LinksCount: int64(linksCount),

		inLinks:  linkIndex.GetInLinks(),
		outLinks: linkIndex.GetOutLinks(),

		stakes: stakeKeeper.GetTotalStakesAmpere(),

		FullTree: fullTree,

		DampingFactor: dampingFactor,
		Tolerance: tolerance,
	}
}

func (c *CalculationContext) GetInLinks() map[graphtypes.CidNumber]graphtypes.CidLinks {
	return c.inLinks
}

func (c *CalculationContext) GetOutLinks() map[graphtypes.CidNumber]graphtypes.CidLinks {
	return c.outLinks
}

func (c *CalculationContext) GetCidsCount() int64 {
	return c.CidsCount
}

func (c *CalculationContext) GetStakes() map[uint64]uint64 {
	return c.stakes
}

func (c *CalculationContext) GetTolerance() float64 {
	return c.Tolerance
}

func (c *CalculationContext) GetDampingFactor() float64 {
	return c.DampingFactor
}

func (с *CalculationContext) GetSortedInLinks(cid graphtypes.CidNumber) (graphtypes.CidLinks, []graphtypes.CidNumber, bool) {
	links := с.inLinks[cid]

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
