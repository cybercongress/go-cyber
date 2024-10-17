package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/spf13/cobra"

	"github.com/cybercongress/go-cyber/v5/x/graph/types"
)

// GetQueryCmd returns
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the graph module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(
		GetCmdGraphStats(),
	)

	return queryCmd
}

func GetCmdGraphStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Query the graph stats",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.GraphStats(
				context.Background(),
				&types.QueryGraphStatsRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
