package types

import (
	"encoding/binary"

	"github.com/cybercongress/go-cyber/v5/types"
)

type CompactLink struct {
	From    uint64
	To      uint64
	Account uint64
}

type LinkFilter func(CompactLink) bool

func NewLink(from CidNumber, to CidNumber, acc types.AccNumber) CompactLink {
	return CompactLink{
		From:    uint64(from),
		To:      uint64(to),
		Account: uint64(acc),
	}
}

func UnmarshalBinaryLink(b []byte) CompactLink {
	return NewLink(
		CidNumber(binary.LittleEndian.Uint64(b[0:8])),
		CidNumber(binary.LittleEndian.Uint64(b[8:16])),
		types.AccNumber(binary.LittleEndian.Uint64(b[16:24])),
	)
}

func (l CompactLink) MarshalBinaryLink() []byte {
	b := make([]byte, 24)
	binary.LittleEndian.PutUint64(b[0:8], l.From)
	binary.LittleEndian.PutUint64(b[8:16], l.To)
	binary.LittleEndian.PutUint64(b[16:24], l.Account)
	return b
}
