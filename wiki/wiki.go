package main

import (
	"fmt"
	"github.com/cybercongress/cyberd/client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

/*
./wiki start
--node=127.0.0.1:34657 --passphrase=1q2w3e4r \
--address=cosmos1g7e74lxpwlcsza8v0nca2hwahluqcv4r4d3p8p \
*/

func main() {

	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start wiki indexing",
		RunE: func(cmd *cobra.Command, args []string) error {
			address := viper.GetString(client.FlagAddress)
			nodeUrl := viper.GetString(client.FlagNode)
			pasphrase := viper.GetString(client.FlagPassphrase)
			cyberdClient := client.NewHttpCyberdClient(nodeUrl, pasphrase, address)
			ContinueIndex(cyberdClient)
			return nil
		},
	}

	cmd.Flags().String(client.FlagAddress, "", "Account to sign transactions")
	cmd.Flags().String(client.FlagPassphrase, "", "Passphrase of account")
	cmd.Flags().String(client.FlagNode, "127.0.0.1:26657", "Url of node communicate with")
	cmd.Flags().String(client.FlagHome, homeDir+"/.cyberdcli", "Cyberd CLI home folder")

	_ = viper.BindPFlag(client.FlagPassphrase, cmd.Flags().Lookup(client.FlagPassphrase))
	_ = viper.BindPFlag(client.FlagAddress, cmd.Flags().Lookup(client.FlagAddress))
	_ = viper.BindPFlag(client.FlagNode, cmd.Flags().Lookup(client.FlagNode))
	_ = viper.BindPFlag(client.FlagHome, cmd.Flags().Lookup(client.FlagHome))

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
