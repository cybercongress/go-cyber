package rank

import (
	"github.com/tendermint/tendermint/libs/log"
)

type ComputeUnit int

const (
	CPU ComputeUnit = iota
	GPU ComputeUnit = iota
)

func CalculateRank(ctx *CalculationContext, rankChan chan []float64, unit ComputeUnit, logger log.Logger) {
	if unit == CPU {
		//used only for development
		calculateRankCPU(ctx, rankChan)
	} else {
		calculateRankGPU(ctx, rankChan, logger)
	}
}
