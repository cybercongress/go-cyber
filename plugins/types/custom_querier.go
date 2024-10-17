package types

import (
	errorsmod "cosmossdk.io/errors"
	"encoding/json"
	"errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bandwidthkeeper "github.com/cybercongress/go-cyber/v5/x/bandwidth/keeper"
	bandwidthtypes "github.com/cybercongress/go-cyber/v5/x/bandwidth/types"
	dmnkeeper "github.com/cybercongress/go-cyber/v5/x/dmn/keeper"
	dmntypes "github.com/cybercongress/go-cyber/v5/x/dmn/types"
	graphkeeper "github.com/cybercongress/go-cyber/v5/x/graph/keeper"
	graphtypes "github.com/cybercongress/go-cyber/v5/x/graph/types"
	gridkeeper "github.com/cybercongress/go-cyber/v5/x/grid/keeper"
	gridtypes "github.com/cybercongress/go-cyber/v5/x/grid/types"
	rankkeeper "github.com/cybercongress/go-cyber/v5/x/rank/keeper"
	ranktypes "github.com/cybercongress/go-cyber/v5/x/rank/types"
	tokenfactorykeeper "github.com/cybercongress/go-cyber/v5/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/cybercongress/go-cyber/v5/x/tokenfactory/types"
)

type ModuleQuerier interface {
	HandleQuery(ctx sdk.Context, query CyberQuery) ([]byte, error)
}

var ErrHandleQuery = errors.New("error handle query")

type QueryPlugin struct {
	moduleQueriers     map[string]ModuleQuerier
	rankKeeper         *rankkeeper.StateKeeper
	graphKeeper        *graphkeeper.GraphKeeper
	dmnKeeper          *dmnkeeper.Keeper
	gridKeeper         *gridkeeper.Keeper
	bandwidthMeter     *bandwidthkeeper.BandwidthMeter
	bankKeeper         *bankkeeper.Keeper
	tokenFactoryKeeper *tokenfactorykeeper.Keeper
}

func NewQueryPlugin(
	moduleQueriers map[string]ModuleQuerier,
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

		switch {
		case contractQuery.ParticleRank != nil:
			return qp.moduleQueriers[ranktypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.GraphStats != nil:
			return qp.moduleQueriers[graphtypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.Thought != nil:
			return qp.moduleQueriers[dmntypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.ThoughtStats != nil:
			return qp.moduleQueriers[dmntypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.ThoughtsFees != nil:
			return qp.moduleQueriers[dmntypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.SourceRoutes != nil:
			return qp.moduleQueriers[gridtypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.SourceRoutedEnergy != nil:
			return qp.moduleQueriers[gridtypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.DestinationRoutedEnergy != nil:
			return qp.moduleQueriers[gridtypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.Route != nil:
			return qp.moduleQueriers[gridtypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.BandwidthLoad != nil:
			return qp.moduleQueriers[bandwidthtypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.BandwidthPrice != nil:
			return qp.moduleQueriers[bandwidthtypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.TotalBandwidth != nil:
			return qp.moduleQueriers[bandwidthtypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.NeuronBandwidth != nil:
			return qp.moduleQueriers[bandwidthtypes.ModuleName].HandleQuery(ctx, contractQuery)
		//case contractQuery.TokenFactory != nil:
		//	return qp.moduleQueriers[tokenfactorytypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.FullDenom != nil:
			return qp.moduleQueriers[tokenfactorytypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.Admin != nil:
			return qp.moduleQueriers[tokenfactorytypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.Metadata != nil:
			return qp.moduleQueriers[tokenfactorytypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.DenomsByCreator != nil:
			return qp.moduleQueriers[tokenfactorytypes.ModuleName].HandleQuery(ctx, contractQuery)
		case contractQuery.Params != nil:
			return qp.moduleQueriers[tokenfactorytypes.ModuleName].HandleQuery(ctx, contractQuery)
		default:
			return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown cyber query variant"}
		}

		// Iterate over the module queriers and dispatch to the appropriate one
		//for _, querier := range qp.moduleQueriers {
		//	resp, err := querier.HandleQuery(ctx, contractQuery)
		//	if err != nil {
		//		if err == ErrHandleQuery {
		//			// This querier cannot handle the query, try the next one
		//			continue
		//		}
		//		// Some other error occurred, return it
		//		return nil, err
		//	}
		//	// Query was handled successfully, return the response
		//	return resp, nil
		//}
		//
		//// If no querier could handle the query, return an error
		//return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown cyber query variant"}
	}
}
