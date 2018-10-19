package main

import (
	"fmt"
	"github.com/cybercongress/cyberd/cosmos/poc/cyberdproxy/proxy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"

	"github.com/TV4/graceful"
	"github.com/rs/cors"
)

var (
	cyberdproxy = &cobra.Command{
		Use:   "cyberdproxy",
		Short: "Http proxy to cyberd node",
	}
)

func main() {

	cyberdproxy.AddCommand(StartCmd())

	if err := cyberdproxy.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

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

			ctx := proxy.NewProxyContext(node)

			mux := http.NewServeMux()
			mux.HandleFunc("/link", proxy.LinkHandlerFn(ctx))
			mux.HandleFunc("/search", proxy.GetWithParamHandlerFn(ctx, "/search", "cid"))
			mux.HandleFunc("/account", proxy.GetWithParamHandlerFn(ctx, "/account", "address"))
			mux.HandleFunc("/health", proxy.GetHandlerFn(ctx, "/health"))
			mux.HandleFunc("/status", proxy.GetHandlerFn(ctx, "/status"))

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
