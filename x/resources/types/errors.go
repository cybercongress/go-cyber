package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrTimeLockCoins 	  = sdkerrors.Register(ModuleName, 2, "err time lock coins")
	ErrIssueCoins 		  = sdkerrors.Register(ModuleName, 3, "error issue coins")
	ErrMintCoins 		  = sdkerrors.Register(ModuleName, 4, "error mint coins")
	ErrBurnCoins 		  = sdkerrors.Register(ModuleName, 5, "error burn coins")
	ErrSendMintedCoins 	  = sdkerrors.Register(ModuleName, 6, "error send minted coins")
	ErrInvalidDenom 	  = sdkerrors.Register(ModuleName, 7, "invalid denom")
	ErrNotAvailableLength = sdkerrors.Register(ModuleName, 8, "not available length")
	ErrInvalidAccountType = sdkerrors.Register(ModuleName, 9, "receiver account type not supported")
	ErrAccountNotFound 	  = sdkerrors.Register(ModuleName, 10, "account not found")
	ErrResourceNotExist   = sdkerrors.Register(ModuleName, 11, "resource not exist")
	ErrInvalidResource    = sdkerrors.Register(ModuleName, 12, "invalid resource")
	ErrFullSlots  		  = sdkerrors.Register(ModuleName, 13, "full slots")
	ErrLessThanOne  	  = sdkerrors.Register(ModuleName, 14, "token amount less than one")
	ErrAmountNotValid 	  = sdkerrors.Register(ModuleName, 15, "amount not valid")
)
