// +build cuda

package keeper

import (
	"github.com/tendermint/tendermint/libs/log"

	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
	"github.com/cybercongress/go-cyber/x/rank/types"

	"time"
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
			RankValues:       make([]float64, 0),
			EntropyValues:    make([]float64, 0),
			KarmaValues:      make([]float64, 0),
		}
	}

	tolerance := ctx.GetTolerance()
	dampingFactor := ctx.GetDampingFactor()

	cidsCount := ctx.GetCidsCount()
	stakesCount := len(ctx.GetStakes())

	rank := make([]float64, cidsCount)
	entropy := make([]float64, cidsCount)
	luminosity := make([]float64, cidsCount)
	karma := make([]float64, stakesCount)

	inLinksCount := make([]uint32, cidsCount)
	outLinksCount := make([]uint32, cidsCount)
	outLinksIns := make([]uint64, 0)
	inLinksOuts := make([]uint64, 0)
	inLinksUsers := make([]uint64, 0)
	outLinksUsers := make([]uint64, 0)
	// will fail if amount of indexed accounts will not equal all accounts
	// distribute current flow through all neuron's cyberlinks
	stakes := make([]uint64, ctx.GetNeuronsCount())
	for neuron, stake := range ctx.GetStakes() {
		neudeg := ctx.GetNeudegs()[neuron]
		if neudeg != 0 {
			stakes[neuron] = stake / neudeg
		} else {
			stakes[neuron] = 0
		}
	}

	for i := int64(0); i < cidsCount; i++ {

		if inLinks, sortedCids, ok := ctx.GetSortedInLinks(graphtypes.CidNumber(i)); ok {
			for _, cid := range sortedCids {
				inLinksCount[i]  += uint32(len(inLinks[cid]))
				for acc := range inLinks[cid] {
					inLinksOuts = append(inLinksOuts, uint64(cid))
					inLinksUsers = append(inLinksUsers, uint64(acc))
				}
			}
		}

		if outLinks, sortedCids, ok := ctx.GetSortedOutLinks(graphtypes.CidNumber(i)); ok {
			for _, cid := range sortedCids {
				outLinksCount[i] += uint32(len(outLinks[cid]))
				for acc := range outLinks[cid] {
					outLinksIns = append(outLinksIns, uint64(cid))
					outLinksUsers = append(outLinksUsers, uint64(acc))
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
		cInLinksOuts,  cOutLinksIns,
		cInLinksUsers, cOutLinksUsers,
		cDampingFactor, cTolerance,
		cRank, cEntropy, cLuminosity, cKarma,
	)
	logger.Info("Rank computation", "duration", time.Since(start).String())

	//return rank
	return types.EMState{
		RankValues:       rank,
		EntropyValues:    entropy,
		KarmaValues:      karma,
	}
}
