package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ctypes "github.com/cybercongress/go-cyber/types"
)

const(
	ActionCreateRoute 	 = "create_route"
	ActionEditRoute 	 = "edit_route"
	ActionDeleteRoute 	 = "delete_route"
	ActionEditRouteAlias = "edit_route_alias"
)


var (
	_ sdk.Msg = &MsgCreateRoute{}
	_ sdk.Msg = &MsgEditRoute{}
	_ sdk.Msg = &MsgDeleteRoute{}
	_ sdk.Msg = &MsgEditRouteAlias{}
)

func NewMsgCreateRoute(src sdk.AccAddress, dst sdk.AccAddress, alias string) *MsgCreateRoute {
	return &MsgCreateRoute{
		Source: 	 src.String(),
		Destination: dst.String(),
		Alias:	     alias,
	}
}

func (msg MsgCreateRoute) Route() string { return RouterKey }

func (msg MsgCreateRoute) Type() string  { return ActionCreateRoute }

func (msg MsgCreateRoute) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid source address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Destination)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid destination address (%s)", err)
	}
	if len(msg.Alias) == 0 && len(msg.Alias) < 32 { // TODO fix validation
		return ErrWrongAlias
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
		Source: 	 src.String(),
		Destination: dst.String(),
		Value:		 value,
	}
}

func (msg MsgEditRoute) Route() string { return RouterKey }

func (msg MsgEditRoute) Type() string  { return ActionEditRoute }

func (msg MsgEditRoute) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid source address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Destination)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid destination address (%s)", err)
	}
	if msg.Value.Denom != ctypes.AMPER && msg.Value.Denom != ctypes.VOLT {
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
		Source: 	 src.String(),
		Destination: dst.String(),
	}
}

func (msg MsgDeleteRoute) Route() string { return RouterKey }

func (msg MsgDeleteRoute) Type() string  { return ActionDeleteRoute }

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

func NewMsgEditRouteAlias(src sdk.AccAddress, dst sdk.AccAddress, alias string) *MsgEditRouteAlias {
	return &MsgEditRouteAlias{
		Source: 	 src.String(),
		Destination: dst.String(),
		Alias:		 alias,
	}
}

func (msg MsgEditRouteAlias) Route() string { return RouterKey }

func (msg MsgEditRouteAlias) Type() string  { return ActionEditRouteAlias }

func (msg MsgEditRouteAlias) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid source address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Destination)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid destination address (%s)", err)
	}
	if len(msg.Alias) == 0 && len(msg.Alias) < 32 {
		return ErrWrongAlias
	}

	return nil
}

func (msg MsgEditRouteAlias) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgEditRouteAlias) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Source)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}


