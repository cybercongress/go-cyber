package storage

// All addition storages introduced by cyberd
type CyberdPersistentStorages struct {
	CidIndex CidIndexStorage
	Links    LinksStorage
	Rank     RankStorage
}
