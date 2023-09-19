package wasm

import (
	"encoding/json"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmplugins "github.com/cybercongress/go-cyber/plugins"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/cybercongress/go-cyber/x/dmn/keeper"
	"github.com/cybercongress/go-cyber/x/dmn/types"
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
	CreateThought         *types.MsgCreateThought         `json:"create_thought,omitempty"`
	ForgetThought         *types.MsgForgetThought         `json:"forget_thought,omitempty"`
	ChangeThoughtInput    *types.MsgChangeThoughtInput    `json:"change_thought_input,omitempty"`
	ChangeThoughtPeriod   *types.MsgChangeThoughtPeriod   `json:"change_thought_period,omitempty"`
	ChangeThoughtBlock    *types.MsgChangeThoughtBlock    `json:"change_thought_block,omitempty"`
	ChangeThoughtGasPrice *types.MsgChangeThoughtGasPrice `json:"change_thought_gas_price,omitempty"`
	ChangeThoughtParticle *types.MsgChangeThoughtParticle `json:"change_thought_particle,omitempty"`
	ChangeThoughtName     *types.MsgChangeThoughtName     `json:"change_thought_name,omitempty"`
}

func (MsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var sdkMsg CosmosMsg
	err := json.Unmarshal(data, &sdkMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to parse link custom msg")
	}

	switch {
	case sdkMsg.CreateThought != nil:
		return []sdk.Msg{sdkMsg.CreateThought}, sdkMsg.CreateThought.ValidateBasic()
	case sdkMsg.ForgetThought != nil:
		return []sdk.Msg{sdkMsg.ForgetThought}, sdkMsg.ForgetThought.ValidateBasic()
	case sdkMsg.ChangeThoughtInput != nil:
		return []sdk.Msg{sdkMsg.ChangeThoughtInput}, sdkMsg.ChangeThoughtInput.ValidateBasic()
	case sdkMsg.ChangeThoughtPeriod != nil:
		return []sdk.Msg{sdkMsg.ChangeThoughtPeriod}, sdkMsg.ChangeThoughtPeriod.ValidateBasic()
	case sdkMsg.ChangeThoughtBlock != nil:
		return []sdk.Msg{sdkMsg.ChangeThoughtBlock}, sdkMsg.ChangeThoughtBlock.ValidateBasic()
	case sdkMsg.ChangeThoughtGasPrice != nil:
		return []sdk.Msg{sdkMsg.ChangeThoughtGasPrice}, sdkMsg.ChangeThoughtGasPrice.ValidateBasic()
	case sdkMsg.ChangeThoughtParticle != nil:
		return []sdk.Msg{sdkMsg.ChangeThoughtParticle}, sdkMsg.ChangeThoughtParticle.ValidateBasic()
	case sdkMsg.ChangeThoughtName != nil:
		return []sdk.Msg{sdkMsg.ChangeThoughtName}, sdkMsg.ChangeThoughtName.ValidateBasic()
	default:
		return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of DMN")
	}
}

//--------------------------------------------------

type QuerierInterface interface {
	Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

type Querier struct {
	keeper.Keeper
}

func NewWasmQuerier(keeper keeper.Keeper) Querier {
	return Querier{keeper}
}

func (Querier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) { return nil, nil }

type CosmosQuery struct {
	Thought      *QueryThoughtParams `json:"thought,omitempty"`
	ThoughtStats *QueryThoughtParams `json:"thought_stats,omitempty"`
	LowestFee    *struct{}           `json:"lowest_fee,omitempty"`
}

type Trigger struct {
	Period uint64 `json:"period"`
	Block  uint64 `json:"block"`
}

type Load struct {
	Input    string           `json:"input"`
	GasPrice wasmvmtypes.Coin `json:"gas_price"`
}

type QueryThoughtParams struct {
	Program string `json:"program"`
	Name    string `json:"name"`
}

type ThoughtQueryResponse struct {
	Program  string  `json:"program"`
	Trigger  Trigger `json:"trigger"`
	Load     Load    `json:"load"`
	Name     string  `json:"name"`
	Particle string  `json:"particle"`
}

type ThoughtStatsQueryResponse struct {
	Program   string `json:"program"`
	Name      string `json:"name"`
	Calls     uint64 `json:"calls"`
	Fees      uint64 `json:"fees"`
	Gas       uint64 `json:"gas"`
	LastBlock uint64 `json:"last_block"`
}

type LowestFeeResponse struct {
	Fee wasmvmtypes.Coin `json:"fee"`
}

func (querier Querier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	switch {
	case query.Thought != nil:

		program, _ := sdk.AccAddressFromBech32(query.Thought.Program)
		thought, found := querier.Keeper.GetThought(ctx, program, query.Thought.Name)
		if !found {
			return nil, sdkerrors.ErrInvalidRequest
		}

		bz, err = json.Marshal(
			ThoughtQueryResponse{
				Program:  thought.Program,
				Trigger:  Trigger(thought.Trigger),
				Load:     convertLoadToWasmLoad(thought.Load),
				Name:     thought.Name,
				Particle: thought.Particle,
			})
	case query.ThoughtStats != nil:
		program, _ := sdk.AccAddressFromBech32(query.ThoughtStats.Program)
		thoughtStats, found := querier.Keeper.GetThoughtStats(ctx, program, query.ThoughtStats.Name)
		if !found {
			return nil, sdkerrors.ErrInvalidRequest
		}

		bz, err = json.Marshal(
			ThoughtStatsQueryResponse{
				Program:   thoughtStats.Program,
				Name:      thoughtStats.Name,
				Calls:     thoughtStats.Calls,
				Fees:      thoughtStats.Fees,
				Gas:       thoughtStats.Gas,
				LastBlock: thoughtStats.LastBlock,
			})
	case query.LowestFee != nil:
		lowestFee := querier.Keeper.GetLowestFee(ctx)
		bz, err = json.Marshal(
			LowestFeeResponse{
				Fee: wasmplugins.ConvertSdkCoinToWasmCoin(lowestFee),
			},
		)
	default:
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown DMN variant"}
	}

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func convertLoadToWasmLoad(load types.Load) Load {
	return Load{
		load.Input,
		wasmplugins.ConvertSdkCoinToWasmCoin(load.GasPrice),
	}
}
