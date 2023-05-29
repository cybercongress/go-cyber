package liquidity_plugin

import (
	"encoding/json"

	"github.com/cybercongress/go-cyber/v2/plugins"
	liquiditytypes "github.com/tendermint/liquidity/x/liquidity/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmTypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/liquidity/x/liquidity/keeper"
)

var (
	_ plugins.WasmQuerierInterface   = WasmQuerier{}
	_ plugins.WasmMsgParserInterface = WasmMsgParser{}
)

//--------------------------------------------------

type WasmMsgParser struct{}

func NewWasmMsgParser() WasmMsgParser {
	return WasmMsgParser{}
}

func (WasmMsgParser) Parse(_ sdk.AccAddress, _ wasmTypes.CosmosMsg) ([]sdk.Msg, error) {
	return nil, nil
}

type CosmosMsg struct {
	CreatePool          *liquiditytypes.MsgCreatePool          `json:"create_pool,omitempty"`
	DepositWithinBatch  *liquiditytypes.MsgDepositWithinBatch  `json:"deposit_within_batch,omitempty"`
	WithdrawWithinBatch *liquiditytypes.MsgWithdrawWithinBatch `json:"withdraw_within_batch,omitempty"`
	SwapWithinBatch     *liquiditytypes.MsgSwapWithinBatch     `json:"swap_within_batch,omitempty"`
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
	} else if sdkMsg.CreatePool != nil {
		return []sdk.Msg{sdkMsg.CreatePool}, sdkMsg.CreatePool.ValidateBasic()
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
	PoolParams    *QueryPoolParams `json:"pool_params,omitempty"`
	PoolLiquidity *QueryPoolParams `json:"pool_liquidity,omitempty"`
	PoolSupply    *QueryPoolParams `json:"pool_supply,omitempty"`
	PoolPrice     *QueryPoolParams `json:"pool_price,omitempty"`
	PoolAddress   *QueryPoolParams `json:"pool_address,omitempty"`
}

type QueryPoolParams struct {
	PoolId uint64 `json:"pool_id"`
}

type PoolParamsResponse struct {
	TypeId                uint32   `json:"type_id"`
	ReserveCoinDenoms     []string `json:"reserve_coin_denoms"`
	ReserveAccountAddress string   `json:"reserve_account_address"`
	PoolCoinDenom         string   `json:"pool_coin_denom"`
}

type PoolLiquidityResponse struct {
	Liquidity wasmTypes.Coins `json:"liquidity"`
}

type PoolSupplyResponse struct {
	Supply wasmTypes.Coin `json:"supply"`
}

type PoolPriceResponse struct {
	Price string `json:"price"`
}

type PoolAddressResponse struct {
	Address string `json:"address"`
}

func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	if query.PoolParams != nil {
		pool, found := querier.Keeper.GetPool(ctx, query.PoolParams.PoolId)
		if found != true {
			return nil, sdkerrors.ErrInvalidRequest
		}

		bz, err = json.Marshal(
			PoolParamsResponse{
				TypeId:                pool.TypeId,
				ReserveCoinDenoms:     pool.ReserveCoinDenoms,
				ReserveAccountAddress: pool.ReserveAccountAddress,
				PoolCoinDenom:         pool.PoolCoinDenom,
			},
		)
	} else if query.PoolLiquidity != nil {
		pool, found := querier.Keeper.GetPool(ctx, query.PoolLiquidity.PoolId)
		if found != true {
			return nil, sdkerrors.ErrInvalidRequest
		}

		reserveCoins := querier.Keeper.GetReserveCoins(ctx, pool)

		bz, err = json.Marshal(
			PoolLiquidityResponse{
				Liquidity: plugins.ConvertSdkCoinsToWasmCoins(reserveCoins),
			},
		)
	} else if query.PoolSupply != nil {
		pool, found := querier.Keeper.GetPool(ctx, query.PoolSupply.PoolId)
		if found != true {
			return nil, sdkerrors.ErrInvalidRequest
		}

		poolSupply := querier.Keeper.GetPoolCoinTotal(ctx, pool)

		bz, err = json.Marshal(
			PoolSupplyResponse{
				Supply: plugins.ConvertSdkCoinToWasmCoin(poolSupply),
			},
		)
	} else if query.PoolPrice != nil {
		pool, found := querier.Keeper.GetPool(ctx, query.PoolPrice.PoolId)
		if found != true {
			return nil, sdkerrors.ErrInvalidRequest
		}

		reserveCoins := querier.Keeper.GetReserveCoins(ctx, pool)
		X := reserveCoins[0].Amount.ToDec()
		Y := reserveCoins[1].Amount.ToDec()

		bz, err = json.Marshal(
			PoolPriceResponse{
				Price: X.Quo(Y).String(),
			},
		)
	} else if query.PoolAddress != nil {
		pool, found := querier.Keeper.GetPool(ctx, query.PoolAddress.PoolId)
		if found != true {
			return nil, sdkerrors.ErrInvalidRequest
		}

		bz, err = json.Marshal(
			PoolAddressResponse{
				Address: pool.ReserveAccountAddress,
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
