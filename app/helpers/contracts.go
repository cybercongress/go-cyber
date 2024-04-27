package helpers

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Execute contract, recover from panic
func ExecuteContract(k wasmtypes.ContractOpsKeeper, childCtx sdk.Context, contractAddr sdk.AccAddress, msgBz []byte, err *error) {
	// Recover from panic, return error
	defer func() {
		if recoveryError := recover(); recoveryError != nil {
			// Determine error associated with panic
			if isOutofGas, msg := IsOutOfGasError(recoveryError); isOutofGas {
				*err = ErrOutOfGas.Wrapf("%s", msg)
			} else {
				*err = ErrContractExecutionPanic.Wrapf("%s", recoveryError)
			}
		}
	}()

	// Execute contract with sudo
	_, *err = k.Sudo(childCtx, contractAddr, msgBz)
}

// Check if error is out of gas error
func IsOutOfGasError(err any) (bool, string) {
	switch e := err.(type) {
	case storetypes.ErrorOutOfGas:
		return true, e.Descriptor
	case storetypes.ErrorGasOverflow:
		return true, e.Descriptor
	default:
		return false, ""
	}
}
