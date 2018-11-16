package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func SetPrefix() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("cbd", "cbd")
	config.SetBech32PrefixForValidator("cbd", "cbd")
	config.SetBech32PrefixForConsensusNode("cbd", "cbd")
	config.Seal()
}
