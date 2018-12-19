// +build !cuda

package rank

import (
	"github.com/tendermint/tendermint/libs/log"
)

func calculateRankGPU(ctx *CalculationContext, logger log.Logger) ([]float64, int) {
	panic("Daemon compiled without gpu support, but started in gpu mode")
}
