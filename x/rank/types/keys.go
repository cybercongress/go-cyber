package types

const (
	ModuleName 	     = "rank"
	StoreKey         = ModuleName
	RouterKey        = ModuleName
	QuerierRoute     = ModuleName
)

var (
	GlobalStoreKeyPrefix  	 = []byte{0x00}

	LatestBlockNumber 		 = append(GlobalStoreKeyPrefix, []byte("latestBlockNumber")...)
	LatestMerkleTree 		 = append(GlobalStoreKeyPrefix, []byte("latestMerkleTree")...)
	NextMerkleTree 		     = append(GlobalStoreKeyPrefix, []byte("nextMerkleTree")...)
	//RankCalcFinished 		 = append(GlobalStoreKeyPrefix, []byte("rankCalcFinished")...)
	NextRankCidCount 		 = append(GlobalStoreKeyPrefix, []byte("nextRankCidCount")...)
	ContextCidCount 		 = append(GlobalStoreKeyPrefix, []byte("contextCidCount")...)
	ContextLinkCount 		 = append(GlobalStoreKeyPrefix, []byte("contextLinkCount")...)
)

