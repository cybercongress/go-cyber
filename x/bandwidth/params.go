package bandwidth

const (
	// Maximum bandwidth of network
	MaxNetworkBandwidth = 100000000

	// todo add more msg prices here
	// Bandwidth cost of specific messages and tx itself
	TxCost         int64 = 10
	LinkMsgCost    int64 = 5
	NonLinkMsgCost int64 = 1

	// Number of blocks to recover full bandwidth
	RecoveryPeriod = 1000

	AdjustPricePeriod = 10

	BaseCreditPrice = 1
)
