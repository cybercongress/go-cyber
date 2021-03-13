package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreateEnergyRoute{}
	_ sdk.Msg = &MsgEditEnergyRoute{}
	_ sdk.Msg = &MsgDeleteEnergyRoute{}
	_ sdk.Msg = &MsgEditEnergyRouteAlias{}
)


func NewMsgCreateEnergyRoute(src sdk.AccAddress, dst sdk.AccAddress, alias string) *MsgCreateEnergyRoute {
	return &MsgCreateEnergyRoute{
		Source: 	 src.String(),
		Destination: dst.String(),
		Alias:	     alias,
	}
}

func (msg MsgCreateEnergyRoute) Route() string { return RouterKey }

func (msg MsgCreateEnergyRoute) Type() string  { return ActionCreateEnergyRoute }

func (msg MsgCreateEnergyRoute) ValidateBasic() error {
	if len(msg.Source) == 0 {
		return ErrEmptySourceAddr
	}
	if len(msg.Destination) == 0 {
		return ErrEmptyDestinationAddr
	}
	if len(msg.Alias) == 0 && len(msg.Alias) < 64 {
		return ErrWrongAlias
	}

	return nil
}

func (msg MsgCreateEnergyRoute) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCreateEnergyRoute) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Source)
	return []sdk.AccAddress{addr}
}


func NewMsgEditEnergyRoute(src sdk.AccAddress, dst sdk.AccAddress, value sdk.Coin) *MsgEditEnergyRoute {
	return &MsgEditEnergyRoute{
		Source: 	 src.String(),
		Destination: dst.String(),
		Value:		 value,
	}
}

func (msg MsgEditEnergyRoute) Route() string { return RouterKey }

func (msg MsgEditEnergyRoute) Type() string  { return ActionEditEnergyRoute }

func (msg MsgEditEnergyRoute) ValidateBasic() error {
	if len(msg.Source) == 0 {
		return ErrEmptySourceAddr
	}
	if len(msg.Destination) == 0 {
		return ErrEmptyDestinationAddr
	}

	return nil
}

func (msg MsgEditEnergyRoute) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgEditEnergyRoute) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Source)
	return []sdk.AccAddress{addr}
}


func NewMsgDeleteEnergyRoute(src sdk.AccAddress, dst sdk.AccAddress) *MsgDeleteEnergyRoute {
	return &MsgDeleteEnergyRoute{
		Source: 	 src.String(),
		Destination: dst.String(),
	}
}

func (msg MsgDeleteEnergyRoute) Route() string { return RouterKey }

func (msg MsgDeleteEnergyRoute) Type() string  { return ActionDeleteEnergyRoute }

func (msg MsgDeleteEnergyRoute) ValidateBasic() error {
	if len(msg.Source) == 0 {
		return ErrEmptySourceAddr
	}
	if len(msg.Destination) == 0 {
		return ErrEmptyDestinationAddr
	}

	return nil
}

func (msg MsgDeleteEnergyRoute) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgDeleteEnergyRoute) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Source)
	return []sdk.AccAddress{addr}
}

func NewMsgEditEnergyRouteAlias(src sdk.AccAddress, dst sdk.AccAddress, alias string) *MsgEditEnergyRouteAlias {
	return &MsgEditEnergyRouteAlias{
		Source: 	 src.String(),
		Destination: dst.String(),
		Alias:		 alias,
	}
}

func (msg MsgEditEnergyRouteAlias) Route() string { return RouterKey }

func (msg MsgEditEnergyRouteAlias) Type() string  { return ActionEditEnergyRouteAlias }

func (msg MsgEditEnergyRouteAlias) ValidateBasic() error {
	if len(msg.Source) == 0 {
		return ErrEmptySourceAddr
	}
	if len(msg.Destination) == 0 {
		return ErrEmptyDestinationAddr
	}
	if len(msg.Alias) == 0 && len(msg.Alias) < 64 {
		return ErrWrongAlias
	}

	return nil
}

func (msg MsgEditEnergyRouteAlias) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgEditEnergyRouteAlias) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Source)
	return []sdk.AccAddress{addr}
}


