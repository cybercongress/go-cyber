package main

import (
	"encoding/json"
	"github.com/cybercongress/cyberd/cosmos/poc/app"
	"github.com/spf13/pflag"
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
)

var (
	FlagAccsCount = "accs-count"
)

func main() {

	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "cyberd",
		Short:             "Cyberd Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	cyberdFlagSet := pflag.NewFlagSet("cyberd-init", pflag.ExitOnError)
	cyberdFlagSet.Int(FlagAccsCount, 1, "Count of initial accounts")

	cyberdAppInit := server.AppInit{
		FlagsAppGenState: cyberdFlagSet,
		AppGenState:      CyberdAppGenState,
		AppGenTx:         CyberdAppGenTx,
	}

	server.AddCommands(ctx, cdc, rootCmd, cyberdAppInit,
		server.ConstructAppCreator(newApp, "cyberd"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "cyberd"))

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.cyberd")
	executor := cli.PrepareBaseCmd(rootCmd, "CBD", rootDir)

	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
	return app.NewCyberdApp(logger, db, baseapp.SetPruning(viper.GetString("pruning")))
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, storeTracer io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	capp := app.NewCyberdApp(logger, db)
	return capp.ExportAppStateAndValidators()
}
