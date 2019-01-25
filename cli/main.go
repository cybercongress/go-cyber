package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	at "github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	dist "github.com/cosmos/cosmos-sdk/x/distribution"
	distClient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	gv "github.com/cosmos/cosmos-sdk/x/gov"
	govClient "github.com/cosmos/cosmos-sdk/x/gov/client"
	sl "github.com/cosmos/cosmos-sdk/x/slashing"
	slashingClient "github.com/cosmos/cosmos-sdk/x/slashing/client"
	st "github.com/cosmos/cosmos-sdk/x/staking"
	stakingClient "github.com/cosmos/cosmos-sdk/x/staking/client"
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
		Short: "Command line interface for interacting with cyberd",
	}

	mc := []sdk.ModuleClients{
		govClient.NewModuleClient(gv.StoreKey, cdc),
		distClient.NewModuleClient(dist.StoreKey, cdc),
		stakingClient.NewModuleClient(st.StoreKey, cdc),
		slashingClient.NewModuleClient(sl.StoreKey, cdc),
	}

	// todo: hack till we don't handle with all merkle proofs
	viper.SetDefault(client.FlagTrustNode, true)
	cyberdcli.PersistentFlags().String(client.FlagChainID, "", "Chain Id of cyberd node")

	// Construct Root Command
	cyberdcli.AddCommand(
		rpc.StatusCommand(),
		queryCmd(cdc, mc),
		txCmd(cdc, mc),
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
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
		authcmd.GetAccountCmd(at.StoreKey, cdc),
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
		authcmd.GetMultiSignCommand(cdc),
		bankcmd.GetBroadcastCommand(cdc),
		client.LineBreak,
	)

	for _, m := range mc {
		txCmd.AddCommand(m.GetTxCmd())
	}

	return txCmd
}
