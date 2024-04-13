package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrTimeLockCoins       = errorsmod.Register(ModuleName, 2, "error timelock coins")
	ErrIssueCoins          = errorsmod.Register(ModuleName, 3, "error issue coins")
	ErrMintCoins           = errorsmod.Register(ModuleName, 4, "error mint coins")
	ErrBurnCoins           = errorsmod.Register(ModuleName, 5, "error burn coins")
	ErrSendMintedCoins     = errorsmod.Register(ModuleName, 6, "error send minted coins")
	ErrNotAvailablePeriod  = errorsmod.Register(ModuleName, 7, "period not available")
	ErrInvalidAccountType  = errorsmod.Register(ModuleName, 8, "receiver account type not supported")
	ErrAccountNotFound     = errorsmod.Register(ModuleName, 9, "account not found")
	ErrResourceNotExist    = errorsmod.Register(ModuleName, 10, "resource does not exist")
	ErrFullSlots           = errorsmod.Register(ModuleName, 11, "all slots are full")
	ErrSmallReturn         = errorsmod.Register(ModuleName, 12, "insufficient resources return amount")
	ErrInvalidBaseResource = errorsmod.Register(ModuleName, 13, "invalid base resource")
)
