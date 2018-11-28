package rank

import (
	. "github.com/cybercongress/cyberd/cosmos/poc/app/storage"
	"github.com/tendermint/tendermint/libs/log"
)

type ComputeUnit int

const (
	CPU ComputeUnit = iota
	GPU ComputeUnit = iota
)

func CalculateRank(data *InMemoryStorage, unit ComputeUnit, logger log.Logger) ([]float64, int) {
	if unit == CPU {
		//used only for development
		return calculateRankCPU(data)
	} else {
		return calculateRankGPU(data, logger)
	}
}
