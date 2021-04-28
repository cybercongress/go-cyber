package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmTypes "github.com/CosmWasm/wasmvm/types"

	"github.com/cybercongress/go-cyber/x/bandwidth/keeper"
)

var _ WasmQuerierInterface = WasmQuerier{}

type WasmQuerierInterface interface {
	Query(ctx sdk.Context, request wasmTypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

type WasmQuerier struct {
	*keeper.BandwidthMeter
}

func NewWasmQuerier(keeper *keeper.BandwidthMeter) WasmQuerier {
	return WasmQuerier{keeper}
}

func (WasmQuerier) Query(_ sdk.Context, _ wasmTypes.QueryRequest) ([]byte, error) { return nil, nil }

type CosmosQuery struct {
	Price     		   *struct{} `json:"get_price,omitempty"`
	Load      		   *struct{} `json:"get_load,omitempty"`
	DesirableBandwidth *struct{} `json:"get_desirable_bandwidth,omitempty"`
	AccountBandwidth   *QueryAccountBandwidthParams `json:"get_account_bandwidth,omitempty"`
}

type QueryAccountBandwidthParams struct {
	Address string `json:"address"`
}

type PriceResponse struct {
	Price string `json:"price"`
}

type LoadResponse struct {
	Load string `json:"load"`
}

type DesirableBandwidthResponse struct {
	DesirableBandwidth uint64 `json:"desirable_bandwidth"`
}

type AccountBandwidthResponse struct {
	Address 		 string `json:"address"`
	RemainedValue 	 uint64 `json:"remained_value"`
	LastUpdatedBlock uint64 `json:"last_updated_block"`
	MaxValue 		 uint64 `json:"max_value"`
}

func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	if query.Price != nil {
		price := querier.BandwidthMeter.GetCurrentCreditPrice()

		bz, err = json.Marshal(PriceResponse{
			Price: price.String(),
		})
	} else if query.Load != nil {
		load := querier.BandwidthMeter.GetCurrentNetworkLoad(ctx)

		bz, err = json.Marshal(LoadResponse{
			Load: load.String(),
		})
	} else if query.DesirableBandwidth != nil {
		desirableBandwidth := querier.BandwidthMeter.GetDesirableBandwidth(ctx)

		bz, err = json.Marshal(DesirableBandwidthResponse{
			DesirableBandwidth: desirableBandwidth,
		})
	} else if query.AccountBandwidth != nil {
		address, _ := sdk.AccAddressFromBech32(query.AccountBandwidth.Address)
		accountBandwidth := querier.BandwidthMeter.GetCurrentAccountBandwidth(ctx, address)

		bz, err = json.Marshal(AccountBandwidthResponse{
			Address: accountBandwidth.Address,
			RemainedValue: accountBandwidth.RemainedValue,
			LastUpdatedBlock: accountBandwidth.LastUpdatedBlock,
			MaxValue: accountBandwidth.MaxValue,
		})
	} else {
		return nil, sdkerrors.ErrInvalidRequest
	}

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}