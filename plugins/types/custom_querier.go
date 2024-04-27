package types

import (
	"encoding/json"
	"errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bandwidthkeeper "github.com/cybercongress/go-cyber/v4/x/bandwidth/keeper"
	dmnkeeper "github.com/cybercongress/go-cyber/v4/x/dmn/keeper"
	graphkeeper "github.com/cybercongress/go-cyber/v4/x/graph/keeper"
	gridkeeper "github.com/cybercongress/go-cyber/v4/x/grid/keeper"
	rankkeeper "github.com/cybercongress/go-cyber/v4/x/rank/keeper"
	tokenfactorykeeper "github.com/cybercongress/go-cyber/v4/x/tokenfactory/keeper"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ModuleQuerier interface {
	HandleQuery(ctx sdk.Context, query CyberQuery) ([]byte, error)
}

var ErrHandleQuery = errors.New("error handle query")

type QueryPlugin struct {
	moduleQueriers     []ModuleQuerier
	rankKeeper         *rankkeeper.StateKeeper
	graphKeeper        *graphkeeper.GraphKeeper
	dmnKeeper          *dmnkeeper.Keeper
	gridKeeper         *gridkeeper.Keeper
	bandwidthMeter     *bandwidthkeeper.BandwidthMeter
	bankKeeper         *bankkeeper.Keeper
	tokenFactoryKeeper *tokenfactorykeeper.Keeper
}

func NewQueryPlugin(
	moduleQueriers []ModuleQuerier,
	rank *rankkeeper.StateKeeper,
	graph *graphkeeper.GraphKeeper,
	dmn *dmnkeeper.Keeper,
	grid *gridkeeper.Keeper,
	bandwidth *bandwidthkeeper.BandwidthMeter,
	bank *bankkeeper.Keeper,
	tokenFactory *tokenfactorykeeper.Keeper,
) *QueryPlugin {
	return &QueryPlugin{
		moduleQueriers:     moduleQueriers,
		rankKeeper:         rank,
		graphKeeper:        graph,
		dmnKeeper:          dmn,
		gridKeeper:         grid,
		bandwidthMeter:     bandwidth,
		bankKeeper:         bank,
		tokenFactoryKeeper: tokenFactory,
	}
}

func CustomQuerier(qp *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var contractQuery CyberQuery
		if err := json.Unmarshal(request, &contractQuery); err != nil {
			return nil, errorsmod.Wrap(err, "cyber query error")
		}

		// Iterate over the module queriers and dispatch to the appropriate one
		for _, querier := range qp.moduleQueriers {
			resp, err := querier.HandleQuery(ctx, contractQuery)
			if err != nil {
				if err == ErrHandleQuery {
					// This querier cannot handle the query, try the next one
					continue
				}
				// Some other error occurred, return it
				return nil, err
			}
			// Query was handled successfully, return the response
			return resp, nil
		}

		// If no querier could handle the query, return an error
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown cyber query variant"}
	}
}
