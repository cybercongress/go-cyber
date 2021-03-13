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
	cdc.RegisterConcrete(&MsgCronAddJob{}, "cyber/MsgCronAddJob", nil)
	cdc.RegisterConcrete(&MsgCronRemoveJob{}, "cyber/MsgCronRemoveJob", nil)
	cdc.RegisterConcrete(&MsgCronChangeJobCID{}, "cyber/MsgCronChangeJobCID", nil)
	cdc.RegisterConcrete(&MsgCronChangeJobLabel{}, "cyber/MsgCronChangeJobLabel", nil)
	cdc.RegisterConcrete(&MsgCronChangeJobCallData{}, "cyber/MsgCronChangeJobCallData", nil)
	cdc.RegisterConcrete(&MsgCronChangeJobGasPrice{}, "cyber/MsgCronChangeJobGasPrice", nil)
	cdc.RegisterConcrete(&MsgCronChangeJobPeriod{}, "cyber/MsgCronChangeJobPeriod", nil)
	cdc.RegisterConcrete(&MsgCronChangeJobBlock{}, "cyber/MsgCronChangeJobBlock", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCronAddJob{},
		&MsgCronRemoveJob{},
		&MsgCronChangeJobCID{},
		&MsgCronChangeJobLabel{},
		&MsgCronChangeJobCallData{},
		&MsgCronChangeJobGasPrice{},
		&MsgCronChangeJobPeriod{},
		&MsgCronChangeJobBlock{},
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
