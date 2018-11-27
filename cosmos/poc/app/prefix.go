package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	AccountAddressPrefix   = "cbd"
	AccountPubKeyPrefix    = "cbdpub"
	ValidatorAddressPrefix = "cbdvaloper"
	ValidatorPubKeyPrefix  = "cbdvaloperpub"
	ConsNodeAddressPrefix  = "cbdvalcons"
	ConsNodePubKeyPrefix   = "cbdvalconspub"
)

func SetPrefix() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.Seal()
}
