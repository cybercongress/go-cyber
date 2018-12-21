package rank

import (
	"github.com/tendermint/tendermint/libs/log"
)

type ComputeUnit int

const (
	CPU ComputeUnit = iota
	GPU ComputeUnit = iota
)

func CalculateRank(ctx *CalculationContext, unit ComputeUnit, logger log.Logger) []float64 {
	if unit == CPU {
		//used only for development
		return calculateRankCPU(ctx)
	} else {
		return calculateRankGPU(ctx, logger)
	}
}

func CalculateRankInParallel(
	ctx *CalculationContext, rankChan chan []float64, err chan error, unit ComputeUnit, logger log.Logger,
) {
	defer func() {
		if r := recover(); r != nil {
			err <- r.(error)
		}
	}()

	rank := CalculateRank(ctx, unit, logger)
	rankChan <- rank
}
