package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	ctypes "github.com/cybercongress/go-cyber/v6/types"
)

const (
	TypeMsgInvestmint = "investmint"
)

var (
	_ sdk.Msg = &MsgInvestmint{}
	_ sdk.Msg = &MsgUpdateParams{}
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

func (msg MsgInvestmint) Type() string { return TypeMsgInvestmint }

func (msg MsgInvestmint) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Neuron)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid neuron address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if msg.Resource != ctypes.VOLT && msg.Resource != ctypes.AMPERE {
		return errorsmod.Wrap(ErrResourceNotExist, msg.Resource)
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

// GetSignBytes implements the LegacyMsg interface.
func (m MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	return m.Params.Validate()
}
