package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ipfs/go-cid"

	//sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cybercongress/go-cyber/types"
	graph "github.com/cybercongress/go-cyber/x/graph/types"
)

func NewMsgCronAddJob(
	address, contract sdk.AccAddress,
	trigger Trigger,
	load Load,
	label string,
	cid string,
) *MsgCronAddJob {
	return &MsgCronAddJob{
		Creator:  address.String(),
		Contract: contract.String(),
		Trigger:  &trigger,
		Load:     &load,
		Label:    label,
		Cid:      cid,
	}
}

func (msg MsgCronAddJob) Route() string { return RouterKey }

func (msg MsgCronAddJob) Type() string  { return ActionCronAddJob }

func (msg MsgCronAddJob) ValidateBasic() error {
	if len(msg.Creator) == 0 {
		return ErrInvalidAddress
	}
	if len(msg.Contract) == 0 {
		return ErrEmptyContract
	}
	if len(msg.Load.CallData) == 0 || len(msg.Load.CallData) > 256 {
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
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if _, err := cid.Decode(string(msg.Cid)); err != nil {
		return graph.ErrInvalidCid
	}

	return nil
}

func (msg MsgCronAddJob) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCronAddJob) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{addr}
}

//______________________________________________________________________

func NewMsgCronRemoveJob(address, contract sdk.AccAddress, label string) *MsgCronRemoveJob {
	return &MsgCronRemoveJob{
		Creator:  address.String(),
		Contract: contract.String(),
		Label: label,
	}
}

func (msg MsgCronRemoveJob) Route() string { return RouterKey }

func (msg MsgCronRemoveJob) Type() string  { return ActionCronRemoveJob }

func (msg MsgCronRemoveJob) ValidateBasic() error {
	if len(msg.Creator) == 0 {
		return ErrInvalidAddress
	}
	if len(msg.Contract) == 0 {
		return ErrEmptyContract
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}

	return nil
}

func (msg MsgCronRemoveJob) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCronRemoveJob) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{addr}
}

//______________________________________________________________________

func NewMsgCronChangeJobLabel(
	address, contract sdk.AccAddress,
	label, newLabel string,
) *MsgCronChangeJobLabel {
	return &MsgCronChangeJobLabel{
		Creator:  address.String(),
		Contract: contract.String(),
		Label:    label,
		NewLabel: newLabel,
	}
}

func (msg MsgCronChangeJobLabel) Route() string { return RouterKey }

func (msg MsgCronChangeJobLabel) Type() string { return ActionCronChangeJobLabel }

func (msg MsgCronChangeJobLabel) ValidateBasic() error {
	if len(msg.Creator) == 0 {
		return ErrInvalidAddress
	}
	if len(msg.Contract) == 0 {
		return ErrEmptyContract
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if len(msg.NewLabel) == 0 || len(msg.NewLabel) > 32 {
		return ErrBadLabel
	}

	return nil
}

func (msg MsgCronChangeJobLabel) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCronChangeJobLabel) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{addr}
}

//______________________________________________________________________

func NewMsgCronChangeJobCID(
	address, contract sdk.AccAddress,
	label string,
	cid string,
) *MsgCronChangeJobCID {
	return &MsgCronChangeJobCID{
		Creator:  address.String(),
		Contract: contract.String(),
		Label: label,
		Cid: cid,
	}
}

func (msg MsgCronChangeJobCID) Route() string { return RouterKey }

func (msg MsgCronChangeJobCID) Type() string { return ActionCronChangeJobCID }

func (msg MsgCronChangeJobCID) ValidateBasic() error {
	if len(msg.Creator) == 0 {
		return ErrInvalidAddress
	}
	if len(msg.Contract) == 0 {
		return ErrEmptyContract
	}
	if _, err := cid.Decode(string(msg.Cid)); err != nil {
		return graph.ErrInvalidCid
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}

	return nil
}

