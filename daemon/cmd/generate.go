package cmd

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
)

func GenerateAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-account",
		Short: "Generate account",
		RunE: func(_ *cobra.Command, args []string) error {

			addr, seed, err := server.GenerateCoinKey()
			if err != nil {
				return err
			}

			fmt.Println("Address: ", addr.String())
			fmt.Println("Seed phrase: ", seed)

			return nil
		},
	}

	return cmd
}
