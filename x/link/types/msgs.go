package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ipfs/go-cid"
)

type MsgCyberlink struct {
	Address sdk.AccAddress `json:"address"`
	Links   []Link         `json:"links"`
}

var _ sdk.Msg = MsgCyberlink{}

func NewMsgCyberlink(address sdk.AccAddress, links []Link) MsgCyberlink {
	return MsgCyberlink{
		Address: address,
		Links:   links,
	}
}

func (msg MsgCyberlink) Name() string { return "link" }
func (MsgCyberlink) Route()    string { return "link" }
func (MsgCyberlink) Type()     string { return "link" }

func (msg MsgCyberlink) ValidateBasic() error {

	if len(msg.Address) == 0 {
		return sdkerrors.ErrInvalidAddress
	}

	if len(msg.Links) == 0 {
		return ErrZeroLinks
	}

	var filter = make(CidsFilter)
	for _, link := range msg.Links {

		if _, err := cid.Decode(string(link.From)); err != nil {
			return ErrInvalidCid
		}

		if _, err := cid.Decode(string(link.To)); err != nil {
			return ErrInvalidCid
		}

		if filter.Contains(link.From, link.To) {
			return ErrCyberlinkExist
		}

		filter.Put(link.From, link.To)
	}

	return nil
}

func (msg MsgCyberlink) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgCyberlink) GetSignBytes() []byte {
	bz := msgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
