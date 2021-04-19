package app

import (
	"github.com/cosmos/cosmos-sdk/std"

	"github.com/cybercongress/go-cyber/app/params"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
)

func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()

	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}

func MakeTestEncodingConfig() simappparams.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()

	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return simappparams.EncodingConfig(encodingConfig)
}
