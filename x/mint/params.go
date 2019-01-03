package mint

const (
	// 31_536_000
	// assuming 1 second block times
	BlocksPerYear = 60 * 60 * 24 * 365
	// todo hardcoded
	GenesisSupply = int64(10 * 1000 * 1000 * 1000 * 1000 * 1000) // 1*10^15
	// percentage value 0....1, 1 means 100% year inflation
	InflationRatePerYear = 2
)
