package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ctypes "github.com/cybercongress/go-cyber/types"
)

const (
	ActionInvestmint = "investmint"
)

func NewMsgInvestmint(
	agent sdk.AccAddress,
	amount sdk.Coin,
	resource string,
	length uint64,
) *MsgInvestmint {
	return &MsgInvestmint{
		Agent:    agent.String(),
		Amount:   amount,
		Resource: resource,
		Length:   length,
	}
}

func (msg MsgInvestmint) Route() string { return RouterKey }

func (msg MsgInvestmint) Type() string { return ActionInvestmint }

func (msg MsgInvestmint) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Agent); if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid agent address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if msg.Amount.Denom != ctypes.SCYB {
		return sdkerrors.Wrap(ErrInvalidBaseResource, msg.Resource)
	}

	if msg.Resource != ctypes.VOLT && msg.Resource != ctypes.AMPERE {
		return sdkerrors.Wrap(ErrResourceNotExist, msg.Resource)
	}

	if msg.Length == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid length")
	}

	return nil
}

func (msg MsgInvestmint) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgInvestmint) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Agent)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

