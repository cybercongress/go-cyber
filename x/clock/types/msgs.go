package types

import (
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	globalerrors "github.com/cybercongress/go-cyber/v4/app/helpers"
)

const (
	// Sudo Message called on the contracts
	BeginBlockSudoMessage = `{"clock_begin_block":{}}`
	EndBlockSudoMessage   = `{"clock_end_block":{}}`
)

// == MsgUpdateParams ==
const (
	TypeMsgRegisterFeePayContract   = "register_clock_contract"
	TypeMsgUnregisterFeePayContract = "unregister_clock_contract"
	TypeMsgUnjailFeePayContract     = "unjail_clock_contract"
	TypeMsgUpdateParams             = "update_clock_params"
)

var (
	_ sdk.Msg = &MsgRegisterClockContract{}
	_ sdk.Msg = &MsgUnregisterClockContract{}
	_ sdk.Msg = &MsgUnjailClockContract{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// Route returns the name of the module
func (msg MsgRegisterClockContract) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgRegisterClockContract) Type() string { return TypeMsgRegisterFeePayContract }

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterClockContract) ValidateBasic() error {
	return validateAddresses(msg.SenderAddress, msg.ContractAddress)
}

// GetSignBytes encodes the message for signing
func (msg *MsgRegisterClockContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRegisterClockContract) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.SenderAddress)
	return []sdk.AccAddress{from}
}

// Route returns the name of the module
func (msg MsgUnregisterClockContract) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUnregisterClockContract) Type() string { return TypeMsgRegisterFeePayContract }

// ValidateBasic runs stateless checks on the message
func (msg MsgUnregisterClockContract) ValidateBasic() error {
	return validateAddresses(msg.SenderAddress, msg.ContractAddress)
}

// GetSignBytes encodes the message for signing
func (msg *MsgUnregisterClockContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUnregisterClockContract) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.SenderAddress)
	return []sdk.AccAddress{from}
}

// Route returns the name of the module
func (msg MsgUnjailClockContract) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUnjailClockContract) Type() string { return TypeMsgRegisterFeePayContract }

// ValidateBasic runs stateless checks on the message
func (msg MsgUnjailClockContract) ValidateBasic() error {
	return validateAddresses(msg.SenderAddress, msg.ContractAddress)
}

// GetSignBytes encodes the message for signing
func (msg *MsgUnjailClockContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUnjailClockContract) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.SenderAddress)
	return []sdk.AccAddress{from}
}

// NewMsgUpdateParams creates new instance of MsgUpdateParams
func NewMsgUpdateParams(
	sender sdk.Address,
	contractGasLimit uint64,
) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: sender.String(),
		Params:    NewParams(contractGasLimit),
	}
}

// Route returns the name of the module
func (msg MsgUpdateParams) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUpdateParams) Type() string { return TypeMsgUpdateParams }

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrap(err, "invalid authority address")
	}

	return msg.Params.Validate()
}

// ValidateAddresses validates the provided addresses
func validateAddresses(addresses ...string) error {
	for _, address := range addresses {
		if _, err := sdk.AccAddressFromBech32(address); err != nil {
			return errors.Wrapf(globalerrors.ErrInvalidAddress, "invalid address: %s", address)
		}
	}

	return nil
}
