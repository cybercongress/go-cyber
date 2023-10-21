package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	ctypes "github.com/cybercongress/go-cyber/v2/types"
)

const (
	ActionInvestmint = "investmint"
)

func NewMsgInvestmint(
	neuron sdk.AccAddress,
	amount sdk.Coin,
	resource string,
	length uint64,
) *MsgInvestmint {
	return &MsgInvestmint{
		Neuron:   neuron.String(),
		Amount:   amount,
		Resource: resource,
		Length:   length,
	}
}

func (msg MsgInvestmint) Route() string { return RouterKey }

func (msg MsgInvestmint) Type() string { return ActionInvestmint }

func (msg MsgInvestmint) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Neuron)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid neuron address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if msg.Resource != ctypes.VOLT && msg.Resource != ctypes.AMPERE {
		return sdkerrors.Wrap(ErrResourceNotExist, msg.Resource)
	}

	if msg.Length == 0 {
		return ErrNotAvailablePeriod
	}

	return nil
}

func (msg MsgInvestmint) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgInvestmint) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Neuron)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
