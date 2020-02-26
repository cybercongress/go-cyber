package types

const (
	// ModuleName is the name of the module
	ModuleName = "rank"

	DefaultParamspace = ModuleName

	// StoreKey is the store key string for bandwidth
	StoreKey = ModuleName

	// RouterKey is the message route for bandwidth
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = ModuleName

	// Query endpoints supported by the minting querier
	QueryParameters         = "parameters"
	QueryCalculationWindow  = "calculation_window"
	QueryDampingFactor      = "damping_factor"
	QueryTolerance          = "tolerance"
)
