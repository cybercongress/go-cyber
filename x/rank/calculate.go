package rank

import (
	"github.com/tendermint/tendermint/libs/log"
)

type ComputeUnit int

const (
	CPU ComputeUnit = iota
	GPU ComputeUnit = iota
)

func CalculateRank(ctx *CalculationContext, unit ComputeUnit, logger log.Logger) ([]float64, int) {
	if unit == CPU {
		//used only for development
		return calculateRankCPU(ctx)
	} else {
		return calculateRankGPU(ctx, logger)
	}
}
