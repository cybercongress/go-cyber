package plugins

import (
	"encoding/json"

	liquiditytypes "github.com/tendermint/liquidity/x/liquidity/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	dmntypes "github.com/cybercongress/go-cyber/x/dmn/types"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
	gridtypes "github.com/cybercongress/go-cyber/x/grid/types"
	resourcestypes "github.com/cybercongress/go-cyber/x/resources/types"
)

type WasmMsgParserInterface interface {
	Parse(contractAddr sdk.AccAddress, msg wasmvmtypes.CosmosMsg) ([]sdk.Msg, error)
	ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error)
}

type MsgParser struct {
	Parsers map[string]WasmMsgParserInterface
}

func NewMsgParser() MsgParser {
	return MsgParser{
		Parsers: make(map[string]WasmMsgParserInterface),
	}
}

type WasmCustomMsg struct {
	Route   string          `json:"route"`
	MsgData json.RawMessage `json:"msg_data"`
}

const (
	WasmMsgParserRouteGraph     = graphtypes.ModuleName
	WasmMsgParserRouteDmn       = dmntypes.ModuleName
	WasmMsgParserRouteGrid      = gridtypes.ModuleName
	WasmMsgParserRouteResources = resourcestypes.ModuleName
	WasmMsgParserLiquidity      = liquiditytypes.ModuleName
)

func (p MsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var customMsg WasmCustomMsg
	err := json.Unmarshal(data, &customMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if parser, ok := p.Parsers[customMsg.Route]; ok {
		return parser.ParseCustom(contractAddr, customMsg.MsgData)
	}

	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, customMsg.Route)
}
