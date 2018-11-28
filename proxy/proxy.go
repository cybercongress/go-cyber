package main

import (
	"fmt"
	"github.com/cybercongress/cyberd/proxy/client"
	"github.com/spf13/cobra"
	"os"
)

var (
	cyberdproxy = &cobra.Command{
		Use:   "cyberdproxy",
		Short: "Http proxy to cyberd node",
	}
)

func main() {

	cyberdproxy.AddCommand(client.StartCmd())

	if err := cyberdproxy.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
