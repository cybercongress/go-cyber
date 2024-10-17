package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v5/x/rank/keeper"

	pluginstypes "github.com/cybercongress/go-cyber/v5/plugins/types"
)

type Querier struct {
	keeper *keeper.StateKeeper
}

func NewWasmQuerier(keeper *keeper.StateKeeper) *Querier {
	return &Querier{keeper}
}

func (querier *Querier) HandleQuery(ctx sdk.Context, query pluginstypes.CyberQuery) ([]byte, error) {
	switch {
	case query.ParticleRank != nil:
		res, err := querier.keeper.Rank(ctx, query.ParticleRank)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to get particle rank")
		}

		responseBytes, err := json.Marshal(res)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to serialize particle rank response")
		}
		return responseBytes, nil
	default:
		return nil, pluginstypes.ErrHandleQuery
	}
}
