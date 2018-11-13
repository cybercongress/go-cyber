package main

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"
)

type CyberdAppState struct {
	Accounts []CyberdInitAccount `json:"accounts"`
}

type CyberdCoin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type CyberdInitAccount struct {
	Address string       `json:"address"`
	Coins   []CyberdCoin `json:"coins"`
}

type CyberdGenTx struct {
	Addresses []sdk.AccAddress `json:"addresses"`
}

func CyberdAppGenTx(cdc *codec.Codec, pk crypto.PubKey) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {

	accsCount := viper.GetInt(flagAccsCount)

	addresses := make([]sdk.AccAddress, 0, accsCount)
	secrets := make(map[string]string)
	for i := 0; i < accsCount; i++ {
		var addr sdk.AccAddress
		var secret string
		addr, secret, err = server.GenerateCoinKey()
		if err != nil {
			return
		}
		addresses = append(addresses, addr)
		secrets[addr.String()] = secret
	}

	var bz []byte
	simpleGenTx := CyberdGenTx{addresses}
	bz, err = cdc.MarshalJSON(simpleGenTx)
	if err != nil {
		return
	}
	appGenTx = json.RawMessage(bz)

	bz, err = cdc.MarshalJSON(secrets)
	if err != nil {
		return
	}
	cliPrint = json.RawMessage(bz)

	validator = tmtypes.GenesisValidator{
		PubKey: pk,
		Power:  10,
	}
	return
}

// create the genesis app state
func CyberdAppGenState(cdc *codec.Codec, genDoc tmtypes.GenesisDoc, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {

	if len(appGenTxs) != 1 {
		err = errors.New("must provide a single genesis transaction")
		return
	}

	var genTx CyberdGenTx
	err = cdc.UnmarshalJSON(appGenTxs[0], &genTx)
	if err != nil {
		return
	}

	accounts := make([]CyberdInitAccount, 0, len(genTx.Addresses))
	for _, addr := range genTx.Addresses {
		accounts = append(accounts, CyberdInitAccount{addr.String(), []CyberdCoin{{"CBD", "9007199254740992"}}})
	}

	appState, err = json.Marshal(CyberdAppState{accounts})
	return
}
