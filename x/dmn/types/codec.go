package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateThought{}, "cyber/MsgCreateThought", nil)
	cdc.RegisterConcrete(&MsgForgetThought{}, "cyber/MsgForgetThought", nil)
	cdc.RegisterConcrete(&MsgChangeThoughtParticle{}, "cyber/MsgChangeThoughtParticle", nil)
	cdc.RegisterConcrete(&MsgChangeThoughtName{}, "cyber/MsgChangeThoughtName", nil)
	cdc.RegisterConcrete(&MsgChangeThoughtCallData{}, "cyber/MsgChangeThoughtCallData", nil)
	cdc.RegisterConcrete(&MsgChangeThoughtGasPrice{}, "cyber/MsgChangeThoughtGasPrice", nil)
	cdc.RegisterConcrete(&MsgChangeThoughtPeriod{}, "cyber/MsgChangeThoughtPeriod", nil)
	cdc.RegisterConcrete(&MsgChangeThoughtBlock{}, "cyber/MsgChangeThoughtBlock", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateThought{},
		&MsgForgetThought{},
		&MsgChangeThoughtParticle{},
		&MsgChangeThoughtName{},
		&MsgChangeThoughtCallData{},
		&MsgChangeThoughtGasPrice{},
		&MsgChangeThoughtPeriod{},
		&MsgChangeThoughtBlock{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
