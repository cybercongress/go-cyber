// +build !cuda

package rank

import (
	"github.com/tendermint/tendermint/libs/log"
)

func calculateRankGPU(ctx *CalculationContext, rankChan chan []float64, logger log.Logger) {
	panic("Daemon compiled without gpu support, but started in gpu mode")
}
