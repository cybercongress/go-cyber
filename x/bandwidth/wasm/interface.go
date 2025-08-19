package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"

	pluginstypes "github.com/cybercongress/go-cyber/v6/plugins/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v6/x/bandwidth/keeper"
)

type Querier struct {
	*keeper.BandwidthMeter
}

func NewWasmQuerier(keeper *keeper.BandwidthMeter) *Querier {
	return &Querier{keeper}
}

func (querier *Querier) HandleQuery(ctx sdk.Context, query pluginstypes.CyberQuery) ([]byte, error) {
	switch {
	case query.BandwidthLoad != nil:
		res, err := querier.BandwidthMeter.Load(ctx, query.BandwidthLoad)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to get bandwidth load")
		}

		responseBytes, err := json.Marshal(res)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to serialize bandwidth load response")
		}
		return responseBytes, nil

	case query.BandwidthPrice != nil:
		res, err := querier.BandwidthMeter.Price(ctx, query.BandwidthPrice)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to get bandwidth price")
		}

		responseBytes, err := json.Marshal(res)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to serialize bandwidth price response")
		}
		return responseBytes, nil
	case query.TotalBandwidth != nil:
		res, err := querier.BandwidthMeter.TotalBandwidth(ctx, query.TotalBandwidth)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to get total bandwidth")
		}

		responseBytes, err := json.Marshal(res)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to serialize total bandwidth response")
		}
		return responseBytes, nil
	case query.NeuronBandwidth != nil:
		res, err := querier.BandwidthMeter.NeuronBandwidth(ctx, query.NeuronBandwidth)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to get neuron bandwidth")
		}

		responseBytes, err := json.Marshal(res)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to serialize neuron bandwidth response")
		}
		return responseBytes, nil
	default:
		return nil, pluginstypes.ErrHandleQuery
	}
}
