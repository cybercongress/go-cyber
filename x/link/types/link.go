package types

import . "github.com/cybercongress/cyberd/types"

type Link struct {
	from CidNumber
	to   CidNumber
	acc  AccNumber
}

func NewLink(from CidNumber, to CidNumber, acc AccNumber) Link {
	return Link{from: from, to: to, acc: acc}
}

func (l Link) From() CidNumber { return l.from }
func (l Link) To() CidNumber   { return l.to }
func (l Link) Acc() AccNumber  { return l.acc }
