package types

type AccountNumber uint64

type Link struct {
	from CidNumber
	to   CidNumber
	acc  AccountNumber
}

func NewLink(from CidNumber, to CidNumber, acc AccountNumber) Link {
	return Link{from: from, to: to, acc: acc}
}

func (l Link) From() CidNumber    { return l.from }
func (l Link) To() CidNumber      { return l.to }
func (l Link) Acc() AccountNumber { return l.acc }
