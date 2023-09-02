package wasm

import (
	"encoding/json"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
	"github.com/cybercongress/go-cyber/x/rank/keeper"
	"github.com/ipfs/go-cid"
)

type WasmQuerierInterface interface {
	Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

var _ WasmQuerierInterface = WasmQuerier{}

type WasmQuerier struct {
	keeper *keeper.StateKeeper
}

func NewWasmQuerier(keeper *keeper.StateKeeper) WasmQuerier {
	return WasmQuerier{keeper}
}

func (WasmQuerier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) { return nil, nil }

type CosmosQuery struct {
	ParticleRank *QueryParticleRankParams `json:"particle_rank,omitempty"`
}

type QueryParticleRankParams struct {
	Particle string `json:"particle"`
}

type ParticleRankResponse struct {
	Rank uint64 `json:"rank"`
}

func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	if query.ParticleRank != nil {
		particle, err := cid.Decode(query.ParticleRank.Particle)
		if err != nil {
			return nil, graphtypes.ErrInvalidParticle
		}

		if particle.Version() != 0 {
			return nil, graphtypes.ErrCidVersion
		}

		rank, err := querier.keeper.GetRankValueByParticle(ctx, query.ParticleRank.Particle)
		if err != nil {
			return nil, err
		}

		bz, err = json.Marshal(ParticleRankResponse{Rank: rank})
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}

		return bz, err
	}

	return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Rank variant"}
}
