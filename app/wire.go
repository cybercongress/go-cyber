package app

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MakeCodec creates a new  codec and registers
// all the necessary types with the codec.
func MakeCodec() *codec.Codec {
	cdc := codec.New()
	ModuleBasics.RegisterCodec(cdc)

	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)

	return cdc.Seal()
}
