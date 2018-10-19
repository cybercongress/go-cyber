package proxy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cybercongress/cyberd/cosmos/poc/app"
)

func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	sdk.RegisterWire(cdc)
	bank.RegisterWire(cdc)
	ibc.RegisterWire(cdc)
	auth.RegisterWire(cdc)
	cdc.RegisterConcrete(app.MsgLink{}, "cyberd/Link", nil)

	cdc.Seal()
	return cdc
}
