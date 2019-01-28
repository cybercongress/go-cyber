package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	AccountAddressPrefix   = "cyber"
	AccountPubKeyPrefix    = "cyberpub"
	ValidatorAddressPrefix = "cybervaloper"
	ValidatorPubKeyPrefix  = "cybervaloperpub"
	ConsNodeAddressPrefix  = "cybervalcons"
	ConsNodePubKeyPrefix   = "cybervalconspub"
)

func SetPrefix() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.Seal()
}
