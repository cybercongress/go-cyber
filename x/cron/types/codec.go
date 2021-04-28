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
	cdc.RegisterConcrete(&MsgAddJob{}, "cyber/MsgAddJob", nil)
	cdc.RegisterConcrete(&MsgRemoveJob{}, "cyber/MsgRemoveJob", nil)
	cdc.RegisterConcrete(&MsgChangeJobCID{}, "cyber/MsgChangeJobCID", nil)
	cdc.RegisterConcrete(&MsgChangeJobLabel{}, "cyber/MsgChangeJobLabel", nil)
	cdc.RegisterConcrete(&MsgChangeJobCallData{}, "cyber/MsgChangeJobCallData", nil)
	cdc.RegisterConcrete(&MsgChangeJobGasPrice{}, "cyber/MsgChangeJobGasPrice", nil)
	cdc.RegisterConcrete(&MsgChangeJobPeriod{}, "cyber/MsgChangeJobPeriod", nil)
	cdc.RegisterConcrete(&MsgChangeJobBlock{}, "cyber/MsgChangeJobBlock", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddJob{},
		&MsgRemoveJob{},
		&MsgChangeJobCID{},
		&MsgChangeJobLabel{},
		&MsgChangeJobCallData{},
		&MsgChangeJobGasPrice{},
		&MsgChangeJobPeriod{},
		&MsgChangeJobBlock{},
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
