package helpers

import (
	errorsmod "cosmossdk.io/errors"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const codespace = "juno-global"

var (
	ErrInvalidAddress            = sdkerrors.ErrInvalidAddress
	ErrContractNotRegistered     = errorsmod.Register(codespace, 1, "contract not registered")
	ErrContractAlreadyRegistered = errorsmod.Register(codespace, 2, "contract already registered")
	ErrContractNotAdmin          = errorsmod.Register(codespace, 3, "sender is not the contract admin")
	ErrContractNotCreator        = errorsmod.Register(codespace, 4, "sender is not the contract creator")
	ErrInvalidCWContract         = errorsmod.Register(codespace, 5, "invalid CosmWasm contract")
	ErrOutOfGas                  = errorsmod.Register(codespace, 6, "contract execution ran out of gas")
	ErrContractExecutionPanic    = errorsmod.Register(codespace, 7, "contract execution panicked")
)
