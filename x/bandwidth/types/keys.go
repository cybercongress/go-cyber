package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName 			= "bandwidth"

	StoreKey			= ModuleName

	QuerierRoute 		= ModuleName
)

var (
	GlobalStoreKeyPrefix  = []byte{0x00}

	AccountStoreKeyPrefix = []byte{0x01}

	BlockStoreKeyPrefix   = []byte{0x02}

	LastBandwidthPrice    = append(GlobalStoreKeyPrefix, []byte("lastBandwidthPrice")...)
)

func AccountStoreKey(addr sdk.AccAddress) []byte {
	return append(AccountStoreKeyPrefix, addr.Bytes()...)
}

func BlockStoreKey(blockNumber uint64) []byte {
	return append(BlockStoreKeyPrefix, sdk.Uint64ToBigEndian(blockNumber)...)
}