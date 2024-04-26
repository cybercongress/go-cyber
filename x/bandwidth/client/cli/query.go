package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/cybercongress/go-cyber/v4/x/bandwidth/types"
)

func GetQueryCmd() *cobra.Command {
	bandwidthQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the bandwidth module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	bandwidthQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryPrice(),
		GetCmdQueryLoad(),
		GetCmdQueryTotalBandwidth(),
		GetCmdQueryNeuron(),
	)

	return bandwidthQueryCmd
}

func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current bandwidth module parameters information",
		Args:  cobra.NoArgs,
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

func GetCmdQueryLoad() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load",
		Short: "Query the bandwidth load",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Load(
				context.Background(),
				&types.QueryLoadRequest{},
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

func GetCmdQueryPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "price",
		Short: "Query the bandwidth price",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Price(
				context.Background(),
				&types.QueryPriceRequest{},
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

func GetCmdQueryTotalBandwidth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total",
		Short: "Query the total bandwidth",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.TotalBandwidth(
				context.Background(),
				&types.QueryTotalBandwidthRequest{},
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

func GetCmdQueryNeuron() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "neuron [address]",
		Short: "Query the neuron bandwidth [address]",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.NeuronBandwidth(
				context.Background(),
				&types.QueryNeuronBandwidthRequest{Neuron: addr.String()},
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
