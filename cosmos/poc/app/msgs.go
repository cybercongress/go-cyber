package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/cosmos/poc/app/storage"
)

type MsgLink struct {
	Address sdk.AccAddress `json:"address"`
	CidFrom storage.Cid    `json:"cid1"`
	CidTo   storage.Cid    `json:"cid2"`
}

var _ sdk.Msg = MsgLink{}

func NewMsgLink(address sdk.AccAddress, fromCid storage.Cid, toCid storage.Cid) MsgLink {
	return MsgLink{Address: address, CidFrom: fromCid, CidTo: toCid}
}

func (msg MsgLink) Name() string {
	return "link"
}

func (MsgLink) Route() string { return "link" }
func (MsgLink) Type() string { return "link" }

func (msg MsgLink) ValidateBasic() sdk.Error {

	if len(msg.Address) == 0 {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}

	if len(msg.CidFrom) == 0 || len(msg.CidTo) == 0 {
		return ErrInvalidCid(DefaultCodespace).TraceSDK("")
	}

	return nil
}

func (msg MsgLink) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgLink) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}
