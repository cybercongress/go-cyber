package wasm

import (
	"encoding/json"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cybercongress/go-cyber/x/bandwidth/keeper"
)

var _ WasmQuerierInterface = WasmQuerier{}

type WasmQuerierInterface interface {
	Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

type WasmQuerier struct {
	*keeper.BandwidthMeter
}

func NewWasmQuerier(keeper *keeper.BandwidthMeter) WasmQuerier {
	return WasmQuerier{keeper}
}

func (WasmQuerier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) { return nil, nil }

type CosmosQuery struct {
	BandwidthPrice  *struct{}                   `json:"bandwidth_price,omitempty"`
	BandwidthLoad   *struct{}                   `json:"bandwidth_load,omitempty"`
	BandwidthTotal  *struct{}                   `json:"bandwidth_total,omitempty"`
	NeuronBandwidth *QueryNeuronBandwidthParams `json:"neuron_bandwidth,omitempty"`
}

type QueryNeuronBandwidthParams struct {
	Neuron string `json:"neuron"`
}

type BandwidthPriceResponse struct {
	Price string `json:"price"`
}

type BandwidthLoadResponse struct {
	Load string `json:"load"`
}

type BandwidthTotalResponse struct {
	Total uint64 `json:"total"`
}

type NeuronBandwidthResponse struct {
	Neuron           string `json:"neuron"`
	RemainedValue    uint64 `json:"remained_value"`
	LastUpdatedBlock uint64 `json:"last_updated_block"`
	MaxValue         uint64 `json:"max_value"`
}

func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	if query.BandwidthPrice != nil {
		price := querier.BandwidthMeter.GetCurrentCreditPrice()

		bz, err = json.Marshal(BandwidthPriceResponse{
			Price: price.String(),
		})
	} else if query.BandwidthLoad != nil {
		load := querier.BandwidthMeter.GetCurrentNetworkLoad(ctx)

		bz, err = json.Marshal(BandwidthLoadResponse{
			Load: load.String(),
		})
	} else if query.BandwidthTotal != nil {
		desirableBandwidth := querier.BandwidthMeter.GetDesirableBandwidth(ctx)

		bz, err = json.Marshal(BandwidthTotalResponse{
			Total: desirableBandwidth,
		})
	} else if query.NeuronBandwidth != nil {
		address, _ := sdk.AccAddressFromBech32(query.NeuronBandwidth.Neuron)
		accountBandwidth := querier.BandwidthMeter.GetCurrentAccountBandwidth(ctx, address)

		bz, err = json.Marshal(NeuronBandwidthResponse{
			Neuron:           accountBandwidth.Neuron,
			RemainedValue:    accountBandwidth.RemainedValue,
			LastUpdatedBlock: accountBandwidth.LastUpdatedBlock,
			MaxValue:         accountBandwidth.MaxValue,
		})
	} else {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Bandwidth variant"}
	}

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
