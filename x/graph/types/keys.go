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
	NeudegTStoreKeyPrefix    = []byte{0x06} // inter-block cache for neurons cyberlink' degree

	LastCidNumber    		 = append(GlobalStoreKeyPrefix, []byte("lastParticleNumber")...)
	LinksCount 				 = append(GlobalStoreKeyPrefix, []byte("cyberlinksAmount")...)
	HasNewLinks 		     = append(GlobalStoreKeyPrefix, []byte("blockHasNewLinks")...)
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

func CyberlinksNewStoreKey(linkKey []byte) []byte {
	return append(CyberlinkStoreKeyPrefix, linkKey...)
}

func CyberlinkRawKey(link CompactLink) []byte {
	keyAsBytes := make([]byte, 24)
	copy(keyAsBytes[0:8],sdk.Uint64ToBigEndian(link.From))
	copy(keyAsBytes[8:16],sdk.Uint64ToBigEndian(link.Account))
	copy(keyAsBytes[16:24],sdk.Uint64ToBigEndian(link.To))
	return keyAsBytes
}
