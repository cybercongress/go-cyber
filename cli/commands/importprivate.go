package commands

import (
	"bytes"
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/mintkey"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/btcd/btcec"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/cli"
)

const (
	flagName = "name"
)

//info = kb.writeLocalKey(secp256k1.PrivKeySecp256k1(derivedPriv), name, passwd)

// LinkTxCmd will create a link tx and sign it with the given key.
func ImportPrivateKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import_private <name>",
		Short: "Import private key",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			rootDir := viper.GetString(cli.HomeFlag)
			kb, _ := keys.GetKeyBaseFromDirWithWritePerm(rootDir)
			//var privateKey [32]byte
			//
			//kbType := reflect.TypeOf(kb)
			//fmt.Println(kbType)
			//
			//kbValue := reflect.ValueOf(kb).Convert(kbType)
			//fmt.Println(kbValue.NumMethod())
			//
			//method := reflect.ValueOf(kb).Convert(kbType).MethodByName("writeLocalKey")
			//in := make([]reflect.Value, 3)
			//in[0] = reflect.ValueOf(secp256k1.PrivKeySecp256k1(privateKey))
			//in[1] = reflect.ValueOf("artur54321")
			//in[2] = reflect.ValueOf("12345668")
			//
			//resp := method.Call(in)
			//fmt.Println(resp)

			//cbd19msfav283ykk2yeedv2hjfwjrp3xfdg0we7vdj
			//cbdpub1addwnpepqthpt3zevh30gn3cfd5eelf3zuw2kxaqj948tjhu3ztrdervejwucsalyeq
			//[38 45 173 80 146 139 15 246 55 75 107 187 25 121 131 220 152 254 99 73 204 227 249 40 107 134 252 237 74 55 4 169]
			//[235 90 233 135 33 2 129 33 115 9 231 5 65 199 75 91 150 174 228 90 235 212 46 136 215 242 138 140 160 175 8 247 179 72 246 3 162 217]
			//armor, _ := kb.Export("artur2")
			//infoBytes, _ := mintkey.UnarmorInfoBytes(armor)
			//fmt.Println(infoBytes)

			name := args[0]

			bufStdin := client.BufferStdin()
			encryptPassword, err := client.GetCheckPassword(
				"Enter a passphrase to encrypt your key to disk:",
				"Repeat the passphrase:", bufStdin)

			privateKey, err := client.GetString(
				"Enter your private key (hex encoded):", bufStdin)
			if err != nil {
				return err
			}
			//pubKey := [33]byte{2, 129, 33, 115, 9, 231, 5, 65, 199, 75, 91, 150, 174, 228, 90, 235, 212, 46, 136, 215, 242, 138, 140, 160, 175, 8, 247, 179, 72, 246, 3, 162, 217}
			//privKey := [32]byte{38, 45, 173, 80, 146, 139, 15, 246, 55, 75, 107, 187, 25, 121, 131, 220, 152, 254, 99, 73, 204, 227, 249, 40, 107, 134, 252, 237, 74, 55, 4, 169}

			b, _ := hex.DecodeString(privateKey)

			privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), b)
			ethPubkey, _ := btcec.ParsePubKey(pubKey.SerializeUncompressed(), btcec.S256())

			var cbdPubKey [33]byte
			copy(cbdPubKey[:], ethPubkey.SerializeCompressed()[:])

			var cbdPribKey [32]byte
			copy(cbdPribKey[:], privKey.Serialize()[:])

			pkArmor := mintkey.EncryptArmorPrivKey(secp256k1.PrivKeySecp256k1(cbdPribKey), encryptPassword)

			buf := bytes.NewBuffer(nil)

			buf.Write([]byte{13, 173, 21, 61, 10})
			amino.EncodeString(buf, name)
			buf.Write([]byte{18, 38, 235, 90, 233, 135, 33})
			buf.Write(cbdPubKey[:])
			buf.Write([]byte{26})
			amino.EncodeString(buf, pkArmor)

			bz := buf.Bytes()
			bufRes := bytes.NewBuffer(nil)
			amino.EncodeUvarint(bufRes, uint64(len(bz)))
			bufRes.Write(bz)

			armorStr := mintkey.ArmorInfoBytes(bufRes.Bytes())
			// import armored bytesm
			kb.Import(name, armorStr)

			//pk, _ := kb.ExportPrivateKeyObject("artur", "1234qwer")
			//fmt.Println(pk)
			//pk, _ := kb.ExportPubKey("artur")
			//pkUn, _ := mintkey.UnarmorPubKeyBytes(pk)
			//fmt.Println(pkUn)
			//
			//pkArmor := mintkey.EncryptArmorPrivKey(pk, "1234qwer")
			//fmt.Println(pkArmor)

			//eager salon spare notice expire throw chaos hub dune poverty stuff trim ancient fame armor domain bring topple thing boost impulse hire best scorpion
			//2 129 ... 162 217
			//info, _ := kb.CreateKey("artur", "eager salon spare notice expire throw chaos hub dune poverty stuff trim ancient fame armor domain bring topple thing boost impulse hire best scorpion", "1234qwer")
			//fmt.Println(info)

			return nil
		},
	}

	cmd.Flags().String(flagName, "", "Name of key owner")

	return cmd
}
