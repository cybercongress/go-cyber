package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/version"
	at "github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	dist "github.com/cosmos/cosmos-sdk/x/distribution/client/rest"
	"github.com/cybercongress/cyberd/app"
	cyberdcmd "github.com/cybercongress/cyberd/cli/commands"
	"github.com/cybercongress/cyberd/cli/commands/keys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	"os"

	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bankrest "github.com/cosmos/cosmos-sdk/x/bank/client/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	slashingrest "github.com/cosmos/cosmos-sdk/x/slashing/client/rest"
	stakingrest "github.com/cosmos/cosmos-sdk/x/staking/client/rest"

	distcmd "github.com/cosmos/cosmos-sdk/x/distribution"
)

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := app.MakeCodec()
	app.SetPrefix()

	cyberdcli := &cobra.Command{
		Use:   "cyberdcli",
		Short: "Command line interface for interacting with cyberd",
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
		// Note: Handle with #870
		panic(err)
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
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		client.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	return txCmd
}

func registerRoutes(rs *lcd.RestServer) {
	authrest.RegisterRoutes(rs.CliCtx, rs.Mux, at.StoreKey)
	bankrest.RegisterRoutes(rs.CliCtx, rs.Mux)
	dist.RegisterRoutes(rs.CliCtx, rs.Mux, distcmd.StoreKey)
	stakingrest.RegisterRoutes(rs.CliCtx, rs.Mux)
	slashingrest.RegisterRoutes(rs.CliCtx, rs.Mux)
	phs := make([]govrest.ProposalRESTHandler, 0)
	govrest.RegisterRoutes(rs.CliCtx, rs.Mux, phs)
}
