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

var (
	GlobalStoreKeyPrefix  	 = []byte{0x00}

	LatestBlockNumber 		 = append(GlobalStoreKeyPrefix, []byte("latestBlockNumber")...)
	LatestMerkleTree 		 = append(GlobalStoreKeyPrefix, []byte("latestMerkleTree")...)
	NextMerkleTree 		     = append(GlobalStoreKeyPrefix, []byte("nextMerkleTree")...)
	RankCalcFinished 		 = append(GlobalStoreKeyPrefix, []byte("rankCalcFinished")...)
	NextRankCidCount 		 = append(GlobalStoreKeyPrefix, []byte("nextRankCidCount")...)
)

