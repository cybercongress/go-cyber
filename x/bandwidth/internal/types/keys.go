package types

const (
	// ModuleName is the name of the module
	ModuleName = "bandwidth"

	// default paramspace for params keeper
	DefaultParamspace = ModuleName

	// StoreKey is the store key string for bandwidth
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the bandwidth store.
	QuerierRoute = ModuleName

	// Query endpoints supported by the bandwidth querier
	QueryParameters         = "parameters"
	QueryDesirableBandwidth = "desirable_bandwidth"
	QueryMaxBlockBandwidth  = "max_block_bandwidth"
	QueryRecoveryPeriod     = "recovery_period"
	QueryAdjustPricePeriod  = "adjust_price_period"
	QueryBaseCreditPrice    = "base_credit_price"
	QueryTxCost             = "tx_cost"
	QueryLinkMsgCost        = "link_msg_cost"
	QueryNonLinkMsgCost     = "non_link_msg_cost"
)
