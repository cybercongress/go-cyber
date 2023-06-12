package wasm

import (
	"encoding/json"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmTypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cybercongress/go-cyber/x/resources/types"
)

var _ WasmMsgParserInterface = WasmMsgParser{}

type WasmMsgParserInterface interface {
	Parse(contractAddr sdk.AccAddress, msg wasmTypes.CosmosMsg) ([]sdk.Msg, error)
	ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error)
}

type WasmMsgParser struct{}

func NewWasmMsgParser() WasmMsgParser {
	return WasmMsgParser{}
}

func (WasmMsgParser) Parse(_ sdk.AccAddress, _ wasmTypes.CosmosMsg) ([]sdk.Msg, error) {
	return nil, nil
}

type CosmosMsg struct {
	Investmint *types.MsgInvestmint `json:"investmint,omitempty"`
}

func (WasmMsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var sdkMsg CosmosMsg
	err := json.Unmarshal(data, &sdkMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to parse link custom msg")
	}

	if sdkMsg.Investmint != nil {
		return []sdk.Msg{sdkMsg.Investmint}, sdkMsg.Investmint.ValidateBasic()
	}

	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown Resources variant")
}
