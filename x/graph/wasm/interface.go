package wasm

import (
	"encoding/json"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/cybercongress/go-cyber/v4/x/graph/keeper"
	"github.com/cybercongress/go-cyber/v4/x/graph/types"
)

var (
	_ QuerierInterface   = Querier{}
	_ MsgParserInterface = MsgParser{}
)

//--------------------------------------------------

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
	Cyberlink *types.MsgCyberlink `json:"cyberlink,omitempty"`
}

func (MsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
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

type QuerierInterface interface {
	Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

type Querier struct {
	keeper.GraphKeeper
}

func NewWasmQuerier(keeper keeper.GraphKeeper) Querier {
	return Querier{keeper}
}

func (Querier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) { return nil, nil }

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

func (querier Querier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	switch {
	case query.ParticlesAmount != nil:
		amount := querier.GraphKeeper.GetCidsCount(ctx)
		bz, err = json.Marshal(ParticlesAmountQueryResponse{ParticlesAmount: amount})
	case query.CyberlinksAmount != nil:
		amount := querier.GraphKeeper.GetLinksCount(ctx)
		bz, err = json.Marshal(CyberlinksAmountQueryResponse{CyberlinksAmount: amount})
	default:
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Graph variant"}
	}

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
