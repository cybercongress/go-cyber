package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(Msg{}, "cyber/Link", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
