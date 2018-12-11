package main

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cybercongress/cyberd/app"
	initCyberd "github.com/cybercongress/cyberd/cyberd/init"
	"github.com/cybercongress/cyberd/cyberd/rpc"
	"github.com/cybercongress/cyberd/x/rank"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"io"
	_ "net/http/pprof"
	"os"
)

const (
	flagGpuEnabled = "compute-rank-on-gpu"
)

func main() {

	cdc := app.MakeCodec()
	app.SetPrefix()
	ctx := server.NewDefaultContext()

	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "cyberd",
		Short:             "Cyberd Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(initCyberd.InitCmd(ctx, cdc))
	rootCmd.AddCommand(initCyberd.GenerateAccountCmd())
	rootCmd.AddCommand(initCyberd.CollectGenTxsCmd(ctx, cdc))
	rootCmd.AddCommand(initCyberd.TestnetFilesCmd(ctx, cdc))
	rootCmd.AddCommand(initCyberd.GenTxCmd(ctx, cdc))
	rootCmd.AddCommand(initCyberd.AddGenesisAccountCmd(ctx, cdc))
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	for _, c := range rootCmd.Commands() {
		if c.Use == "start" {
			c.Flags().Bool(flagGpuEnabled, true, "Run cyberd with cuda calculations")
		}
	}

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
	pruning := baseapp.SetPruning(viper.GetString("pruning"))
	computeUnit := rank.CPU
	if viper.GetBool(flagGpuEnabled) {
		computeUnit = rank.GPU
	}
	cyberdApp := app.NewCyberdApp(logger, db, computeUnit, pruning)
	rpc.SetCyberdApp(cyberdApp)
	return cyberdApp
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	capp := app.NewCyberdApp(logger, db, rank.CPU)
	return capp.ExportAppStateAndValidators()
}
