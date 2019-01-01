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

func CalculateRank(ctx *CalculationContext, unit ComputeUnit, logger log.Logger) (rank Rank) {
	start := time.Now()
	if unit == CPU {
		//used only for development
		rank = NewRank(calculateRankCPU(ctx), ctx.fullTree)
	} else {
		rank = NewRank(calculateRankGPU(ctx, logger), ctx.fullTree)
	}
	logger.Info(
		"Rank calculated", "time", time.Since(start), "links", ctx.linksCount, "cids", ctx.cidsCount,
	)
	return
}

func CalculateRankInParallel(
	ctx *CalculationContext, rankChan chan Rank, err chan error, unit ComputeUnit, logger log.Logger,
) {
	defer func() {
		if r := recover(); r != nil {
			err <- r.(error)
		}
	}()

	rank := CalculateRank(ctx, unit, logger)
	rankChan <- rank
}
