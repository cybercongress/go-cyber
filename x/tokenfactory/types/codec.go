package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// TODO revisit naming
	cdc.RegisterConcrete(&MsgCreateDenom{}, "cyber/tokenfactory/MsgCreateDenom", nil)
	cdc.RegisterConcrete(&MsgMint{}, "cyber/tokenfactory/MsgMint", nil)
	cdc.RegisterConcrete(&MsgBurn{}, "cyber/tokenfactory/MsgBurn", nil)
	cdc.RegisterConcrete(&MsgChangeAdmin{}, "cyber/tokenfactory/MsgChangeAdmin", nil)
	cdc.RegisterConcrete(&MsgSetDenomMetadata{}, "cyber/tokenfactory/MsgSetDenomMetadata", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateDenom{},
		&MsgMint{},
		&MsgBurn{},
		&MsgChangeAdmin{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterCodec(amino)
	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	// Note: these 3 are inlines from authz/codec in 0.46 so we can be compatible with 0.45
	sdk.RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	codec.RegisterEvidences(amino)

	amino.Seal()
}
