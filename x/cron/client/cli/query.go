package cli

import (
	"fmt"
	"strings"
	"context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cybercongress/go-cyber/x/cron/types"
)

func GetQueryCmd() *cobra.Command {
	cronQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cronQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryJob(),
		GetCmdQueryJobStats(),
		GetCmdQueryJobs(),
		GetCmdQueryJobsStats(),
	)

	return cronQueryCmd
}

func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current energy module parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as staking parameters.

Example:
$ %s query energy params
`,
				version.Version,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{},
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

func GetCmdQueryJob() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "job [contract] [creator] [label]",
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s query cron job cosmosvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.Version,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Job(
				context.Background(),
				&types.QueryJobParamsRequest{
					args[1], args[0], args[2],
				},
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

func GetCmdQueryJobStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "job-stats [contract] [creator] [label]",
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s query cron job-stats cosmosvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.Version,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.JobStats(
				context.Background(),
				&types.QueryJobParamsRequest{
					args[1], args[0], args[2],
				},
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

func GetCmdQueryJobs() *cobra.Command {
	cmd :=  &cobra.Command{
		Use:   "jobs",
		Short: "Query all jobs",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s query cron jobs
`,
				version.Version,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Jobs(
				context.Background(),
				&types.QueryJobsRequest{},
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

func GetCmdQueryJobsStats() *cobra.Command {
	cmd :=  &cobra.Command{
		Use:   "jobs-stats",
		Short: "Query all jobs stats",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s query cron jobs-stats
`,
				version.Version,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.JobsStats(
				context.Background(),
				&types.QueryJobsStatsRequest{},
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

