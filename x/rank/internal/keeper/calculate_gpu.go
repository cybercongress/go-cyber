// +build cuda

package keeper

import (
	"github.com/cybercongress/cyberd/x/link"
	"github.com/cybercongress/cyberd/x/rank/internal/types"
	"github.com/tendermint/tendermint/libs/log"

	"math"
	"sync"
	"time"
)

/*
#cgo CFLAGS: -I/usr/lib/
#cgo LDFLAGS: -L/usr/local/cuda/lib64 -lcbdrank -lcudart
#include "cbdrank.h"
*/
import "C"

func calculateRankGPU(ctx *types.CalculationContext, logger log.Logger) []float64 {
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

	ch := make(chan int64, 100000)
	var wg sync.WaitGroup
	var lock1 sync.Mutex
	var lock2 sync.Mutex
	wg.Add(int(cidsCount))

	// the worker's function
	f := func(count int64) {
		defer wg.Done()
		if inLinks, sortedCids, ok := ctx.GetSortedInLinks(link.CidNumber(count)); ok {
			for _, cid := range sortedCids {
				inLinksCount[count] += uint32(len(inLinks[cid]))
				for acc := range inLinks[cid] {
					lock2.Lock()
					inLinksOuts = append(inLinksOuts, uint64(cid))
					inLinksUsers = append(inLinksUsers, uint64(acc))
					lock2.Unlock()
				}
			}
			linksCount += uint64(inLinksCount[count])
		}

		if outLinks, ok := outLinks[link.CidNumber(count)]; ok {
			for _, accs := range outLinks {
				outLinksCount[count] += uint32(len(accs))
				for acc := range accs {
					lock1.Lock()
					outLinksUsers = append(outLinksUsers, uint64(acc))
					lock1.Unlock()
				}
			}
		}
	}

	countWorkers := int64(math.Min(10000, float64(cidsCount)))

	// here the workers start
	for i:=int64(0); i < countWorkers; i++ {
		go func() {
			var count int64
			for {
				count = <- ch
				f(count)
			}
		}()
	}

	// data is added to the channel for workers
	for i := int64(0); i < cidsCount; i++ {
		ch <- i
	}

	// waiting for a while all workers will finish work
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

	logger.Info("Rank: data for gpu preparation", "time", time.Since(start))

	start = time.Now()
	cRank := (*C.double)(&rank[0])
	C.calculate_rank(
		cStakes, cStakesSize, cCidsSize, cLinksSize,
		cInLinksCount, cOutLinksCount,
		cInLinksOuts, cInLinksUsers, cOutLinksUsers,
		cRank,
	)
	logger.Info("Rank: gpu calculations", "time", time.Since(start))

	return rank
}
