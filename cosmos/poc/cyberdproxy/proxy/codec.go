package proxy

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cybercongress/cyberd/cosmos/poc/app"
)

func MakeCodec() *codec.Codec {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	ibc.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	cdc.RegisterConcrete(app.MsgLink{}, "cyberd/Link", nil)

	cdc.Seal()
	return cdc
}
