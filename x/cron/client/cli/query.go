package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		Short: "Query the current cron module parameters information",
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
		Short: "Query job",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			creator, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			res, err := queryClient.Job(
				context.Background(),
				&types.QueryJobParamsRequest{
					contract.String(), creator.String(), args[2],
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
		Short: "Query job stats",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			creator, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			res, err := queryClient.JobStats(
				context.Background(),
				&types.QueryJobParamsRequest{
					contract.String(), creator.String(), args[2],
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

