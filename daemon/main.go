package main

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/app"
	"github.com/cybercongress/cyberd/daemon/cmd"
	"github.com/cybercongress/cyberd/daemon/rpc"
	"github.com/cybercongress/cyberd/util"
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
	flagNotToSealAccPrefix    = "not-to-seal-acc-prefix"
)

func main() {

	rootDir := os.ExpandEnv("$HOME/.cyberd")

	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "cyberd",
		Short:             "Cyberd Daemon (server)",
		PersistentPreRunE: util.ConcatCobraCmdFuncs(server.PersistentPreRunEFn(ctx), setAppPrefix),
	}

	rootCmd.AddCommand(cmd.InitCmd(ctx, cdc))
	rootCmd.AddCommand(cmd.GenerateAccountsCmd(ctx, cdc))
	rootCmd.AddCommand(cmd.CollectGenTxsCmd(ctx, cdc))
	rootCmd.AddCommand(cmd.TestnetFilesCmd(ctx, cdc))
	rootCmd.AddCommand(cmd.GenTxCmd(ctx, cdc))
	rootCmd.AddCommand(cmd.AddGenesisAccountCmd(ctx, cdc))
	rootCmd.AddCommand(cmd.GenesisCmds(ctx, cdc))
	rootCmd.AddCommand(client.NewCompletionCmd(rootCmd, true))
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	for _, c := range rootCmd.Commands() {
		if c.Use == "start" {
			c.Flags().Bool(flagGpuEnabled, true, "Run cyberd with cuda calculations")
			c.Flags().Bool(flagSearchRpcQueryEnabled, true, "Build index of links with ranks and allow to query search through RPC")
		}
	}

	executor := cli.PrepareBaseCmd(rootCmd, "CBD", rootDir)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
	// todo use constant here
	pruning := baseapp.SetPruning(sdk.NewPruningOptions(60*60*24, 0))
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

func setAppPrefix(_ *cobra.Command, args []string) error {
	for _, arg := range args {
		if arg == flagNotToSealAccPrefix {
			return nil
		}
	}
	app.SetPrefix()
	return nil
}
