package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ipfs/go-cid"

	//sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cybercongress/go-cyber/types"
	graph "github.com/cybercongress/go-cyber/x/graph/types"
)

const (
	ActionCronAddJob 			= "add_job"
	ActionCronRemoveJob 		= "remove_job"
	ActionCronChangeJobLabel 	= "change_label"
	ActionCronChangeJobCID 		= "change_cid"
	ActionCronChangeJobCallData = "change_call_data"
	ActionCronChangeJobGasPrice = "change_gas_price"
	ActionCronChangeJobPeriod   = "change_period"
	ActionCronChangeJobBlock    = "change_block"
)

func NewMsgAddJob(
	program sdk.AccAddress,
	trigger Trigger,
	load Load,
	label string,
	cid string,
) *MsgAddJob {
	return &MsgAddJob{
		Program:  program.String(),
		Trigger:  trigger,
		Load:     load,
		Label:    label,
		Cid:      cid,
	}
}

func (msg MsgAddJob) Route() string { return RouterKey }

func (msg MsgAddJob) Type() string  { return ActionCronAddJob }

func (msg MsgAddJob) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if msg.Load.CallData == "" || len(msg.Load.CallData) > 512 {
		return ErrBadCallData
	}
	if msg.Load.GasPrice.Denom != types.CYB {
		return ErrBadGasPrice
	}
	if !msg.Load.GasPrice.Amount.IsPositive() {
		return ErrBadGasPrice
	}
	if (msg.Trigger.Period == 0 && msg.Trigger.Block == 0) {
		return ErrBadTrigger
	}
	if (msg.Trigger.Period > 0 && msg.Trigger.Block > 0) {
		return ErrBadTrigger
	}
	if msg.Label == "" || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if _, err := cid.Decode(string(msg.Cid)); err != nil {
		return graph.ErrInvalidCid
	}

	return nil
}

func (msg MsgAddJob) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgAddJob) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

//______________________________________________________________________

func NewMsgRemoveJob(program sdk.AccAddress, label string) *MsgRemoveJob {
	return &MsgRemoveJob{
		Program:  program.String(),
		Label: label,
	}
}

func (msg MsgRemoveJob) Route() string { return RouterKey }

func (msg MsgRemoveJob) Type() string  { return ActionCronRemoveJob }

func (msg MsgRemoveJob) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if msg.Label == "" || len(msg.Label) > 32 {
		return ErrBadLabel
	}

	return nil
}

func (msg MsgRemoveJob) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgRemoveJob) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgChangeJobLabel(
	program sdk.AccAddress,
	label, newLabel string,
) *MsgChangeJobLabel {
	return &MsgChangeJobLabel{
		Program:  program.String(),
		Label:    label,
		NewLabel: newLabel,
	}
}

func (msg MsgChangeJobLabel) Route() string { return RouterKey }

func (msg MsgChangeJobLabel) Type() string { return ActionCronChangeJobLabel }

func (msg MsgChangeJobLabel) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if msg.Label == "" || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if msg.NewLabel == "" || len(msg.NewLabel) > 32 {
		return ErrBadLabel
	}

	return nil
}

func (msg MsgChangeJobLabel) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeJobLabel) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgChangeJobCID(
	program sdk.AccAddress,
	label string,
	cid string,
) *MsgChangeJobCID {
	return &MsgChangeJobCID{
		Program: program.String(),
		Label: label,
		Cid: cid,
	}
}

func (msg MsgChangeJobCID) Route() string { return RouterKey }

func (msg MsgChangeJobCID) Type() string { return ActionCronChangeJobCID }

func (msg MsgChangeJobCID) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if _, err := cid.Decode(msg.Cid); err != nil {
		return graph.ErrInvalidCid
	}
	if msg.Label == "" || len(msg.Label) > 32 {
		return ErrBadLabel
	}

	return nil
}

func (msg MsgChangeJobCID) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeJobCID) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgChangeCallData(
	program sdk.AccAddress,
	label string,
	calldata string,
) *MsgChangeJobCallData {
	return &MsgChangeJobCallData{
		Program:  program.String(),
		Label:    label,
		CallData: calldata,
	}
}

func (msg MsgChangeJobCallData) Route() string { return RouterKey }

func (msg MsgChangeJobCallData) Type() string { return ActionCronChangeJobCallData }

func (msg MsgChangeJobCallData) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if msg.Label == "" || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if msg.CallData == "" || len(msg.CallData) > 512 {
		return ErrBadCallData
	}

	return nil
}

func (msg MsgChangeJobCallData) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeJobCallData) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgChangeGasPrice(
	program sdk.AccAddress,
	label string,
	gasprice sdk.Coin,
) *MsgChangeJobGasPrice {
	return &MsgChangeJobGasPrice{
		Program:  program.String(),
		Label:    label,
		GasPrice: gasprice,
	}
}

func (msg MsgChangeJobGasPrice) Route() string { return RouterKey }

func (msg MsgChangeJobGasPrice) Type() string { return ActionCronChangeJobGasPrice }

func (msg MsgChangeJobGasPrice) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if msg.Label == "" || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if msg.GasPrice.Denom != types.CYB {
		return ErrBadGasPrice
	}
	if !msg.GasPrice.Amount.IsPositive() {
		return ErrBadGasPrice
	}

	return nil
}

func (msg MsgChangeJobGasPrice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeJobGasPrice) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgChangeJobPeriod(
	program sdk.AccAddress,
	label string,
	period uint64,
) *MsgChangeJobPeriod {
	return &MsgChangeJobPeriod{
		Program:  program.String(),
		Label:	  label,
		Period:   period,
	}
}

func (msg MsgChangeJobPeriod) Route() string { return RouterKey }

func (msg MsgChangeJobPeriod) Type() string { return ActionCronChangeJobPeriod }

func (msg MsgChangeJobPeriod) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if msg.Label == "" || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if msg.Period == 0 {
		return ErrBadTrigger
	}

	return nil
}

func (msg MsgChangeJobPeriod) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeJobPeriod) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgChangeJobBlock(
	program sdk.AccAddress,
	label string,
	block uint64,
) *MsgChangeJobBlock {
	return &MsgChangeJobBlock{
		Program:  program.String(),
		Label:    label,
		Block:    block,
	}
}

func (msg MsgChangeJobBlock) Route() string { return RouterKey }

func (msg MsgChangeJobBlock) Type() string { return ActionCronChangeJobBlock }

func (msg MsgChangeJobBlock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid program address (%s)", err)
	}
	if msg.Label == "" || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if msg.Block == 0 {
		return ErrBadTrigger
	}

	return nil
}

func (msg MsgChangeJobBlock) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Program)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgChangeJobBlock) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}
