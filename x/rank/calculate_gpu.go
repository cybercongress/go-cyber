// +build cuda

package rank

import (
	"github.com/cybercongress/cyberd/x/link/types"
	"github.com/tendermint/tendermint/libs/log"
	"sync"
	"time"
)

/*
#cgo CFLAGS: -I/usr/lib/
#cgo LDFLAGS: -L/usr/local/cuda/lib64 -lcbdrank -lcudart
#include "cbdrank.h"
*/
import "C"

func calculateRankGPU(ctx *CalculationContext, logger log.Logger) []float64 {
	start := time.Now()
	if ctx.GetCidsCount() == 0 {
		return make([]float64, 0)
	}

	outLinks := ctx.GetOutLinks()

	cidsCount := ctx.GetCidsCount()
	linksCount := uint64(0)
	stakesCount := len(ctx.GetStakes())

	rank := make([]float64, cidsCount)
	inLinksCount := make([]uint32, cidsCount)
	outLinksCount := make([]uint32, cidsCount)

	inLinksOuts := make([]uint64, 0)
	inLinksUsers := make([]uint64, 0)
	outLinksUsers := make([]uint64, 0)

	// todo reduce size of stake by passing only participating in linking stakes.
	stakes := make([]uint64, stakesCount)
	for acc, stake := range ctx.GetStakes() {
		stakes[uint64(acc)] = stake
	}

	var wg sync.WaitGroup
	wg.Add(int(cidsCount))

	for i := int64(0); i < cidsCount; i++ {
		/* Fill values */
		go func(count int64) {
			defer wg.Done()
			if inLinks, sortedCids, ok := ctx.GetSortedInLinks(types.CidNumber(i)); ok {
				for _, cid := range sortedCids {
					inLinksCount[count] += uint32(len(inLinks[cid]))
					for acc := range inLinks[cid] {
						inLinksOuts = append(inLinksOuts, uint64(cid))
						inLinksUsers = append(inLinksUsers, uint64(acc))
					}
				}
				linksCount += uint64(inLinksCount[count])
			}

			if outLinks, ok := outLinks[types.CidNumber(count)]; ok {
				for _, accs := range outLinks {
					outLinksCount[count] += uint32(len(accs))
					for acc := range accs {
						outLinksUsers = append(outLinksUsers, uint64(acc))
					}
				}
			}
		}(i)
	}

	wg.Wait()

	/* Convert to C types */
	cStakes := (*C.ulong)(&stakes[0])

	cStakesSize := C.ulong(len(stakes))
	cCidsSize := C.ulong(len(inLinksCount))
	cLinksSize := C.ulong(len(inLinksOuts))

	cInLinksCount := (*C.uint)(&inLinksCount[0])
	cOutLinksCount := (*C.uint)(&outLinksCount[0])

	cInLinksOuts := (*C.ulong)(&inLinksOuts[0])
	cInLinksUsers := (*C.ulong)(&inLinksUsers[0])
	cOutLinksUsers := (*C.ulong)(&outLinksUsers[0])

	logger.Debug("Rank: data for gpu preparation", "time", time.Since(start))

	start = time.Now()
	cRank := (*C.double)(&rank[0])
	C.calculate_rank(
		cStakes, cStakesSize, cCidsSize, cLinksSize,
		cInLinksCount, cOutLinksCount,
		cInLinksOuts, cInLinksUsers, cOutLinksUsers,
		cRank,
	)
	logger.Debug("Rank: gpu calculations", "time", time.Since(start))

	return rank
}
