package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	AccountAddressPrefix   = Bech32Prefix
	AccountPubKeyPrefix    = AccountAddressPrefix + "pub"
	ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
	ValidatorPubKeyPrefix  = AccountAddressPrefix + "valoperpub"
	ConsNodeAddressPrefix  = AccountAddressPrefix + "valcons"
	ConsNodePubKeyPrefix   = AccountAddressPrefix + "valconspub"
	CoinType               = 118
)

// SetConfig sets the configuration for the network using mainnet or testnet
func SetConfig() {
	config := sdk.GetConfig()
	config.SetCoinType(uint32(CoinType))
	config.SetPurpose(uint32(44))
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.Seal()
}
