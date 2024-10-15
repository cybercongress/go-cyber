package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "bandwidth"
	StoreKey   = ModuleName
	TStoreKey  = "transient_bandwidth"
)

var (
	GlobalStoreKeyPrefix  = []byte{0x00}
	AccountStoreKeyPrefix = []byte{0x01}
	BlockStoreKeyPrefix   = []byte{0x02}
	ParamsKey             = []byte{0x02}

	LastBandwidthPrice = append(GlobalStoreKeyPrefix, []byte("lastBandwidthPrice")...)
	TotalBandwidth     = append(GlobalStoreKeyPrefix, []byte("totalBandwidth")...)
	BlockBandwidth     = append(GlobalStoreKeyPrefix, []byte("blockBandwidth")...)
)

func AccountStoreKey(addr string) []byte {
	return append(AccountStoreKeyPrefix, []byte(addr)...)
}

func BlockStoreKey(blockNumber uint64) []byte {
	return append(BlockStoreKeyPrefix, sdk.Uint64ToBigEndian(blockNumber)...)
}
