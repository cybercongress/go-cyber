package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var msgCdc = codec.New()

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCyberlink{}, "cyber/Link", nil)
}

func init() {
	RegisterCodec(msgCdc)
}
