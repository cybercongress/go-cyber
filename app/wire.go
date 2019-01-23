package app

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cybercongress/cyberd/x/link"
)

// MakeCodec creates a new  codec and registers all the necessary types
// with the codec.
func MakeCodec() *codec.Codec {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)

	sdk.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	ibc.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	link.RegisterCodec(cdc)

	cdc.Seal()
	return cdc
}
