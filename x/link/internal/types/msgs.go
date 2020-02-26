package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ipfs/go-cid"
)

type Msg struct {
	Address sdk.AccAddress `json:"address"`
	Links   []Link         `json:"links"`
}

var _ sdk.Msg = Msg{}

func NewMsg(address sdk.AccAddress, links []Link) Msg {
	return Msg{
		Address: address,
		Links:   links,
	}
}

func (msg Msg) Name() string { return "link" }

func (Msg) Route() string { return "link" }
func (Msg) Type() string  { return "link" }

func (msg Msg) ValidateBasic() error {

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

func (msg Msg) GetSignBytes() []byte {

	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg Msg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}
