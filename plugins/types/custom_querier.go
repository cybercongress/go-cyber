package types

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
