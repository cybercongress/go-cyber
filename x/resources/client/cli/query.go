package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	ctypes "github.com/cybercongress/go-cyber/v2/types"
	"github.com/cybercongress/go-cyber/v2/x/resources/types"
)

func GetQueryCmd() *cobra.Command {
	energyQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	energyQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryInvestmintAmount(),
	)

	return energyQueryCmd
}

func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current resources module parameters information",
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

func GetCmdQueryInvestmintAmount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "investmint [amount] [resource] [length]",
		Short: "Query resources return on investmint",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return fmt.Errorf("coin %s not a valid coin, please input a valid coin", args[0])
			}

			if amount.Denom != ctypes.SCYB {
				return fmt.Errorf("coin %s not a valid coin, please input a valid coin", args[0])
			}

			if args[1] != ctypes.VOLT && args[1] != ctypes.AMPERE {
				return fmt.Errorf("resource %s not a valid resource, please input a valid resource", args[1])
			}

			length, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("block period %s not a valid uint, please input a valid block period", args[1])
			}

			res, err := queryClient.Investmint(
				context.Background(),
				&types.QueryInvestmintRequest{
					Amount:   amount,
					Resource: args[1],
					Length:   length,
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
