package keys

import (
	"bufio"
	"bytes"
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
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	//"github.com/tendermint/tendermint/libs/cli"
	//"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/cyberd/util"
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

			privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), b)
			ethPubkey, _ := btcec.ParsePubKey(pubKey.SerializeUncompressed(), btcec.S256())

			var cbdPubKey [33]byte
			copy(cbdPubKey[:], ethPubkey.SerializeCompressed()[:])

			var cbdPribKey [32]byte
			copy(cbdPribKey[:], privKey.Serialize()[:])

			pkArmor := mintkey.EncryptArmorPrivKey(secp256k1.PrivKeySecp256k1(cbdPribKey), passphrase, string(keys.Secp256k1))

			buffer := bytes.NewBuffer(nil)

			buffer.Write([]byte{13, 173, 21, 61, 10})
			amino.EncodeString(buffer, args[0])
			buffer.Write([]byte{18, 38, 235, 90, 233, 135, 33})
			buffer.Write(cbdPubKey[:])
			buffer.Write([]byte{26})
			amino.EncodeString(buffer, pkArmor)

			bz := buffer.Bytes()
			bufRes := bytes.NewBuffer(nil)
			amino.EncodeUvarint(bufRes, uint64(len(bz)))
			bufRes.Write(bz)

			armorStr := mintkey.ArmorInfoBytes(bufRes.Bytes())
			// import armored bytesm
			return kb.Import(args[0], armorStr)
		},
	}

	return cmd
}
