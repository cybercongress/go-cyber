package wasm

import (
	"encoding/json"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/cybercongress/go-cyber/x/graph/keeper"
	"github.com/cybercongress/go-cyber/x/graph/types"
)

var _ WasmQuerierInterface = WasmQuerier{}
var _ WasmMsgParserInterface = WasmMsgParser{}

//--------------------------------------------------

type WasmMsgParserInterface interface {
	Parse(contractAddr sdk.AccAddress, msg wasmvmtypes.CosmosMsg) ([]sdk.Msg, error)
	ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error)
}

type WasmMsgParser struct{}

func NewWasmMsgParser() WasmMsgParser {
	return WasmMsgParser{}
}

func (WasmMsgParser) Parse(_ sdk.AccAddress, _ wasmvmtypes.CosmosMsg) ([]sdk.Msg, error) { return nil, nil }

type CosmosMsg struct {
	Cyberlink *types.MsgCyberlink `json:"cyberlink,omitempty"`
}

func (WasmMsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var sdkMsg CosmosMsg
	err := json.Unmarshal(data, &sdkMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to parse graph custom msg")
	}

	if sdkMsg.Cyberlink != nil {
		return []sdk.Msg{sdkMsg.Cyberlink}, sdkMsg.Cyberlink.ValidateBasic()
	}

	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Graph")
}

//--------------------------------------------------

type WasmQuerierInterface interface {
	Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

type WasmQuerier struct {
	keeper.GraphKeeper
}

func NewWasmQuerier(keeper keeper.GraphKeeper) WasmQuerier {
	return WasmQuerier{keeper}
}

func (WasmQuerier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) { return nil, nil }

type CosmosQuery struct {
	ParticlesAmount  *struct{} `json:"particles_amount,omitempty"`
	CyberlinksAmount *struct{} `json:"cyberlinks_amount,omitempty"`
}

type ParticlesAmountQueryResponse struct {
	ParticlesAmount uint64 `json:"particles_amount"`
}

type CyberlinksAmountQueryResponse struct {
	CyberlinksAmount uint64 `json:"cyberlinks_amount"`
}

func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	if query.ParticlesAmount != nil {
		amount := querier.GraphKeeper.GetCidsCount(ctx)
		bz, err = json.Marshal(ParticlesAmountQueryResponse{ParticlesAmount: amount})
	} else if query.CyberlinksAmount != nil {
		amount := querier.GraphKeeper.GetLinksCount(ctx)
		bz, err = json.Marshal(CyberlinksAmountQueryResponse{CyberlinksAmount: amount})
	} else {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Graph variant"}
	}

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}