package bandwidth

const (
	// Maximum bandwidth of network
	DesirableNetworkBandwidthForRecoveryPeriod = 1000000000

	// todo add more msg prices here
	// Bandwidth cost of specific messages and tx itself
	TxCost         int64 = 10
	LinkMsgCost    int64 = 5
	NonLinkMsgCost int64 = 1

	// Number of blocks to recover full bandwidth
	RecoveryPeriod = 60 * 60 * 24
	// Number of blocks before next adjust price
	AdjustPricePeriod = 60 * 10
	BaseCreditPrice   = 1
)
