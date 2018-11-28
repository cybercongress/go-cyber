package client

import (
	"fmt"
	"github.com/TV4/graceful"
	"github.com/cybercongress/cyberd/cosmos/poc/claim/common"
	"github.com/cybercongress/cyberd/cosmos/poc/claim/context"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
)

func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start claim service",
		RunE: func(cmd *cobra.Command, args []string) error {
			port := viper.GetString(common.FlagPort)

			ctx, err := context.NewClaimContext()
			if err != nil {
				return err
			}

			mux := http.NewServeMux()
			mux.HandleFunc("/claim", ClaimHandlerFn(ctx))

			c := cors.New(cors.Options{
				AllowedOrigins: []string{"*"},
			})

			handler := c.Handler(mux)

			fmt.Println("Starting claim service at port " + port)

			graceful.ListenAndServe(&http.Server{
				Addr:    ":" + port,
				Handler: handler,
			})

			return nil
		},
	}

	cmd.Flags().String(common.FlagPort, "26666", "Port to run on")
	cmd.Flags().String(common.FlagName, "", "Name of account to claim from")
	cmd.Flags().String(common.FlagPassphrase, "", "Passphrase of account to claim from")
	cmd.Flags().String(common.FlagChainId, "", "Chain Id")
	cmd.Flags().String(common.FlagAddress, "", "ClaimFrom of account to claim from")
	cmd.Flags().String(common.FlagNode, "127.0.0.1:26657", "Node url to connect")

	viper.BindPFlag(common.FlagPort, cmd.Flags().Lookup(common.FlagPort))
	viper.BindPFlag(common.FlagName, cmd.Flags().Lookup(common.FlagName))
	viper.BindPFlag(common.FlagPassphrase, cmd.Flags().Lookup(common.FlagPassphrase))
	viper.BindPFlag(common.FlagChainId, cmd.Flags().Lookup(common.FlagChainId))
	viper.BindPFlag(common.FlagAddress, cmd.Flags().Lookup(common.FlagAddress))
	viper.BindPFlag(common.FlagNode, cmd.Flags().Lookup(common.FlagNode))

	return cmd

}
