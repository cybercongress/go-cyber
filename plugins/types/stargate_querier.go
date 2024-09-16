package types

import (
	"fmt"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StargateQuerier dispatches whitelisted stargate queries
func StargateQuerier(queryRouter baseapp.GRPCQueryRouter, cdc codec.Codec) func(ctx sdk.Context, request *wasmvmtypes.StargateQuery) ([]byte, error) {
	return func(ctx sdk.Context, request *wasmvmtypes.StargateQuery) ([]byte, error) {
		protoResponseType, err := GetWhitelistedQuery(request.Path)
		if err != nil {
			return nil, err
		}

		// no matter what happens after this point, we must return
		// the response type to prevent sync.Pool from leaking.
		defer ReturnStargateResponseToPool(request.Path, protoResponseType)

		route := queryRouter.Route(request.Path)
		if route == nil {
			return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("No route to query '%s'", request.Path)}
		}

		res, err := route(ctx, abci.RequestQuery{
			Data: request.Data,
			Path: request.Path,
		})
		if err != nil {
			return nil, err
		}

		if res.Value == nil {
			return nil, fmt.Errorf("Res returned from abci query route is nil")
		}

		bz, err := ConvertProtoToJSONMarshal(protoResponseType, res.Value, cdc)
		if err != nil {
			return nil, err
		}

		return bz, nil
	}
}

// ConvertProtoToJSONMarshal unmarshals the given bytes into a proto message and then marshals it to json.
// This is done so that clients calling stargate queries do not need to define their own proto unmarshalers,
// being able to use response directly by json marshalling, which is supported in cosmwasm.
func ConvertProtoToJSONMarshal(protoResponseType codec.ProtoMarshaler, bz []byte, cdc codec.Codec) ([]byte, error) {
	// unmarshal binary into stargate response data structure
	err := cdc.Unmarshal(bz, protoResponseType)
	if err != nil {
		return nil, wasmvmtypes.Unknown{}
	}

	bz, err = cdc.MarshalJSON(protoResponseType)
	if err != nil {
		return nil, wasmvmtypes.Unknown{}
	}

	protoResponseType.Reset()

	return bz, nil
}
