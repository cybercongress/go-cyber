package rank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/cybercongress/cyberd/types"
	"github.com/cybercongress/cyberd/x/bank"
	"github.com/cybercongress/cyberd/x/link/keeper"
	. "github.com/cybercongress/cyberd/x/link/types"
	"sort"
)

type CalculationContext struct {
	cidsCount int64

	inLinks  map[CidNumber]CidLinks
	outLinks map[CidNumber]CidLinks

	stakes map[AccNumber]uint64
}

func NewCalcContext(
	ctx sdk.Context, linkIndex keeper.LinkIndexedKeeper,
	numberKeeper keeper.CidNumberKeeper, indexedKeeper bank.IndexedKeeper,
) *CalculationContext {

	return &CalculationContext{
		cidsCount: int64(numberKeeper.GetCidsCount(ctx)),

		inLinks:  linkIndex.GetInLinks(),
		outLinks: linkIndex.GetOutLinks(),

		stakes: indexedKeeper.GetStakes(),
	}
}

func (c *CalculationContext) GetInLinks() map[CidNumber]CidLinks {
	return c.inLinks
}

func (c *CalculationContext) GetOutLinks() map[CidNumber]CidLinks {
	return c.outLinks
}

func (c *CalculationContext) GetCidsCount() int64 {
	return c.cidsCount
}

func (c *CalculationContext) GetStakes() map[AccNumber]uint64 {
	return c.stakes
}

func (с *CalculationContext) GetSortedInLinks(cid CidNumber) (CidLinks, []CidNumber, bool) {
	links := с.inLinks[cid]

	if len(links) == 0 {
		return nil, nil, false
	}

	numbers := make([]CidNumber, 0, len(links))
	for num := range links {
		numbers = append(numbers, num)
	}

	sort.Slice(numbers, func(i, j int) bool { return numbers[i] < numbers[j] })

	return links, numbers, true
}
