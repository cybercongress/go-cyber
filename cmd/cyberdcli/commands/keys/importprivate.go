package keys

import (
	"bufio"
	"encoding/hex"
	//"fmt"
	//"os"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"

	//"github.com/cosmos/cosmos-sdk/client"
	//"github.com/cosmos/cosmos-sdk/client/input"
	//"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/mintkey"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/btcd/btcec"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	//"github.com/tendermint/tendermint/libs/cli"
	//"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/util"
	//"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
)

const hashPrefix = "0x"

func importPrivateKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "private <name>",
		Short: "Import ethereum private key into the local keybase",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			buf := bufio.NewReader(cmd.InOrStdin())

			kb, err := keys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), buf)
			if err != nil {
				return err
			}

			passphrase, err := input.GetPassword("Enter passphrase to encrypt/decrypt your key:", buf)
			if err != nil {
				return err
			}
			privateKey, err := input.GetString(
				"Enter your private key (hex encoded):", buf)
			if err != nil {
				return err
			}


			if util.HasPrefixIgnoreCase(privateKey, hashPrefix) {
				privateKey = privateKey[len(hashPrefix):]
			}

			b, _ := hex.DecodeString(privateKey)

			privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), b)

			var cbdPribKey [32]byte
			copy(cbdPribKey[:], privKey.Serialize()[:])

			pkArmor := mintkey.EncryptArmorPrivKey(secp256k1.PrivKeySecp256k1(cbdPribKey), passphrase, string(keys.Secp256k1))

			return kb.ImportPrivKey(args[0], pkArmor, passphrase)
		},
	}

	return cmd
}
