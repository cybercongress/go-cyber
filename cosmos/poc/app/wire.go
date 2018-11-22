package app

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgLink{}, "cyberd/Link", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}

// MakeCodec creates a new  codec and registers all the necessary types
// with the codec.
func MakeCodec() *codec.Codec {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)

	sdk.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	ibc.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	RegisterCodec(cdc)

	cdc.Seal()
	return cdc
}

// GetAccountDecoder returns the AccountDecoder function for the custom
// AppAccount.
func GetAccountDecoder(cdc *codec.Codec) auth.AccountDecoder {
	return func(accBytes []byte) (auth.Account, error) {
		if len(accBytes) == 0 {
			return nil, sdk.ErrTxDecode("accBytes are empty")
		}

		acct := new(auth.BaseAccount)
		err := cdc.UnmarshalBinaryBare(accBytes, &acct)
		if err != nil {
			panic(err)
		}

		return acct, err
	}
}
