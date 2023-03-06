package wasm

import (
	"encoding/json"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmTypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cybercongress/go-cyber/x/dmn/keeper"
	"github.com/cybercongress/go-cyber/x/dmn/types"
)

var _ WasmQuerierInterface = WasmQuerier{}
var _ WasmMsgParserInterface = WasmMsgParser{}

//--------------------------------------------------

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
	CreateThought         *types.MsgCreateThought         `json:"create_thought,omitempty"`
	ForgetThought         *types.MsgForgetThought         `json:"forget_thought,omitempty"`
	ChangeThoughtInput    *types.MsgChangeThoughtInput    `json:"change_thought_input,omitempty"`
	ChangeThoughtPeriod   *types.MsgChangeThoughtPeriod   `json:"change_thought_period,omitempty"`
	ChangeThoughtBlock    *types.MsgChangeThoughtBlock    `json:"change_thought_block,omitempty"`
	ChangeThoughtGasPrice *types.MsgChangeThoughtGasPrice `json:"change_thought_gas_price,omitempty"`
	ChangeThoughtParticle *types.MsgChangeThoughtParticle `json:"change_thought_particle,omitempty"`
	ChangeThoughtName     *types.MsgChangeThoughtName     `json:"change_thought_name,omitempty"`
}

func (WasmMsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var sdkMsg CosmosMsg
	err := json.Unmarshal(data, &sdkMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to parse link custom msg")
	}

	if sdkMsg.CreateThought != nil {
		return []sdk.Msg{sdkMsg.CreateThought}, sdkMsg.CreateThought.ValidateBasic()
	} else if sdkMsg.ForgetThought != nil {
		return []sdk.Msg{sdkMsg.ForgetThought}, sdkMsg.ForgetThought.ValidateBasic()
	} else if sdkMsg.ChangeThoughtInput != nil {
		return []sdk.Msg{sdkMsg.ChangeThoughtInput}, sdkMsg.ChangeThoughtInput.ValidateBasic()
	} else if sdkMsg.ChangeThoughtPeriod != nil {
		return []sdk.Msg{sdkMsg.ChangeThoughtPeriod}, sdkMsg.ChangeThoughtPeriod.ValidateBasic()
	} else if sdkMsg.ChangeThoughtBlock != nil {
		return []sdk.Msg{sdkMsg.ChangeThoughtBlock}, sdkMsg.ChangeThoughtBlock.ValidateBasic()
	} else if sdkMsg.ChangeThoughtGasPrice != nil {
		return []sdk.Msg{sdkMsg.ChangeThoughtGasPrice}, sdkMsg.ChangeThoughtGasPrice.ValidateBasic()
	} else if sdkMsg.ChangeThoughtParticle != nil {
		return []sdk.Msg{sdkMsg.ChangeThoughtParticle}, sdkMsg.ChangeThoughtParticle.ValidateBasic()
	} else if sdkMsg.ChangeThoughtName != nil {
		return []sdk.Msg{sdkMsg.ChangeThoughtName}, sdkMsg.ChangeThoughtName.ValidateBasic()
	}

	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of DMN")
}

//--------------------------------------------------

type WasmQuerierInterface interface {
	Query(ctx sdk.Context, request wasmTypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

type WasmQuerier struct {
	keeper.Keeper
}

func NewWasmQuerier(keeper keeper.Keeper) WasmQuerier {
	return WasmQuerier{keeper}
}

func (WasmQuerier) Query(_ sdk.Context, _ wasmTypes.QueryRequest) ([]byte, error) { return nil, nil }

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

func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	if query.Thought != nil {

		program, _ := sdk.AccAddressFromBech32(query.Thought.Program)
		thought, found := querier.Keeper.GetThought(ctx, program, query.Thought.Name)
		if found != true {
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
	} else if query.ThoughtStats != nil {
		program, _ := sdk.AccAddressFromBech32(query.ThoughtStats.Program)
		thoughtStats, found := querier.Keeper.GetThoughtStats(ctx, program, query.ThoughtStats.Name)
		if found != true {
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
	} else if query.LowestFee != nil {
		lowestFee := querier.Keeper.GetLowestFee(ctx)
		bz, err = json.Marshal(
			LowestFeeResponse{
				Fee: wasmkeeper.ConvertSdkCoinToWasmCoin(lowestFee),
			},
		)
	} else {
		return nil, sdkerrors.ErrInvalidRequest
	}
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func convertLoadToWasmLoad(load types.Load) Load {
	return Load{
		load.Input,
		wasmkeeper.ConvertSdkCoinToWasmCoin(load.GasPrice),
	}
}
