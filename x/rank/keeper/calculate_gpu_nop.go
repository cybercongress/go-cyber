//go:build !cuda
// +build !cuda

package keeper

import (
	"github.com/cybercongress/go-cyber/v6/x/rank/types"

	"github.com/cometbft/cometbft/libs/log"
)

func calculateRankGPU(ctx *types.CalculationContext, logger log.Logger) types.EMState {
	panic("Daemon compiled without gpu support, but started in gpu mode")
}
