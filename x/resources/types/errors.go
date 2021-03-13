package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrTimeLockCoins = sdkerrors.Register(ModuleName, 2, "err time lock coins")
	ErrIssueCoins = sdkerrors.Register(ModuleName, 3, "error issue coins")
	ErrMintCoins = sdkerrors.Register(ModuleName, 4, "error mint coins")
	ErrSendMintedCoins = sdkerrors.Register(ModuleName, 5, "error send minted coins")
	ErrInvalidDenom = sdkerrors.Register(ModuleName, 6, "invalid denom")
	ErrNotAvailableLength = sdkerrors.Register(ModuleName, 7, "not available length")
	ErrInvalidAccountType = sdkerrors.Register(ModuleName, 8, "receiver account type not supported")
	ErrAccountNotFound = sdkerrors.Register(ModuleName, 9, "account not found")
	ErrNotDefined = sdkerrors.Register(ModuleName, 10, "error not defined")
	ErrResourceNotFound = sdkerrors.Register(ModuleName, 11, "resource not found")
	ErrNotAuthorized = sdkerrors.Register(ModuleName, 12, "not authorized")
	ErrInvalidResource  = sdkerrors.Register(ModuleName, 13, "invalid resource")

)
