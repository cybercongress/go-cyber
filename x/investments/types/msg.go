package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	ctypes "github.com/cybercongress/go-cyber/types"
	//"github.com/cybercongress/go-cyber/x/investment/types"
)

func NewMsgInvest(
	investor sdk.AccAddress,
	amount sdk.Coin,
	resource string,
	endtime int64,
) *MsgInvest {
	return &MsgInvest{
		Investor: investor.String(),
		Amount:   amount,
		Resource: resource,
		EndTime:  endtime,
	}
}

func (msg MsgInvest) Route() string { return RouterKey }

func (msg MsgInvest) Type() string { return ActionInvest }

func (msg MsgInvest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Investor); if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid investor address: %s", err)
	}

	if msg.Amount.Denom != ctypes.CYB {
		return sdkerrors.Wrap(ErrInvalidResource, msg.Resource)
	}

	if msg.Resource != ctypes.VOLT && msg.Resource != ctypes.AMPER {
		return sdkerrors.Wrap(ErrInvalidResource, msg.Resource)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if msg.EndTime <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid end time")
	}

	return nil
}

func (msg MsgInvest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgInvest) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Investor)
	return []sdk.AccAddress{addr}
}

//------------------------------

func NewMsgCreateResource(
	sender sdk.AccAddress,
	resource sdk.Coin,
	receiver sdk.AccAddress,
) *MsgCreateResource {
	return &MsgCreateResource{
		Sender: sender.String(),
		Resource: resource,
		Receiver: receiver.String(),
	}
}

func (msg MsgCreateResource) Route() string { return RouterKey }

func (msg MsgCreateResource) Type() string { return ActionCreateResource }

func (msg MsgCreateResource) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender); if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Receiver); if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address: %s", err)
	}

	if msg.Resource.IsZero() || !msg.Resource.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid resource %s", msg.Resource)
	}

	return nil
}

func (msg MsgCreateResource) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCreateResource) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

//------------------------------

func NewMsgRedeemResource(
	sender sdk.AccAddress,
	resource sdk.Coin,
	receiver sdk.AccAddress,
) *MsgRedeemResource {
	return &MsgRedeemResource{
		Sender: sender.String(),
		Resource: resource,
	}
}

func (msg MsgRedeemResource) Route() string { return RouterKey }

func (msg MsgRedeemResource) Type() string { return ActionRedeemResource }

func (msg MsgRedeemResource) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender); if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}

	if msg.Resource.IsZero() || !msg.Resource.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid resource %s", msg.Resource)
	}

	return nil
}

func (msg MsgRedeemResource) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgRedeemResource) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

