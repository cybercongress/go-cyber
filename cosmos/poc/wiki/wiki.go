package main

import (
	"fmt"
	"github.com/cybercongress/cyberd/cosmos/poc/claim/common"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

type Link struct {
	from string
	to   string
}

var (
	wiki = &cobra.Command{
		Use:   "wiki",
		Short: "Use to index wiki",
	}
)

/*
./wiki start
--node=127.0.0.1:34657 --passphrase=1q2w3e4r \
--address=cosmos1g7e74lxpwlcsza8v0nca2hwahluqcv4r4d3p8p \
--chain-id=test-chain-gRXWCL
*/

func main() {

	wiki.AddCommand(StartCmd())
	homeDir, err := getHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.SetDefault("home", homeDir+"/.cyberdwiki")

	if err := wiki.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func StartCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start wiki indexing",
		RunE: func(cmd *cobra.Command, args []string) error {
			Index()
			return nil
		},
	}

	cmd.Flags().String(common.FlagPassphrase, "", "Passphrase of account to claim from")
	cmd.Flags().String(common.FlagChainId, "", "Chain Id")
	cmd.Flags().String(common.FlagAddress, "", "ClaimFrom of account to claim from")
	cmd.Flags().String(common.FlagNode, "127.0.0.1:26657", "Node url to connect")

	viper.BindPFlag(common.FlagPassphrase, cmd.Flags().Lookup(common.FlagPassphrase))
	viper.BindPFlag(common.FlagChainId, cmd.Flags().Lookup(common.FlagChainId))
	viper.BindPFlag(common.FlagAddress, cmd.Flags().Lookup(common.FlagAddress))
	viper.BindPFlag(common.FlagNode, cmd.Flags().Lookup(common.FlagNode))
	return cmd
}

func getHomeDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return home, nil
}
