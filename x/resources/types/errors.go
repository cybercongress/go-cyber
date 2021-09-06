package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrTimeLockCoins 	  = sdkerrors.Register(ModuleName, 2, "error time lock coins")
	ErrIssueCoins 		  = sdkerrors.Register(ModuleName, 3, "error issue coins")
	ErrMintCoins 		  = sdkerrors.Register(ModuleName, 4, "error mint coins")
	ErrBurnCoins 		  = sdkerrors.Register(ModuleName, 5, "error burn coins")
	ErrSendMintedCoins 	  = sdkerrors.Register(ModuleName, 6, "error send minted coins")
	ErrNotAvailablePeriod = sdkerrors.Register(ModuleName, 7, "not available period")
	ErrInvalidAccountType = sdkerrors.Register(ModuleName, 8, "receiver account type not supported")
	ErrAccountNotFound 	  = sdkerrors.Register(ModuleName, 9, "account not found")
	ErrResourceNotExist   = sdkerrors.Register(ModuleName, 10, "resource not exist")
	ErrFullSlots           = sdkerrors.Register(ModuleName, 11, "all slots are full")
	ErrSmallReturn         = sdkerrors.Register(ModuleName, 12, "small resource's return amount")
	ErrInvalidBaseResource = sdkerrors.Register(ModuleName, 13, "invalid base resource")
)
