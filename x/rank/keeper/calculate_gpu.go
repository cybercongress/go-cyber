//go:build cuda
// +build cuda

package keeper

import (
	"time"

	"github.com/cometbft/cometbft/libs/log"

	graphtypes "github.com/cybercongress/go-cyber/v4/x/graph/types"
	"github.com/cybercongress/go-cyber/v4/x/rank/types"
)

/*
#cgo CFLAGS: -I/usr/lib/
#cgo LDFLAGS: -L/usr/local/cuda/lib64 -lcbdrank -lcudart
#include "cbdrank.h"
*/
import "C"

func calculateRankGPU(ctx *types.CalculationContext, logger log.Logger) types.EMState {
	start := time.Now()
	if ctx.GetCidsCount() == 0 {
		return types.EMState{
			RankValues:    make([]float64, 0),
			EntropyValues: make([]float64, 0),
			KarmaValues:   make([]float64, 0),
		}
	}

	tolerance := ctx.GetTolerance()
	dampingFactor := ctx.GetDampingFactor()

	cidsCount := ctx.GetCidsCount()
	stakesCount := ctx.GetNeuronsCount()
	linksCount := ctx.LinksCount

	rank := make([]float64, cidsCount)
	entropy := make([]float64, cidsCount)
	luminosity := make([]float64, cidsCount)
	karma := make([]float64, stakesCount)

	inLinksCount := make([]uint32, cidsCount)
	outLinksCount := make([]uint32, cidsCount)
	outLinksIns := make([]uint64, linksCount)
	inLinksOuts := make([]uint64, linksCount)
	inLinksUsers := make([]uint64, linksCount)
	outLinksUsers := make([]uint64, linksCount)

	// will fail if amount of indexed accounts will not equal all accounts
	// distribute current flow through all neuron's cyberlinks
	stakes := make([]uint64, stakesCount)
	for neuron, stake := range ctx.GetStakes() {
		neudeg := ctx.GetNeudegs()[neuron]
		if neudeg != 0 {
			stakes[neuron] = stake / neudeg
		} else {
			stakes[neuron] = 0
		}
	}

	var pointer1 uint32 = 0
	var pointer2 uint32 = 0
	for i := int64(0); i < cidsCount; i++ {

		if inLinks, sortedCids, ok := ctx.GetSortedInLinks(graphtypes.CidNumber(i)); ok {
			for _, cid := range sortedCids {
				inLinksCount[i] += uint32(len(inLinks[cid]))
				for acc := range inLinks[cid] {
					inLinksOuts[pointer1] = uint64(cid)
					inLinksUsers[pointer1] = uint64(acc)
					pointer1++
				}
			}
		}

		if outLinks, sortedCids, ok := ctx.GetSortedOutLinks(graphtypes.CidNumber(i)); ok {
			for _, cid := range sortedCids {
				outLinksCount[i] += uint32(len(outLinks[cid]))
				for acc := range outLinks[cid] {
					outLinksIns[pointer2] = uint64(cid)
					outLinksUsers[pointer2] = uint64(acc)
					pointer2++
				}
			}
		}
	}

	/* Convert to C types */
	cStakes := (*C.ulong)(&stakes[0])

	cStakesSize := C.ulong(len(stakes))
	cCidsSize := C.ulong(len(inLinksCount))
	cLinksSize := C.ulong(len(inLinksOuts))

	cInLinksCount := (*C.uint)(&inLinksCount[0])
	cOutLinksCount := (*C.uint)(&outLinksCount[0])

	cOutLinksIns := (*C.ulong)(&outLinksIns[0])
	cInLinksOuts := (*C.ulong)(&inLinksOuts[0])

	cInLinksUsers := (*C.ulong)(&inLinksUsers[0])
	cOutLinksUsers := (*C.ulong)(&outLinksUsers[0])

	cDampingFactor := C.double(dampingFactor)
	cTolerance := C.double(tolerance)

	logger.Info("Data transform", "duration", time.Since(start).String())

	start = time.Now()
	cRank := (*C.double)(&rank[0])
	cEntropy := (*C.double)(&entropy[0])
	cLuminosity := (*C.double)(&luminosity[0])
	cKarma := (*C.double)(&karma[0])
	C.calculate_rank(
		cStakes, cStakesSize, cCidsSize, cLinksSize,
		cInLinksCount, cOutLinksCount,
		cInLinksOuts, cOutLinksIns,
		cInLinksUsers, cOutLinksUsers,
		cDampingFactor, cTolerance,
		cRank, cEntropy, cLuminosity, cKarma,
	)
	logger.Info("Rank computation", "duration", time.Since(start).String())

	// return rank
	return types.EMState{
		RankValues:    rank,
		EntropyValues: entropy,
		KarmaValues:   karma,
	}
}
