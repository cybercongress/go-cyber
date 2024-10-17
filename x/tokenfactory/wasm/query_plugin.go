package wasm

import (
	"encoding/json"
	"fmt"
	pluginstypes "github.com/cybercongress/go-cyber/v5/plugins/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	bindingstypes "github.com/cybercongress/go-cyber/v5/x/tokenfactory/wasm/types"
)

// CustomQuerier dispatches custom CosmWasm wasm queries.
func (querier *Querier) HandleQuery(ctx sdk.Context, query pluginstypes.CyberQuery) ([]byte, error) {
	//var contractQuery = query.TokenFactory

	switch {
	case query.FullDenom != nil:
		creator := query.FullDenom.CreatorAddr
		subdenom := query.FullDenom.Subdenom

		fullDenom, err := GetFullDenom(creator, subdenom)
		if err != nil {
			return nil, errorsmod.Wrap(err, "osmo full denom query")
		}

		res := bindingstypes.FullDenomResponse{
			Denom: fullDenom,
		}

		bz, err := json.Marshal(res)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to marshal FullDenomResponse")
		}

		return bz, nil

	case query.Admin != nil:
		res, err := querier.GetDenomAdmin(ctx, query.Admin.Denom)
		if err != nil {
			return nil, err
		}

		bz, err := json.Marshal(res)
		if err != nil {
			return nil, fmt.Errorf("failed to JSON marshal AdminResponse: %w", err)
		}

		return bz, nil

	case query.Metadata != nil:
		res, err := querier.GetMetadata(ctx, query.Metadata.Denom)
		if err != nil {
			return nil, err
		}

		bz, err := json.Marshal(res)
		if err != nil {
			return nil, fmt.Errorf("failed to JSON marshal MetadataResponse: %w", err)
		}

		return bz, nil

	case query.DenomsByCreator != nil:
		res, err := querier.GetDenomsByCreator(ctx, query.DenomsByCreator.Creator)
		if err != nil {
			return nil, err
		}

		bz, err := json.Marshal(res)
		if err != nil {
			return nil, fmt.Errorf("failed to JSON marshal DenomsByCreatorResponse: %w", err)
		}

		return bz, nil

	case query.Params != nil:
		res, err := querier.GetParams(ctx)
		if err != nil {
			return nil, err
		}

		bz, err := json.Marshal(res)
		if err != nil {
			return nil, fmt.Errorf("failed to JSON marshal ParamsResponse: %w", err)
		}

		return bz, nil

	default:
		return nil, pluginstypes.ErrHandleQuery
	}
	//}
}

// ConvertSdkCoinsToWasmCoins converts sdk type coins to wasm vm type coins
func ConvertSdkCoinsToWasmCoins(coins []sdk.Coin) wasmvmtypes.Coins {
	var toSend wasmvmtypes.Coins
	for _, coin := range coins {
		c := ConvertSdkCoinToWasmCoin(coin)
		toSend = append(toSend, c)
	}
	return toSend
}

// ConvertSdkCoinToWasmCoin converts a sdk type coin to a wasm vm type coin
func ConvertSdkCoinToWasmCoin(coin sdk.Coin) wasmvmtypes.Coin {
	return wasmvmtypes.Coin{
		Denom: coin.Denom,
		// Note: tokenfactory tokens have 18 decimal places, so 10^22 is common, no longer in u64 range
		Amount: coin.Amount.String(),
	}
}