func (msg MsgCronChangeJobCID) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCronChangeJobCID) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{addr}
}

//______________________________________________________________________

func NewMsgCronChangeCallData(
	address, contract sdk.AccAddress,
	label string,
	calldata string,
) *MsgCronChangeJobCallData {
	return &MsgCronChangeJobCallData{
		Creator:  address.String(),
		Contract: contract.String(),
		Label:    label,
		CallData: calldata,
	}
}

func (msg MsgCronChangeJobCallData) Route() string { return RouterKey }

func (msg MsgCronChangeJobCallData) Type() string { return ActionCronChangeJobCallData }

func (msg MsgCronChangeJobCallData) ValidateBasic() error {
	if len(msg.Creator) == 0 {
		return ErrInvalidAddress
	}
	if len(msg.Contract) == 0 {
		return ErrEmptyContract
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if len(msg.CallData) == 0 || len(msg.CallData) > 256 {
		return ErrBadCallData
	}

	return nil
}

func (msg MsgCronChangeJobCallData) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCronChangeJobCallData) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{addr}
}

//______________________________________________________________________

func NewMsgCronChangeGasPrice(
	address, contract sdk.AccAddress,
	label string,
	gasprice sdk.Coin,
) *MsgCronChangeJobGasPrice {
	return &MsgCronChangeJobGasPrice{
		Creator:  address.String(),
		Contract: contract.String(),
		Label:    label,
		GasPrice: gasprice,
	}
}

func (msg MsgCronChangeJobGasPrice) Route() string { return RouterKey }

func (msg MsgCronChangeJobGasPrice) Type() string { return ActionCronChangeJobGasPrice }

func (msg MsgCronChangeJobGasPrice) ValidateBasic() error {
	if len(msg.Creator) == 0 {
		return ErrInvalidAddress
	}
	if len(msg.Contract) == 0 {
		return ErrEmptyContract
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
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

func (msg MsgCronChangeJobGasPrice) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCronChangeJobGasPrice) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{addr}
}

//______________________________________________________________________

func NewMsgCronChangeJobPeriod(
	address, contract sdk.AccAddress,
	label string,
	period uint64,
) *MsgCronChangeJobPeriod {
	return &MsgCronChangeJobPeriod{
		Creator:  address.String(),
		Contract: contract.String(),
		Label:	  label,
		Period:   period,
	}
}

func (msg MsgCronChangeJobPeriod) Route() string { return RouterKey }

func (msg MsgCronChangeJobPeriod) Type() string { return ActionCronChangeJobPeriod }

func (msg MsgCronChangeJobPeriod) ValidateBasic() error {
	if len(msg.Creator) == 0 {
		return ErrInvalidAddress
	}
	if len(msg.Contract) == 0 {
		return ErrEmptyContract
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if msg.Period == 0 {
		return ErrBadTrigger
	}

	return nil
}

func (msg MsgCronChangeJobPeriod) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCronChangeJobPeriod) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{addr}
}

//______________________________________________________________________

func NewMsgCronChangeJobBlock(
	address, contract sdk.AccAddress,
	label string,
	block uint64,
) *MsgCronChangeJobBlock {
	return &MsgCronChangeJobBlock{
		Creator:  address.String(),
		Contract: contract.String(),
		Label:    label,
		Block:    block,
	}
}

func (msg MsgCronChangeJobBlock) Route() string { return RouterKey }

func (msg MsgCronChangeJobBlock) Type() string { return ActionCronChangeJobBlock }

func (msg MsgCronChangeJobBlock) ValidateBasic() error {
	if len(msg.Creator) == 0 {
		return ErrInvalidAddress
	}
	if len(msg.Contract) == 0 {
		return ErrEmptyContract
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if msg.Block == 0 {
		return ErrBadTrigger
	}

	return nil
}

func (msg MsgCronChangeJobBlock) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{addr}
}

func (msg MsgCronChangeJobBlock) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
