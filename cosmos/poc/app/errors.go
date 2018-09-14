package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Cyberd errors reserve 4200 ~ 4299.
const (
	DefaultCodespace sdk.CodespaceType = 2

	CodeInvalidCid  sdk.CodeType = 4201
)

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case CodeInvalidCid:
		return "invalid cid"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

//----------------------------------------
// Error constructors

func ErrInvalidCid(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidCid, "")
}

//----------------------------------------

func msgOrDefaultMsg(msg string, code sdk.CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func newError(codespace sdk.CodespaceType, code sdk.CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(codespace, code, msg)
}

