package rank

import . "github.com/cybercongress/cyberd/cosmos/poc/app/storage"

type ComputeUnit int

const (
	CPU ComputeUnit = iota
	GPU ComputeUnit = iota
)

func CalculateRank(data *InMemoryStorage, unit ComputeUnit) ([]float64, int) {
	if unit == CPU {
		return calculateRankCPU(data)
	} else {
		return calculateRankGPU(data)
	}
}
