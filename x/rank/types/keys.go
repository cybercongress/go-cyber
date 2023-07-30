package types

const (
	ModuleName   = "rank"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	GlobalStoreKeyPrefix = []byte{0x00}

	LatestBlockNumber = append(GlobalStoreKeyPrefix, []byte("latestBlockNumber")...)
	LatestMerkleTree  = append(GlobalStoreKeyPrefix, []byte("latestMerkleTree")...)
	NextMerkleTree    = append(GlobalStoreKeyPrefix, []byte("nextMerkleTree")...)
	NextRankCidCount  = append(GlobalStoreKeyPrefix, []byte("nextRankParticlesAmount")...)
	ContextCidCount   = append(GlobalStoreKeyPrefix, []byte("contextParticlesAmount")...)
	ContextLinkCount  = append(GlobalStoreKeyPrefix, []byte("contextLinkAmount")...)
)
