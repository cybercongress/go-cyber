// +build !cuda

package rank

import (
	. "github.com/cybercongress/cyberd/cosmos/poc/app/storage"
)

func calculateRankGPU(data *InMemoryStorage) ([]float64, int) {
	panic("Not Supported")
}
