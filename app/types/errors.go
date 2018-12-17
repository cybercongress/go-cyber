package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Base error codes
	CodeOK                 sdk.CodeType = 0
	CodeLinkAlreadyExist   sdk.CodeType = 1
	CodeInvalidCid         sdk.CodeType = 2
	CodeCidNotFound        sdk.CodeType = 3
	CodeNotEnoughBandwidth sdk.CodeType = 4

	// Code space
	CodespaceCbd sdk.CodespaceType = "cyberd"
)

func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case CodeInvalidCid:
		return "invalid cid"
	case CodeCidNotFound:
		return "cid not found"
	case CodeLinkAlreadyExist:
		return "link already exists"
	case CodeNotEnoughBandwidth:
		return "not enough bandwidth to make transaction"
	default:
		return fmt.Sprintf("unknown error: code %d", code)
	}
}

//----------------------------------------
// Error constructors

func ErrInvalidCid() sdk.Error {
	return newError(CodespaceCbd, CodeInvalidCid)
}

func ErrNotEnoughBandwidth() sdk.Error {
	return newError(CodespaceCbd, CodeNotEnoughBandwidth)
}

func ErrCidNotFound() sdk.Error {
	return newError(CodespaceCbd, CodeCidNotFound)
}

func newError(codespace sdk.CodespaceType, code sdk.CodeType) sdk.Error {
	msg := codeToDefaultMsg(code)
	return sdk.NewError(codespace, code, msg)
}
