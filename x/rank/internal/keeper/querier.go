package keeper

import (
	//"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cybercongress/cyberd/x/rank/exported"
	"github.com/cybercongress/cyberd/x/rank/internal/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewQuerier returns a minting Querier handler. k exported.StateKeeper
func NewQuerier(k exported.StateKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, _ abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParams(ctx, k)

		case types.QueryCalculationWindow:
			return queryCalculationWindow(ctx, k)

		case types.QueryDampingFactor:
			return queryDampingFactor(ctx, k)

		case types.QueryTolerance:
			return queryTolerance(ctx, k)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, k exported.StateKeeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryCalculationWindow(ctx sdk.Context, k exported.StateKeeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params.CalculationPeriod)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDampingFactor(ctx sdk.Context, k exported.StateKeeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params.DampingFactor)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryTolerance(ctx sdk.Context, k exported.StateKeeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params.Tolerance)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
