package types

import . "github.com/cybercongress/cyberd/types"

type Link struct {
	From Cid
	To   Cid
}

type CompactLink struct {
	from CidNumber
	to   CidNumber
	acc  AccNumber
}

func NewLink(from CidNumber, to CidNumber, acc AccNumber) CompactLink {
	return CompactLink{from: from, to: to, acc: acc}
}

func (l CompactLink) From() CidNumber { return l.from }
func (l CompactLink) To() CidNumber   { return l.to }
func (l CompactLink) Acc() AccNumber  { return l.acc }
