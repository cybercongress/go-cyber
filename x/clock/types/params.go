package types

import (
	errorsmod "cosmossdk.io/errors"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefaultParams returns default parameters
func DefaultParams() Params {
	return Params{
		ContractGasLimit: 100_000,
	}
}

// NewParams creates a new Params object
func NewParams(
	contractGasLimit uint64,
) Params {
	return Params{
		ContractGasLimit: contractGasLimit,
	}
}

// Validate performs basic validation.
func (p Params) Validate() error {
	minimumGas := uint64(100_000)
	if p.ContractGasLimit < minimumGas {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"invalid contract gas limit: %d. Must be above %d", p.ContractGasLimit, minimumGas,
		)
	}

	return nil
}
