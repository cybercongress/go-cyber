// +build !cuda

package rank

import (
	. "github.com/cybercongress/cyberd/cosmos/poc/app/storage"
	"github.com/tendermint/tendermint/libs/log"
)

func calculateRankGPU(data *InMemoryStorage, logger log.Logger) ([]float64, int) {
	panic("Not Supported")
}
