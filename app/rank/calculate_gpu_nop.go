// +build !cuda

package rank

import (
	. "github.com/cybercongress/cyberd/app/storage"
	"github.com/tendermint/tendermint/libs/log"
)

func calculateRankGPU(data *InMemoryStorage, logger log.Logger) ([]float64, int) {
	panic("Daemon compiled without gpu support, but started in gpu mode")
}
