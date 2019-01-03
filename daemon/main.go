package main

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cybercongress/cyberd/app"
	"github.com/cybercongress/cyberd/daemon/genesis"
	initCyberd "github.com/cybercongress/cyberd/daemon/init"
	"github.com/cybercongress/cyberd/daemon/rpc"
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
	flagGpuEnabled            = "compute-rank-on-gpu"
	flagSearchRpcQueryEnabled = "allow-search-rpc-query"
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
	rootCmd.AddCommand(initCyberd.GenerateAccountsCmd(ctx, cdc))
	rootCmd.AddCommand(initCyberd.CollectGenTxsCmd(ctx, cdc))
	rootCmd.AddCommand(initCyberd.TestnetFilesCmd(ctx, cdc))
	rootCmd.AddCommand(initCyberd.GenTxCmd(ctx, cdc))
	rootCmd.AddCommand(initCyberd.AddGenesisAccountCmd(ctx, cdc))
	rootCmd.AddCommand(genesis.GenerateEulerGenesisFile(ctx, cdc))
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	for _, c := range rootCmd.Commands() {
		if c.Use == "start" {
			c.Flags().Bool(flagGpuEnabled, true, "Run cyberd with cuda calculations")
			c.Flags().Bool(flagSearchRpcQueryEnabled, true, "Build index of links with ranks and allow to query search through RPC")
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
	cyberdApp := app.NewCyberdApp(logger, db, computeUnit, viper.GetBool(flagSearchRpcQueryEnabled), pruning)
	rpc.SetCyberdApp(cyberdApp)
	return cyberdApp
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	capp := app.NewCyberdApp(logger, db, rank.CPU, false)
	return capp.ExportAppStateAndValidators()
}
