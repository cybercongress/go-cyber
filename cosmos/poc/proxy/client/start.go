package client

import (
	"fmt"
	"github.com/TV4/graceful"
	"github.com/cybercongress/cyberd/cosmos/poc/proxy/core"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
)

const (
	flagNode = "node"
	flagPort = "port"
)

func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			node := viper.GetString(flagNode)
			port := viper.GetString(flagPort)

			ctx := core.NewProxyContext(node)

			mux := http.NewServeMux()
			mux.HandleFunc("/link", core.TxHandlerFn(ctx, core.UnmarshalLinkRequest))
			mux.HandleFunc("/send", core.TxHandlerFn(ctx, core.UnmarshalSendRequest))
			mux.HandleFunc("/search", core.GetWithParamsHandlerFn(ctx, "/search", []string{"cid", "page", "perPage"}))
			mux.HandleFunc("/account", core.GetWithParamsHandlerFn(ctx, "/account", []string{"address"}))
			mux.HandleFunc("/health", core.GetHandlerFn(ctx, "/health"))
			mux.HandleFunc("/status", core.GetHandlerFn(ctx, "/status"))

			c := cors.New(cors.Options{
				AllowedOrigins: []string{"*"},
			})

			handler := c.Handler(mux)

			fmt.Println("Connecting to node " + node)
			fmt.Println("Running on port " + port)

			graceful.ListenAndServe(&http.Server{
				Addr:    ":" + port,
				Handler: handler,
			})

			return nil
		},
	}

	cmd.Flags().String(flagNode, "http://localhost:26657", "Node url")
	cmd.Flags().String(flagPort, "26660", "Port to run on")

	viper.BindPFlag(flagNode, cmd.Flags().Lookup(flagNode))
	viper.BindPFlag(flagPort, cmd.Flags().Lookup(flagPort))

	return cmd

}
