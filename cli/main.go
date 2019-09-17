package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/cybercongress/cyberd/app"
	cyberdcmd "github.com/cybercongress/cyberd/cli/commands"
	"github.com/cybercongress/cyberd/cli/commands/keys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	"os"

)

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := app.MakeCodec()
	app.SetPrefix()

	cyberdcli := &cobra.Command{
		Use:   "cyberdcli",
		Short: "Command Line Interface for interacting with cyberd",
	}

	// todo: hack till we don't handle with all merkle proofs
	viper.SetDefault(client.FlagTrustNode, true)

	cyberdcli.PersistentFlags().String(client.FlagChainID, "", "Chain Id of cyberd node")

	// Construct Root Command
	cyberdcli.AddCommand(
		rpc.StatusCommand(),
		queryCmd(cdc),
		txCmd(cdc),
		client.LineBreak,
		keys.Commands(),
		client.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		client.LineBreak,
		version.Cmd,
		client.NewCompletionCmd(cyberdcli, true),
	)

	cyberdcli.AddCommand(
		client.PostCommands(
			cyberdcmd.LinkTxCmd(cdc),
		)...)

	executor := cli.PrepareMainCmd(cyberdcli, "CBD", os.ExpandEnv("$HOME/.cyberdcli"))

	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(cdc),
		client.LineBreak,
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcmd.QueryTxCmd(cdc),
		authcmd.QueryTxsByEventsCmd(cdc),
		client.LineBreak,
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		client.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		client.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		client.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	// remove auth and bank commands as they're mounted under the root tx command
	var cmdsToRemove []*cobra.Command

	for _, cmd := range txCmd.Commands() {
		if cmd.Use == auth.ModuleName || cmd.Use == bank.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	txCmd.RemoveCommand(cmdsToRemove...)

	return txCmd
}

func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}
