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

//______________________________________________________________________

type MsgCreateEnergyRoute struct {
	Source 		 sdk.AccAddress `json:"source" yaml:"source"`
	Destination  sdk.AccAddress `json:"destination" yaml:"destination"`
	Alias 		 string 		`json:"alias" yaml:"alias"`
}

func NewMsgCreateEnergyRoute(src sdk.AccAddress, dst sdk.AccAddress, alias string) MsgCreateEnergyRoute {
	return MsgCreateEnergyRoute{
		Source: 	 src,
		Destination: dst,
		Alias:	     alias,
	}
}

func (msg MsgCreateEnergyRoute) Route() string { return RouterKey }

func (msg MsgCreateEnergyRoute) Type() string  { return "create_energy_route" }

func (msg MsgCreateEnergyRoute) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Source}
}

func (msg MsgCreateEnergyRoute) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateEnergyRoute) ValidateBasic() error {
	if msg.Source.Empty() {
		return ErrEmptySourceAddr
	}
	if msg.Destination.Empty() {
		return ErrEmptyDestinationAddr
	}
	if len(msg.Alias) == 0 && len(msg.Alias) < 64 {
		return ErrWrongAlias
	}

	return nil
}

//______________________________________________________________________

type MsgEditEnergyRoute struct {
	Source 		 sdk.AccAddress `json:"source" yaml:"source"`
	Destination  sdk.AccAddress `json:"destination" yaml:"destination"`
	Value 		 sdk.Coin 		`json:"value" yaml:"value"`
}

func NewMsgEditEnergyRoute(src sdk.AccAddress, dst sdk.AccAddress, value sdk.Coin) MsgEditEnergyRoute {
	return MsgEditEnergyRoute{
		Source: 	 src,
		Destination: dst,
		Value:		 value,
	}
}

func (msg MsgEditEnergyRoute) Route() string { return RouterKey }

func (msg MsgEditEnergyRoute) Type() string  { return "edit_energy_route" }

func (msg MsgEditEnergyRoute) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Source}
}

func (msg MsgEditEnergyRoute) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgEditEnergyRoute) ValidateBasic() error {
	if msg.Source.Empty() {
		return ErrEmptySourceAddr
	}
	if msg.Destination.Empty() {
		return ErrEmptyDestinationAddr
	}

	return nil
}

//______________________________________________________________________

type MsgDeleteEnergyRoute struct {
	Source 		 sdk.AccAddress `json:"source" yaml:"source"`
	Destination  sdk.AccAddress `json:"destination" yaml:"destination"`
}

func NewMsgDeleteEnergyRoute(src sdk.AccAddress, dst sdk.AccAddress) MsgDeleteEnergyRoute {
	return MsgDeleteEnergyRoute{
		Source: 	 src,
		Destination: dst,
	}
}

func (msg MsgDeleteEnergyRoute) Route() string { return RouterKey }

func (msg MsgDeleteEnergyRoute) Type() string  { return "delete_energy_route" }

func (msg MsgDeleteEnergyRoute) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Source}
}

func (msg MsgDeleteEnergyRoute) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgDeleteEnergyRoute) ValidateBasic() error {
	if msg.Source.Empty() {
		return ErrEmptySourceAddr
	}
	if msg.Destination.Empty() {
		return ErrEmptyDestinationAddr
	}

	return nil
}

//______________________________________________________________________

type MsgEditEnergyRouteAlias struct {
	Source 		 sdk.AccAddress `json:"source" yaml:"source"`
	Destination  sdk.AccAddress `json:"destination" yaml:"destination"`
	Alias 		 string 		`json:"alias" yaml:"alias"`
}

func NewMsgEditEnergyRouteAlias(src sdk.AccAddress, dst sdk.AccAddress, alias string) MsgEditEnergyRouteAlias {
	return MsgEditEnergyRouteAlias{
		Source: 	 src,
		Destination: dst,
		Alias:		 alias,
	}
}

func (msg MsgEditEnergyRouteAlias) Route() string { return RouterKey }

func (msg MsgEditEnergyRouteAlias) Type() string  { return "edit_energy_route_alias" }

func (msg MsgEditEnergyRouteAlias) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Source}
}

func (msg MsgEditEnergyRouteAlias) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgEditEnergyRouteAlias) ValidateBasic() error {
	if msg.Source.Empty() {
		return ErrEmptySourceAddr
	}
	if msg.Destination.Empty() {
		return ErrEmptyDestinationAddr
	}
	if len(msg.Alias) == 0 && len(msg.Alias) < 64 {
		return ErrWrongAlias
	}

	return nil
}


