package bandwidth

const (
	// Maximum bandwidth of network
	DesirableNetworkBandwidthForRecoveryPeriod = LinkMsgCost * 1500 * 1000

	// Bandwidth cost of specific messages and tx itself
	LinkMsgCost    int64 = 100
	TxCost               = LinkMsgCost * 3
	NonLinkMsgCost       = LinkMsgCost * 100

	// Number of blocks to recover full bandwidth
	RecoveryPeriod = 60 * 60 * 24
	// Number of blocks before next adjust price
	AdjustPricePeriod = 60 * 10
	BaseCreditPrice   = 1
)
