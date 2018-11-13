package main

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cybercongress/cyberd/cosmos/poc/app"
	"github.com/cybercongress/cyberd/cosmos/poc/cyberd/rpc"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/p2p"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	gaiaInit "github.com/cosmos/cosmos-sdk/cmd/gaia/init"
)

const (
	flagClientHome = "home-client"
	flagAccsCount = "accs-count"
)

func main() {

	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "cyberd",
		Short:             "Cyberd Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	cyberdAppInit := server.AppInit{
		AppGenState: CyberdAppGenState,
	}

	rootCmd.AddCommand(InitCmd(ctx, cdc, cyberdAppInit))
	server.AddCommands(ctx, cdc, rootCmd, cyberdAppInit, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.cyberd")
	executor := cli.PrepareBaseCmd(rootCmd, "CBD", rootDir)

	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

func InitCmd(ctx *server.Context, cdc *codec.Codec, appInit server.AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {

			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
			}

			nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
			if err != nil {
				return err
			}
			nodeID := string(nodeKey.ID())

			pk := gaiaInit.ReadOrCreatePrivValidator(config.PrivValidatorFile())
			genTx, appMessage, validator, err := CyberdAppGenTx(cdc, pk)
			if err != nil {
				return err
			}

			appState, err := appInit.AppGenState(
				cdc, tmtypes.GenesisDoc{}, []json.RawMessage{genTx})
			if err != nil {
				return err
			}
			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			toPrint := struct {
				ChainID    string          `json:"chain_id"`
				NodeID     string          `json:"node_id"`
				AppMessage json.RawMessage `json:"app_message"`
			}{
				chainID,
				nodeID,
				appMessage,
			}
			out, err := codec.MarshalJSONIndent(cdc, toPrint)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "%s\n", string(out))
			return gaiaInit.ExportGenesisFile(config.GenesisFile(), chainID,
				[]tmtypes.GenesisValidator{validator}, appStateJSON)
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, app.DefaultCLIHome, "client's home directory")
	cmd.Flags().String(client.FlagChainID, "",
		"genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(client.FlagName, "", "validator's moniker")
	cmd.Flags().Int(flagAccsCount, 1, "Count of initial accounts")
	return cmd
}

func newApp(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
	cyberdApp := app.NewCyberdApp(logger, db, baseapp.SetPruning(viper.GetString("pruning")))
	rpc.SetCyberdApp(cyberdApp)
	return cyberdApp
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, storeTracer io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	capp := app.NewCyberdApp(logger, db)
	return capp.ExportAppStateAndValidators()
}
