package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/cybercongress/go-cyber/x/bandwidth/types"
)

// GetQueryCmd returns the cli query commands for the bandwidth module.
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	bandwidthQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the bandwidth module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	bandwidthQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQueryParams(cdc),
			//GetCmdQueryPrice(cdc), TODO Amino:JSON float* support requires `amino:"unsafe"`.
			//GetCmdQueryLoad(cdc), TODO Amino:JSON float* support requires `amino:"unsafe"`.
			GetCmdQueryAccount(cdc),
		)...,
	)

	return bandwidthQueryCmd
}

func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current bandwidth parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := cdc.UnmarshalJSON(res, &params); err != nil {
				return err
			}

			return cliCtx.PrintOutput(params)
		},
	}
}

func GetCmdQueryLoad(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "load",
		Short: "Query the bandwidth load",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLoad)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var resp types.ResultLoad
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(resp)
		},
	}
}

func GetCmdQueryPrice(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "price",
		Short: "Query the bandwidth price",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPrice)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var resp types.ResultPrice
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}
}

func GetCmdQueryAccount(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "status [account]",
		Short: "Query the account bandwidth [account-addr]",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.NewQueryAccountParams(addr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAccount)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var resp types.AcсountBandwidth
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(resp)
		},
	}
}