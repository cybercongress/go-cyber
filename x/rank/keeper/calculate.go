package keeper

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cybercongress/go-cyber/x/rank/types"

	"github.com/tendermint/tendermint/libs/log"

	"fmt"
	"time"
)

func CalculateRank(ctx *types.CalculationContext, unit types.ComputeUnit, logger log.Logger) (rank types.Rank) {
	start := time.Now()
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), "rank_calculation")

	if unit == types.CPU {
		//used only for development
		rank = types.NewRank(calculateRankCPU(ctx), logger, ctx.FullTree)
	} else {
		rank = types.NewRank(calculateRankGPU(ctx, logger), logger, ctx.FullTree)
	}

	diff := time.Since(start)

	logger.Info(
		"cyber~Rank calculated", "duration", diff.String(),
		"cyberlinks", ctx.LinksCount, "cids", ctx.CidsCount,
		"hash", fmt.Sprintf("%X", rank.MerkleTree.ExportSubtreesRoots()), // TODO remove this line before release
	)

	return
}

func CalculateRankInParallel(
	ctx *types.CalculationContext, rankChan chan types.Rank, err chan error, unit types.ComputeUnit, logger log.Logger,
) {
	defer func() {
		if r := recover(); r != nil {
			err <- r.(error)
		}
	}()

	rank := CalculateRank(ctx, unit, logger)
	rankChan <- rank
}
