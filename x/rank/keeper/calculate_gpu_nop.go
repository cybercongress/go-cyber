//go:build !cuda
// +build !cuda

package keeper

import (
	"github.com/cybercongress/go-cyber/v2/x/rank/types"

	"github.com/tendermint/tendermint/libs/log"
)

func calculateRankGPU(ctx *types.CalculationContext, logger log.Logger) types.EMState {
	panic("Daemon compiled without gpu support, but started in gpu mode")
}
