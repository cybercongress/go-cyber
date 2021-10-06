package wasm

import (
	"encoding/json"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	wasmplugins "github.com/cybercongress/go-cyber/plugins"

	wasmTypes "github.com/CosmWasm/wasmvm/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cybercongress/go-cyber/x/cron/keeper"
	"github.com/cybercongress/go-cyber/x/cron/types"
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

func (WasmMsgParser) Parse(_ sdk.AccAddress, _ wasmTypes.CosmosMsg) ([]sdk.Msg, error) { return nil, nil }

type CosmosMsg struct {
	AddJob            *types.MsgAddJob            `json:"add_job,omitempty"`
	RemoveJob         *types.MsgRemoveJob         `json:"remove_job,omitempty"`
	ChangeJobCallData *types.MsgChangeJobCallData `json:"change_job_call_data,omitempty"`
	ChangeJobPeriod   *types.MsgChangeJobPeriod   `json:"change_job_period,omitempty"`
	ChangeJobBlock    *types.MsgChangeJobBlock    `json:"change_job_block,omitempty"`
}

func (WasmMsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var sdkMsg CosmosMsg
	err := json.Unmarshal(data, &sdkMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to parse link custom msg")
	}

	if sdkMsg.AddJob != nil {
		return []sdk.Msg{sdkMsg.AddJob}, sdkMsg.AddJob.ValidateBasic()
	} else if sdkMsg.RemoveJob != nil {
		return []sdk.Msg{sdkMsg.RemoveJob}, sdkMsg.RemoveJob.ValidateBasic()
	} else if sdkMsg.ChangeJobCallData != nil {
		return []sdk.Msg{sdkMsg.ChangeJobCallData}, sdkMsg.ChangeJobCallData.ValidateBasic()
	} else if sdkMsg.ChangeJobPeriod != nil {
		return []sdk.Msg{sdkMsg.ChangeJobPeriod}, sdkMsg.ChangeJobPeriod.ValidateBasic()
	} else if sdkMsg.ChangeJobBlock != nil {
		return []sdk.Msg{sdkMsg.ChangeJobBlock}, sdkMsg.ChangeJobBlock.ValidateBasic()
	}

	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Cron")
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
	Job      *QueryJobParams `json:"get_job,omitempty"`
	JobStats *QueryJobParams `json:"get_job_stats,omitempty"`
	GetLowestFee *struct{}   `json:"get_lowest_fee,omitempty"`
}

type Trigger struct {
	Period uint64 `json:"period"`
	Block  uint64 `json:"block"`
}

type Load struct {
	CallData string `json:"call_data"`
	GasPrice wasmvmtypes.Coin `json:"gas_price"`
}

type QueryJobParams struct {
	Program  string `json:"program"`
	Label    string `json:"label"`
}

type JobQueryResponse struct {
	Program  string `json:"program"`
	Trigger  Trigger `json:"trigger"`
	Load 	 Load 	`json:"load"`
	Label    string `json:"label"`
	Particle string `json:"particle"`
}

type JobStatsQueryResponse struct {
	Program   string `json:"program"`
	Label     string `json:"label"`
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

	if query.Job != nil {

		program, _ := sdk.AccAddressFromBech32(query.Job.Program)
		job, found := querier.Keeper.GetJob(ctx, program, query.Job.Label); if found != true {
			return nil, sdkerrors.ErrInvalidRequest
		}

		bz, err = json.Marshal(
			JobQueryResponse{
				Program:  job.Program,
				Trigger:  Trigger(job.Trigger),
				Load:     convertLoadToWasmLoad(job.Load),
				Label:    job.Label,
				Particle:      job.Particle,
		})
	} else if query.JobStats != nil {
		program, _ := sdk.AccAddressFromBech32(query.JobStats.Program)
		jobStats, found := querier.Keeper.GetJobStats(ctx, program, query.JobStats.Label); if found != true {
			return nil, sdkerrors.ErrInvalidRequest
		}

		bz, err = json.Marshal(
			JobStatsQueryResponse{
				Program:   jobStats.Program,
				Label:     jobStats.Label,
				Calls:     jobStats.Calls,
				Fees:      jobStats.Fees,
				Gas:       jobStats.Gas,
				LastBlock: jobStats.LastBlock,
		})
	} else if query.GetLowestFee != nil {
		lowestFee := querier.Keeper.GetLowestFee(ctx)
		bz, err = json.Marshal(
			LowestFeeResponse{
				Fee: wasmplugins.ConvertSdkCoinToWasmCoin(lowestFee),
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
		load.CallData,
		wasmplugins.ConvertSdkCoinToWasmCoin(load.GasPrice),
	}
}