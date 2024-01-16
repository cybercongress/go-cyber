package wasm

import (
	"encoding/json"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/deep-foundation/deep-chain/x/resources/types"
)

var _ MsgParserInterface = MsgParser{}

type MsgParserInterface interface {
	Parse(contractAddr sdk.AccAddress, msg wasmvmtypes.CosmosMsg) ([]sdk.Msg, error)
	ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error)
}

type MsgParser struct{}

func NewWasmMsgParser() MsgParser {
	return MsgParser{}
}

func (MsgParser) Parse(_ sdk.AccAddress, _ wasmvmtypes.CosmosMsg) ([]sdk.Msg, error) {
	return nil, nil
}

type CosmosMsg struct {
	Investmint *types.MsgInvestmint `json:"investmint,omitempty"`
}

func (MsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var sdkMsg CosmosMsg
	err := json.Unmarshal(data, &sdkMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to parse link custom msg")
	}

	if sdkMsg.Investmint != nil {
		return []sdk.Msg{sdkMsg.Investmint}, sdkMsg.Investmint.ValidateBasic()
	}

	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Resources")
}
