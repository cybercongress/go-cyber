package main

import (
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/tendermint/go-amino"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cybercongress/cyberd/app"
	cyberdcmd "github.com/cybercongress/cyberd/cli/commands"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"

	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	distClient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	govClient "github.com/cosmos/cosmos-sdk/x/gov/client"
	slashingClient "github.com/cosmos/cosmos-sdk/x/slashing/client"
	stakeClient "github.com/cosmos/cosmos-sdk/x/stake/client"
)

const (
	storeAcc      = "acc"
	storeGov      = "gov"
	storeSlashing = "slashing"
	storeStake    = "stake"
	storeDist     = "distr"
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

	mc := []sdk.ModuleClients{
		govClient.NewModuleClient(storeGov, cdc),
		distClient.NewModuleClient(storeDist, cdc),
		stakeClient.NewModuleClient(storeStake, cdc),
		slashingClient.NewModuleClient(storeSlashing, cdc),
	}

	// Construct Root Command
	cyberdcli.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(),
		queryCmd(cdc, mc),
		txCmd(cdc, mc),
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
		cyberdcmd.ImportPrivateKeyCmd(),
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

func queryCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		tx.SearchTxCmd(cdc),
		tx.QueryTxCmd(cdc),
		client.LineBreak,
		authcmd.GetAccountCmd(storeAcc, cdc),
	)

	for _, m := range mc {
		queryCmd.AddCommand(m.GetQueryCmd())
	}

	return queryCmd
}

func txCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		client.LineBreak,
		authcmd.GetSignCommand(cdc),
		bankcmd.GetBroadcastCommand(cdc),
		client.LineBreak,
	)

	for _, m := range mc {
		txCmd.AddCommand(m.GetTxCmd())
	}

	return txCmd
}
