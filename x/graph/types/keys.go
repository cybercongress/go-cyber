package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName 	  		= "graph"
	RouterKey 	  		= ModuleName
	StoreKey      	    = ModuleName
	QuerierRoute  		= ModuleName

	TStoreKey           = "transient_index"
)

var (
	GlobalStoreKeyPrefix  	 = []byte{0x00}
	CidStoreKeyPrefix	  	 = []byte{0x01}
	CidReverseStoreKeyPrefix = []byte{0x02}
	CyberlinkStoreKeyPrefix  = []byte{0x03}
	CyberlinkTStoreKeyPrefix = []byte{0x04} // inter-block cache for cyberlinks
	NeudegStoreKeyPrefix     = []byte{0x05}
	NeudegTStoreKeyPrefix     = []byte{0x06}

	LastCidNumber    		 = append(GlobalStoreKeyPrefix, []byte("lastCidNumber")...)
	LinksCount 				 = append(GlobalStoreKeyPrefix, []byte("linksCount")...)
	HasNewLinks 		     = append(GlobalStoreKeyPrefix, []byte("hasNewLinks")...)
)

func CidStoreKey(cid Cid) []byte {
	return append(CidStoreKeyPrefix, []byte(cid)...)
}

func CidReverseStoreKey(num CidNumber) []byte {
	return append(CidReverseStoreKeyPrefix, sdk.Uint64ToBigEndian(uint64(num))...)
}

func CyberlinksStoreKey(id uint64) []byte {
	return append(CyberlinkStoreKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func CyberlinksTStoreKey(link []byte) []byte {
	return append(CyberlinkTStoreKeyPrefix, link...)
}

func NeudegStoreKey(accNumber uint64) []byte {
	return append(NeudegStoreKeyPrefix, sdk.Uint64ToBigEndian(accNumber)...)
}

func NeudegTStoreKey(accNumber uint64) []byte {
	return append(NeudegTStoreKeyPrefix, sdk.Uint64ToBigEndian(accNumber)...)
}
