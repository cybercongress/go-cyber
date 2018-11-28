package main

import (
	"fmt"
	"github.com/cybercongress/cyberd/claim/client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cyberdclaim = &cobra.Command{
		Use:   "cyberdclaim",
		Short: "Service to claim tokens in cyberd zeronet",
	}
)

func main() {

	cyberdclaim.AddCommand(client.StartCmd())

	homeDir, err := getHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.SetDefault("home", homeDir+"/.cyberdcli")

	if err := cyberdclaim.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getHomeDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return home, nil
}
