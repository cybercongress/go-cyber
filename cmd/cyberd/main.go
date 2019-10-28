package main

import (
	"encoding/json"
	"github.com/cybercongress/cyberd/cmd/cyberd/rpc"

	//"github.com/cybercongress/cyberd/cyberd/rpc"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cybercongress/cyberd/app"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	genaccscli "github.com/cosmos/cosmos-sdk/x/genaccounts/client/cli"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"

	//"github.com/cybercongress/cyberd/cyberd/cmd"
	//"github.com/cybercongress/cyberd/cyberd/rpc"
	"github.com/cybercongress/cyberd/x/rank"
)

const (
	//flagGpuEnabled                = "compute-rank-on-gpu"
	//flagSearchEnabled             = "allow-search"
	flagInvCheckPeriod            = "inv-check-period"
)

var invCheckPeriod uint

func main() {

	//rootDir := os.ExpandEnv("$HOME/.cyberd")

	cdc := app.MakeCodec()

	app.SetPrefix()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "cyberd",
		Short:             "Cyberd Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.MigrateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(genutilcli.GenTxCmd(ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{},
		genaccounts.AppModuleBasic{}, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics))
	rootCmd.AddCommand(genaccscli.AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(client.NewCompletionCmd(rootCmd, true))
	rootCmd.AddCommand(testnetCmd(ctx, cdc, app.ModuleBasics, genaccounts.AppModuleBasic{}))
	//rootCmd.AddCommand(replayCmd())

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	executor := cli.PrepareBaseCmd(rootCmd, "CBD", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	//computeUnit := rank.CPU
	//if viper.GetBool(flagGpuEnabled) {
	//	computeUnit = rank.GPU
	//}
	//gpu := viper.GetBool(flagGpuEnabled)
	//if (gpu) {
	//
	//}

	cyberdApp := app.NewCyberdApp(
		logger, db, traceStore, int64(-1), invCheckPeriod,
		//computeUnit,
		rank.CPU,
		false,
		//viper.GetBool(flagSearchEnabled),
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
		baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
	)
	rpc.SetCyberdApp(cyberdApp) // ?
	return cyberdApp
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		// debug here
		capp := app.NewCyberdApp(logger, db, traceStore, height, uint(1), rank.CPU, false)
		return capp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	capp := app.NewCyberdApp(logger, db, traceStore, height, uint(1), rank.CPU, false)
	return capp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

