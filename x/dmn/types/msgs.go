package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ipfs/go-cid"

	// sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cybercongress/go-cyber/types"
	graph "github.com/cybercongress/go-cyber/x/graph/types"
)

const (
	TypeMsgCreateThought         = "create_thought"
	TypeMsgForgetThought         = "forget_thought"
	TypeMsgChangeThoughtName     = "change_thought_name"
	TypeMsgChangeThoughtParticle = "change_thought_particle"
	TypeMsgChangeThoughtInput    = "change_thought_input"
	TypeMsgChangeThoughtGasPrice = "change_thought_gas_price"
	TypeMsgChangeThoughtPeriod   = "change_thought_period"
	TypeMsgChangeThoughtBlock    = "change_thought_block"
)

func NewMsgCreateThought(
	program sdk.AccAddress,
	trigger Trigger,
	load Load,
	name string,
	particle string,
) *MsgCreateThought {
	return &MsgCreateThought{
		Program:  program.String(),
		Trigger:  trigger,
		Load:     load,
		Name:     name,
		Particle: particle,
	}
}

func (msg MsgCreateThought) Route() string { return RouterKey }

func (msg MsgCreateThought) Type() string { return TypeMsgCreateThought }

func (msg MsgCreateThought) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if msg.Load.Input == "" || len(msg.Load.Input) > 2048 {
		return ErrBadCallData
	}
	if msg.Load.GasPrice.Denom != types.CYB {
		return ErrBadGasPrice
	}
	if !msg.Load.GasPrice.Amount.IsPositive() {
		return ErrBadGasPrice
	}
	if msg.Trigger.Period == 0 && msg.Trigger.Block == 0 {
		return ErrBadTrigger
	}
	if msg.Trigger.Period > 0 && msg.Trigger.Block > 0 {
		return ErrBadTrigger
	}
	if msg.Name == "" || len(msg.Name) > 32 {
		return ErrBadName
	}
	particle, err := cid.Decode(msg.Particle)
	if err != nil {
		return graph.ErrInvalidParticle
	}
	if particle.Version() != 0 {
		return graph.ErrCidVersion
	}

	return nil
}

func (msg MsgCreateThought) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateThought) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgForgetThought(program sdk.AccAddress, label string) *MsgForgetThought {
	return &MsgForgetThought{
		Program: program.String(),
		Name:    label,
	}
}

func (msg MsgForgetThought) Route() string { return RouterKey }

func (msg MsgForgetThought) Type() string { return TypeMsgForgetThought }

func (msg MsgForgetThought) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if msg.Name == "" || len(msg.Name) > 32 {
		return ErrBadName
	}

	return nil
}

func (msg MsgForgetThought) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgForgetThought) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgChangeThoughtName(
	program sdk.AccAddress,
	name, newName string,
) *MsgChangeThoughtName {
	return &MsgChangeThoughtName{
		Program: program.String(),
		Name:    name,
		NewName: newName,
	}
}

func (msg MsgChangeThoughtName) Route() string { return RouterKey }

func (msg MsgChangeThoughtName) Type() string { return TypeMsgChangeThoughtName }

func (msg MsgChangeThoughtName) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if msg.Name == "" || len(msg.Name) > 32 {
		return ErrBadName
	}
	if msg.NewName == "" || len(msg.NewName) > 32 {
		return ErrBadName
	}

	return nil
}

func (msg MsgChangeThoughtName) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeThoughtName) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgChangeThoughtParticle(
	program sdk.AccAddress,
	label string,
	cid string,
) *MsgChangeThoughtParticle {
	return &MsgChangeThoughtParticle{
		Program:  program.String(),
		Name:     label,
		Particle: cid,
	}
}

func (msg MsgChangeThoughtParticle) Route() string { return RouterKey }

func (msg MsgChangeThoughtParticle) Type() string { return TypeMsgChangeThoughtParticle }

func (msg MsgChangeThoughtParticle) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	particle, err := cid.Decode(msg.Particle)
	if err != nil {
		return graph.ErrInvalidParticle
	}
	if particle.Version() != 0 {
		return graph.ErrCidVersion
	}
	if msg.Name == "" || len(msg.Name) > 32 {
		return ErrBadName
	}

	return nil
}

func (msg MsgChangeThoughtParticle) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeThoughtParticle) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgChangeCallInput(
	program sdk.AccAddress,
	name string,
	input string,
) *MsgChangeThoughtInput {
	return &MsgChangeThoughtInput{
		Program: program.String(),
		Name:    name,
		Input:   input,
	}
}

func (msg MsgChangeThoughtInput) Route() string { return RouterKey }

func (msg MsgChangeThoughtInput) Type() string { return TypeMsgChangeThoughtInput }

func (msg MsgChangeThoughtInput) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if msg.Name == "" || len(msg.Name) > 32 {
		return ErrBadName
	}
	if msg.Input == "" || len(msg.Input) > 2048 {
		return ErrBadCallData
	}

	return nil
}

func (msg MsgChangeThoughtInput) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeThoughtInput) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgChangeThoughtGasPrice(
	program sdk.AccAddress,
	name string,
	gasprice sdk.Coin,
) *MsgChangeThoughtGasPrice {
	return &MsgChangeThoughtGasPrice{
		Program:  program.String(),
		Name:     name,
		GasPrice: gasprice,
	}
}

func (msg MsgChangeThoughtGasPrice) Route() string { return RouterKey }

func (msg MsgChangeThoughtGasPrice) Type() string { return TypeMsgChangeThoughtGasPrice }

func (msg MsgChangeThoughtGasPrice) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if msg.Name == "" || len(msg.Name) > 32 {
		return ErrBadName
	}
	if msg.GasPrice.Denom != types.CYB {
		return ErrBadGasPrice
	}
	if !msg.GasPrice.Amount.IsPositive() {
		return ErrBadGasPrice
	}

	return nil
}

func (msg MsgChangeThoughtGasPrice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeThoughtGasPrice) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgChangeThoughtPeriod(
	program sdk.AccAddress,
	name string,
	period uint64,
) *MsgChangeThoughtPeriod {
	return &MsgChangeThoughtPeriod{
		Program: program.String(),
		Name:    name,
		Period:  period,
	}
}

func (msg MsgChangeThoughtPeriod) Route() string { return RouterKey }

func (msg MsgChangeThoughtPeriod) Type() string { return TypeMsgChangeThoughtPeriod }

func (msg MsgChangeThoughtPeriod) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if msg.Name == "" || len(msg.Name) > 32 {
		return ErrBadName
	}
	if msg.Period == 0 {
		return ErrBadTrigger
	}

	return nil
}

func (msg MsgChangeThoughtPeriod) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeThoughtPeriod) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgChangeThoughtBlock(
	program sdk.AccAddress,
	name string,
	block uint64,
) *MsgChangeThoughtBlock {
	return &MsgChangeThoughtBlock{
		Program: program.String(),
		Name:    name,
		Block:   block,
	}
}

func (msg MsgChangeThoughtBlock) Route() string { return RouterKey }

func (msg MsgChangeThoughtBlock) Type() string { return TypeMsgChangeThoughtBlock }

func (msg MsgChangeThoughtBlock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if msg.Name == "" || len(msg.Name) > 32 {
		return ErrBadName
	}
	if msg.Block == 0 {
		return ErrBadTrigger
	}

	return nil
}

func (msg MsgChangeThoughtBlock) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgChangeThoughtBlock) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}
