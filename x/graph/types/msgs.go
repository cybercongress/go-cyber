package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ipfs/go-cid"
)

const (
	TypeMsgCyberlink    = "cyberlink"
)

var (
	_ sdk.Msg = &MsgCyberlink{}
)

func NewMsgCyberlink(address sdk.AccAddress, links []Link) *MsgCyberlink {
	return &MsgCyberlink{
		Neuron: address.String(),
		Links:   links,
	}
}

func (msg MsgCyberlink) Route() string { return RouterKey }

func (msg MsgCyberlink) Type()  string { return TypeMsgCyberlink }

func (msg MsgCyberlink) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Neuron)
	return []sdk.AccAddress{addr}
}

func (msg MsgCyberlink) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCyberlink) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(msg.Neuron)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid neuron address: %s", err)
	}

	if len(msg.Links) == 0 {
		return ErrZeroLinks
	}

	var filter = make(CidsFilter)
	for _, link := range msg.Links {
		if link.From == link.To {
			return ErrSelfLink
		}

		fromCid, err := cid.Decode(link.From)
		if err != nil {
			return ErrInvalidParticle
		}

		if fromCid.Version() != 0 {
			return ErrCidVersion
		}

		toCid, err := cid.Decode(link.To)
		if err != nil {
			return ErrInvalidParticle
		}

		if toCid.Version() != 0 {
			return ErrCidVersion
		}

		if filter.Contains(Cid(link.From), Cid(link.To)) {
			return ErrCyberlinkExist
		}

		filter.Put(Cid(link.From), Cid(link.To))
	}

	return nil
}