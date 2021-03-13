package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ipfs/go-cid"
)

func NewMsgCyberlink(address sdk.AccAddress, links []Link) *MsgCyberlink {
	return &MsgCyberlink{
		Address: address.String(),
		Links:   links,
	}
}

func (MsgCyberlink) Route() string { return RouterKey }

func (MsgCyberlink) Type()  string { return ActionCyberlink }

func (msg MsgCyberlink) ValidateBasic() error {

	if len(msg.Address) == 0 {
		return sdkerrors.ErrInvalidAddress
	}

	if len(msg.Links) == 0 {
		return ErrZeroLinks
	}

	var filter = make(CidsFilter)
	for _, link := range msg.Links {
		if link.From == link.To {
			return ErrSelfLink
		}

		if _, err := cid.Decode(link.From); err != nil {
			return ErrInvalidCid
		}

		if _, err := cid.Decode(link.To); err != nil {
			return ErrInvalidCid
		}

		if filter.Contains(Cid(link.From), Cid(link.To)) {
			return ErrCyberlinkExist
		}

		filter.Put(Cid(link.From), Cid(link.To))
	}

	return nil
}

func (msg MsgCyberlink) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCyberlink) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Address)
	return []sdk.AccAddress{addr}
}