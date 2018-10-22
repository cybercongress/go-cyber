package main

import (
	"fmt"
	"github.com/cybercongress/cyberd/cosmos/poc/claim/client"
	"github.com/spf13/cobra"
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

	if err := cyberdclaim.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

