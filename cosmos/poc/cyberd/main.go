package main

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cybercongress/cyberd/cosmos/poc/app"
	"github.com/cybercongress/cyberd/cosmos/poc/app/rank"
	"github.com/cybercongress/cyberd/cosmos/poc/cyberd/rpc"
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
	flagClientHome = "home-client"
	flagAccsCount  = "accs-count"
	flagGpuEnabled = "compute-rank-on-gpu"
)

func main() {

	cdc := app.MakeCodec()
	app.SetPrefix()
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

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, storeTracer io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	capp := app.NewCyberdApp(logger, db, rank.CPU)
	return capp.ExportAppStateAndValidators()
}
