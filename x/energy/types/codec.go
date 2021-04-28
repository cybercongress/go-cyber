package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateRoute{}, "cyber/MsgCreateRoute", nil)
	cdc.RegisterConcrete(&MsgEditRoute{}, "cyber/MsgEditRoute", nil)
	cdc.RegisterConcrete(&MsgDeleteRoute{}, "cyber/MsgDeleteRoute", nil)
	cdc.RegisterConcrete(&MsgEditRouteAlias{}, "cyber/MsgEditRouteAlias", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateRoute{},
		&MsgEditRoute{},
		&MsgDeleteRoute{},
		&MsgEditRouteAlias{},
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
