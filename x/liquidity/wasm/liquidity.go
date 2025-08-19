package wasm

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	liquiditytypes "github.com/cybercongress/go-cyber/v6/x/liquidity/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cybercongress/go-cyber/v6/x/liquidity/keeper"
)

// TODO refactor to the new query and message handler or remove

var (
// _ WasmQuerierInterface   = WasmQuerier{}
// _ WasmMsgParserInterface = WasmMsgParser{}
)

//--------------------------------------------------

type WasmMsgParser struct{}

func NewLiquidityWasmMsgParser() WasmMsgParser {
	return WasmMsgParser{}
}

func (WasmMsgParser) Parse(_ sdk.AccAddress, _ wasmvmtypes.CosmosMsg) ([]sdk.Msg, error) {
	return nil, nil
}

type CosmosMsg struct {
	CreatePool          *liquiditytypes.MsgCreatePool          `json:"create_pool,omitempty"`
	DepositWithinBatch  *liquiditytypes.MsgDepositWithinBatch  `json:"deposit_within_batch,omitempty"`
	WithdrawWithinBatch *liquiditytypes.MsgWithdrawWithinBatch `json:"withdraw_within_batch,omitempty"`
	SwapWithinBatch     *liquiditytypes.MsgSwapWithinBatch     `json:"swap_within_batch,omitempty"`
}

// Deprecated
func (WasmMsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var sdkMsg CosmosMsg
	err := json.Unmarshal(data, &sdkMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to parse link custom msg")
	}

	switch {
	case sdkMsg.SwapWithinBatch != nil:
		return []sdk.Msg{sdkMsg.SwapWithinBatch}, sdkMsg.SwapWithinBatch.ValidateBasic()
	case sdkMsg.DepositWithinBatch != nil:
		return []sdk.Msg{sdkMsg.DepositWithinBatch}, sdkMsg.DepositWithinBatch.ValidateBasic()
	case sdkMsg.WithdrawWithinBatch != nil:
		return []sdk.Msg{sdkMsg.WithdrawWithinBatch}, sdkMsg.WithdrawWithinBatch.ValidateBasic()
	case sdkMsg.CreatePool != nil:
		return []sdk.Msg{sdkMsg.CreatePool}, sdkMsg.CreatePool.ValidateBasic()
	default:
		return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Liquidity")
	}
}

//--------------------------------------------------

type WasmQuerier struct {
	keeper.Keeper
}

func NewLiquidityWasmQuerier(keeper keeper.Keeper) WasmQuerier {
	return WasmQuerier{keeper}
}

func (WasmQuerier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) { return nil, nil }

type CosmosQuery struct {
	PoolParams    *QueryPoolParams `json:"pool_params,omitempty"`
	PoolLiquidity *QueryPoolParams `json:"pool_liquidity,omitempty"`
	PoolSupply    *QueryPoolParams `json:"pool_supply,omitempty"`
	PoolPrice     *QueryPoolParams `json:"pool_price,omitempty"`
	PoolAddress   *QueryPoolParams `json:"pool_address,omitempty"`
}

type QueryPoolParams struct {
	PoolID uint64 `json:"pool_id"`
}

type PoolParamsResponse struct {
	TypeID                uint32   `json:"type_id"`
	ReserveCoinDenoms     []string `json:"reserve_coin_denoms"`
	ReserveAccountAddress string   `json:"reserve_account_address"`
	PoolCoinDenom         string   `json:"pool_coin_denom"`
}

type PoolLiquidityResponse struct {
	Liquidity wasmvmtypes.Coins `json:"liquidity"`
}

type PoolSupplyResponse struct {
	Supply wasmvmtypes.Coin `json:"supply"`
}

type PoolPriceResponse struct {
	Price string `json:"price"`
}

type PoolAddressResponse struct {
	Address string `json:"address"`
}

// Deprecated
func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	switch {
	case query.PoolParams != nil:
		pool, found := querier.Keeper.GetPool(ctx, query.PoolParams.PoolID)
		if !found {
			return nil, sdkerrors.ErrInvalidRequest
		}

		bz, err = json.Marshal(
			PoolParamsResponse{
				TypeID:                pool.TypeId,
				ReserveCoinDenoms:     pool.ReserveCoinDenoms,
				ReserveAccountAddress: pool.ReserveAccountAddress,
				PoolCoinDenom:         pool.PoolCoinDenom,
			},
		)
	case query.PoolLiquidity != nil:
		pool, found := querier.Keeper.GetPool(ctx, query.PoolLiquidity.PoolID)
		if !found {
			return nil, sdkerrors.ErrInvalidRequest
		}

		reserveCoins := querier.Keeper.GetReserveCoins(ctx, pool)

		bz, err = json.Marshal(
			PoolLiquidityResponse{
				Liquidity: wasmkeeper.ConvertSdkCoinsToWasmCoins(reserveCoins),
			},
		)
	case query.PoolSupply != nil:
		pool, found := querier.Keeper.GetPool(ctx, query.PoolSupply.PoolID)
		if !found {
			return nil, sdkerrors.ErrInvalidRequest
		}

		poolSupply := querier.Keeper.GetPoolCoinTotal(ctx, pool)

		bz, err = json.Marshal(
			PoolSupplyResponse{
				Supply: wasmkeeper.ConvertSdkCoinToWasmCoin(poolSupply),
			},
		)
	case query.PoolPrice != nil:
		pool, found := querier.Keeper.GetPool(ctx, query.PoolPrice.PoolID)
		if !found {
			return nil, sdkerrors.ErrInvalidRequest
		}

		reserveCoins := querier.Keeper.GetReserveCoins(ctx, pool)
		X := reserveCoins[0].Amount.ToLegacyDec()
		Y := reserveCoins[1].Amount.ToLegacyDec()

		bz, err = json.Marshal(
			PoolPriceResponse{
				Price: X.Quo(Y).String(),
			},
		)
	case query.PoolAddress != nil:
		pool, found := querier.Keeper.GetPool(ctx, query.PoolAddress.PoolID)
		if !found {
			return nil, sdkerrors.ErrInvalidRequest
		}

		bz, err = json.Marshal(
			PoolAddressResponse{
				Address: pool.ReserveAccountAddress,
			},
		)
	default:
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Liquidity variant"}
	}

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
