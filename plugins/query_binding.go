package plugins

import (
	"encoding/json"

	liquiditytypes "github.com/gravity-devs/liquidity/x/liquidity/types"

	wasm "github.com/CosmWasm/wasmd/x/wasm"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	bandwidthtypes "github.com/cybercongress/go-cyber/x/bandwidth/types"
	dmntypes "github.com/cybercongress/go-cyber/x/dmn/types"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
	gridtypes "github.com/cybercongress/go-cyber/x/grid/types"
	ranktypes "github.com/cybercongress/go-cyber/x/rank/types"
)

type WasmQuerierInterface interface {
	Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

type Querier struct {
	Queriers map[string]WasmQuerierInterface
}

func NewQuerier() Querier {
	return Querier{
		Queriers: make(map[string]WasmQuerierInterface),
	}
}

type WasmCustomQuery struct {
	Route     string          `json:"route"`
	QueryData json.RawMessage `json:"query_data"`
}

const (
	WasmQueryRouteGraph     = graphtypes.ModuleName
	WasmQueryRouteRank      = ranktypes.ModuleName
	WasmQueryRouteDmn       = dmntypes.ModuleName
	WasmQueryRouteGrid      = gridtypes.ModuleName
	WasmQueryRouteBandwidth = bandwidthtypes.ModuleName
	WasmQueryRouteLiquidity = liquiditytypes.ModuleName
)

func (q Querier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var customQuery WasmCustomQuery
	err := json.Unmarshal(data, &customQuery)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if querier, ok := q.Queriers[customQuery.Route]; ok {
		return querier.QueryCustom(ctx, customQuery.QueryData)
	} else {
		return nil, sdkerrors.Wrap(wasm.ErrQueryFailed, customQuery.Route)
	}
}
