package wasm

import (
	"encoding/json"

	"github.com/ipfs/go-cid"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	graphtypes "github.com/cybercongress/go-cyber/v4/x/graph/types"

	"github.com/cybercongress/go-cyber/v4/x/rank/keeper"
)

type QuerierInterface interface {
	Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

var _ QuerierInterface = Querier{}

type Querier struct {
	keeper *keeper.StateKeeper
}

func NewWasmQuerier(keeper *keeper.StateKeeper) Querier {
	return Querier{keeper}
}

func (Querier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) { return nil, nil }

type CosmosQuery struct {
	ParticleRank *QueryParticleRankParams `json:"particle_rank,omitempty"`
}

type QueryParticleRankParams struct {
	Particle string `json:"particle"`
}

type ParticleRankResponse struct {
	Rank uint64 `json:"rank"`
}

func (querier Querier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
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
