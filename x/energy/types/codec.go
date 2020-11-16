package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateEnergyRoute{}, "cyber/MsgCreateEnergyRoute", nil)
	cdc.RegisterConcrete(MsgEditEnergyRoute{}, "cyber/MsgEditEnergyRoute", nil)
	cdc.RegisterConcrete(MsgDeleteEnergyRoute{}, "cyber/MsgDeleteEnergyRoute", nil)
	cdc.RegisterConcrete(MsgEditEnergyRouteAlias{}, "cyber/MsgEditEnergyRouteAlias", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
