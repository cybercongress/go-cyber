package liquidity_plugin

import (
	"encoding/json"
	"github.com/cybercongress/go-cyber/plugins"
	liquiditytypes "github.com/tendermint/liquidity/x/liquidity/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmTypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/liquidity/x/liquidity/keeper"
)

var _ plugins.WasmQuerierInterface = WasmQuerier{}
var _ plugins.WasmMsgParserInterface = WasmMsgParser{}

//--------------------------------------------------

type WasmMsgParser struct{}

func NewWasmMsgParser() WasmMsgParser {
	return WasmMsgParser{}
}

func (WasmMsgParser) Parse(_ sdk.AccAddress, _ wasmTypes.CosmosMsg) ([]sdk.Msg, error) { return nil, nil }

type CosmosMsg struct {
	DepositWithinBatch	  *liquiditytypes.MsgDepositWithinBatch  `json:"deposit_within_batch,omitempty"`
	WithdrawWithinBatch   *liquiditytypes.MsgWithdrawWithinBatch `json:"withdraw_within_batch,omitempty"`
	SwapWithinBatch       *liquiditytypes.MsgSwapWithinBatch     `json:"swap_within_batch,omitempty"`
}

func (WasmMsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var sdkMsg CosmosMsg
	err := json.Unmarshal(data, &sdkMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to parse link custom msg")
	}

	if sdkMsg.SwapWithinBatch != nil {
		return []sdk.Msg{sdkMsg.SwapWithinBatch}, sdkMsg.SwapWithinBatch.ValidateBasic()
	} else if sdkMsg.DepositWithinBatch != nil {
		return []sdk.Msg{sdkMsg.DepositWithinBatch}, sdkMsg.DepositWithinBatch.ValidateBasic()
	} else if sdkMsg.WithdrawWithinBatch != nil {
		return []sdk.Msg{sdkMsg.WithdrawWithinBatch}, sdkMsg.WithdrawWithinBatch.ValidateBasic()
	}

	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Liquidity")
}

//--------------------------------------------------


type WasmQuerier struct {
	keeper.Keeper
}

func NewWasmQuerier(keeper keeper.Keeper) WasmQuerier {
	return WasmQuerier{keeper}
}

func (WasmQuerier) Query(_ sdk.Context, _ wasmTypes.QueryRequest) ([]byte, error) { return nil, nil }

type CosmosQuery struct {
	PoolParams   *QueryPoolParams `json:"pool_params,omitempty"`
}

type QueryPoolParams struct {
	PoolId uint64 `json:"pool_id"`
}

type PoolParamsResponse struct {
	Id 					  uint64 `json:"id"`
	TypeId 				  uint32 `json:"type_id"`
	ReserveCoinDenoms 	  []string `json:"reserve_coin_denoms"`
	ReserveAccountAddress string `json:"reserve_account_address"`
	PoolCoinDenom 		  string `json:"pool_coin_denom"`
}

func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	if query.PoolParams != nil {
		pool, found := querier.Keeper.GetPool(ctx, query.PoolParams.PoolId); if found != true {
			return nil, sdkerrors.ErrInvalidRequest
		}

		bz, err = json.Marshal(
			PoolParamsResponse{
				Id: 			   	   pool.Id,
				TypeId: 		   	   pool.TypeId,
				ReserveCoinDenoms: 	   pool.ReserveCoinDenoms,
				ReserveAccountAddress: pool.ReserveAccountAddress,
				PoolCoinDenom: 	       pool.PoolCoinDenom,
			})
	} else {
		return nil, sdkerrors.ErrInvalidRequest
	}
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
