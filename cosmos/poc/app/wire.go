package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgLink{}, "cyberd/Link", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}

// MakeCodec creates a new wire codec and registers all the necessary types
// with the codec.
func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()

	wire.RegisterCrypto(cdc)
	sdk.RegisterWire(cdc)
	bank.RegisterWire(cdc)
	ibc.RegisterWire(cdc)
	auth.RegisterWire(cdc)
	RegisterWire(cdc)

	cdc.Seal()

	return cdc
}
