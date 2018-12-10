package bandwidth

const (
	// Maximum bandwidth of network
	MaxNetworkBandwidth = 100000000

	// Bandwidth cost of specific messages and tx itself
	TxCost         int64 = 10
	LinkMsgCost    int64 = 5
	NonLinkMsgCost int64 = 1

	// Number of blocks to recover full bandwidth
	RecoveryPeriod = 1000

	AdjustPriceInterval = 10

	BaseCreditPrice = 1
)
