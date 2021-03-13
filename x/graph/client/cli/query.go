package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/spf13/cobra"

	"github.com/cybercongress/go-cyber/x/graph/types"
)

// GetQueryCmd returns
func GetQueryCmd() *cobra.Command {
	graphQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the graph module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	graphQueryCmd.AddCommand(
		GetCmdInLinks(),
		GetCmdOutLinks(),
		GetCmdLinksAmount(),
		GetCmdCidsAmount(),
		GetCmdGraphStats(),
	)

	return graphQueryCmd
}

func GetCmdInLinks() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "in [cid]",
		Short: "Query the current in links by given CID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.InLinks(
				context.Background(),
				&types.QueryLinksRequest{Cid: args[0]},
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

func GetCmdOutLinks() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "out [cid]",
		Short: "Query the current out links by given CID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.OutLinks(
				context.Background(),
				&types.QueryLinksRequest{Cid: args[0]},
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

func GetCmdLinksAmount() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "amount-links",
		Short: "Query the links amount",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.LinksAmount(
				context.Background(),
				&types.QueryLinksAmountRequest{},
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

func GetCmdCidsAmount() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "amount-cids",
		Short: "Query the cids amount",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CidsAmount(
				context.Background(),
				&types.QueryCidsAmountRequest{},
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

func GetCmdGraphStats() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "graph-stats",
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