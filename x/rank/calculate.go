package rank

import (
	"github.com/tendermint/tendermint/libs/log"
	"time"
)

type ComputeUnit int

const (
	CPU ComputeUnit = iota
	GPU ComputeUnit = iota
)

func CalculateRank(ctx *CalculationContext, unit ComputeUnit, logger log.Logger) (rank []float64) {
	start := time.Now()
	if unit == CPU {
		//used only for development
		rank = calculateRankCPU(ctx)

	} else {
		rank = calculateRankGPU(ctx, logger)
	}
	logger.Info(
		"Rank calculated", "time", time.Since(start), "links", ctx.linksCount, "cids", ctx.cidsCount,
	)
	return
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
