package init

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/go-bip39"
	"github.com/cybercongress/cyberd/app"
	. "github.com/cybercongress/cyberd/app/genesis"
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
)

// Used for development purposes.
// Takes about 4min to generated 250kk genesis accs.
func GenerateAccountsCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-random-genesis-accounts [count]",
		Short: "Add randoms accounts to genesis.json. Used to populate tests chain during performance testing.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {

			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			count, ok := sdk.NewIntFromString(args[0])
			if !ok {
				return fmt.Errorf("failed to parse accs count: %s", args[0])
			}

			genFile := config.GenesisFile()
			if !common.FileExists(genFile) {
				return fmt.Errorf("%s does not exist, run `cyberd init` first", genFile)
			}

			genDoc, err := loadGenesisDoc(cdc, genFile)
			if err != nil {
				return err
			}

			var appState GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			amount := int64(10)
			kb := client.MockKeyBase()
			addresses := make([]GenesisAccount, count.Int64())

			var wg sync.WaitGroup
			wg.Add(int(count.Int64()))
			for i := int64(0); i < count.Int64(); i++ {
				go func(position int64) {
					defer wg.Done()
					entropySeed, _ := bip39.NewEntropy(256)
					mnemonic, _ := bip39.NewMnemonic(entropySeed[:])
					info, _ := kb.CreateKey(string(position), mnemonic, "")
					addresses[position] = GenesisAccount{Address: info.GetAddress(), Amount: amount}
				}(i)
			}
			wg.Wait()

			appState.Accounts = append(appState.Accounts, addresses...)
			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			return ExportGenesisFile(genFile, genDoc.ChainID, nil, appStateJSON)
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	return cmd
}
