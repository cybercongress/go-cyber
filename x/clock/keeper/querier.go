package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	globalerrors "github.com/cybercongress/go-cyber/v4/app/helpers"
	"github.com/cybercongress/go-cyber/v4/x/clock/types"
)

var _ types.QueryServer = &Querier{}

type Querier struct {
	keeper Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{
		keeper: k,
	}
}

// ContractModules returns contract addresses which are using the clock
func (q Querier) ClockContracts(stdCtx context.Context, req *types.QueryClockContracts) (*types.QueryClockContractsResponse, error) {
	ctx := sdk.UnwrapSDKContext(stdCtx)

	contracts, err := q.keeper.GetPaginatedContracts(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	return contracts, nil
}

// ClockContract returns the clock contract information
func (q Querier) ClockContract(stdCtx context.Context, req *types.QueryClockContract) (*types.QueryClockContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(stdCtx)

	// Ensure the contract address is valid
	if _, err := sdk.AccAddressFromBech32(req.ContractAddress); err != nil {
		return nil, globalerrors.ErrInvalidAddress
	}

	contract, err := q.keeper.GetClockContract(ctx, req.ContractAddress)
	if err != nil {
		return nil, err
	}

	return &types.QueryClockContractResponse{
		ClockContract: *contract,
	}, nil
}

// Params returns the total set of clock parameters.
func (q Querier) Params(stdCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(stdCtx)

	p := q.keeper.GetParams(ctx)

	return &types.QueryParamsResponse{
		Params: &p,
	}, nil
}
