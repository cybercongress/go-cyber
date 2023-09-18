package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	ctypes "github.com/cybercongress/go-cyber/types"
)

const (
	ActionCreateRoute   = "create_route"
	ActionEditRoute     = "edit_route"
	ActionDeleteRoute   = "delete_route"
	ActionEditRouteName = "edit_route_name"
)

var (
	_ sdk.Msg = &MsgCreateRoute{}
	_ sdk.Msg = &MsgEditRoute{}
	_ sdk.Msg = &MsgDeleteRoute{}
	_ sdk.Msg = &MsgEditRouteName{}
)

func NewMsgCreateRoute(src sdk.AccAddress, dst sdk.AccAddress, name string) *MsgCreateRoute {
	return &MsgCreateRoute{
		Source:      src.String(),
		Destination: dst.String(),
		Name:        name,
	}
}

func (msg MsgCreateRoute) Route() string { return RouterKey }

func (msg MsgCreateRoute) Type() string { return ActionCreateRoute }

func (msg MsgCreateRoute) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid source address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Destination)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid destination address (%s)", err)
	}
	if len(msg.Name) == 0 || len(msg.Name) > 32 { // TODO fix validation
		return ErrWrongName
	}

	return nil
}

func (msg MsgCreateRoute) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateRoute) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgEditRoute(src sdk.AccAddress, dst sdk.AccAddress, value sdk.Coin) *MsgEditRoute {
	return &MsgEditRoute{
		Source:      src.String(),
		Destination: dst.String(),
		Value:       value,
	}
}

func (msg MsgEditRoute) Route() string { return RouterKey }

func (msg MsgEditRoute) Type() string { return ActionEditRoute }

func (msg MsgEditRoute) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid source address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Destination)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid destination address (%s)", err)
	}
	if msg.Value.Denom != ctypes.AMPERE && msg.Value.Denom != ctypes.VOLT {
		return ErrWrongValueDenom
	}

	return nil
}

func (msg MsgEditRoute) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgEditRoute) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgDeleteRoute(src sdk.AccAddress, dst sdk.AccAddress) *MsgDeleteRoute {
	return &MsgDeleteRoute{
		Source:      src.String(),
		Destination: dst.String(),
	}
}

func (msg MsgDeleteRoute) Route() string { return RouterKey }

func (msg MsgDeleteRoute) Type() string { return ActionDeleteRoute }

func (msg MsgDeleteRoute) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid source address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Destination)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid destination address (%s)", err)
	}

	return nil
}

func (msg MsgDeleteRoute) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgDeleteRoute) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgEditRouteName(src sdk.AccAddress, dst sdk.AccAddress, name string) *MsgEditRouteName {
	return &MsgEditRouteName{
		Source:      src.String(),
		Destination: dst.String(),
		Name:        name,
	}
}

func (msg MsgEditRouteName) Route() string { return RouterKey }

func (msg MsgEditRouteName) Type() string { return ActionEditRouteName }

func (msg MsgEditRouteName) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid source address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Destination)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid destination address (%s)", err)
	}
	if len(msg.Name) == 0 || len(msg.Name) > 32 {
		return ErrWrongName
	}

	return nil
}

func (msg MsgEditRouteName) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgEditRouteName) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
