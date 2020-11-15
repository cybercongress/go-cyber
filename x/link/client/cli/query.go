package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	//"github.com/litvintech/cyber/x/link/internal/types"

	"github.com/cybercongress/go-cyber/x/link/types"
)

// GetQueryCmd returns
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	linkQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the link module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	linkQueryCmd.AddCommand(flags.GetCommands(
		GetCmdInLinks(cdc),
		GetCmdOutLinks(cdc),
		GetCmdLinksAmount(cdc),
	)...)

	return linkQueryCmd
}

func GetCmdInLinks(cdc *codec.Codec) *cobra.Command{
	return &cobra.Command{
		Use:   "in",
		Short: "Query the current in links by given CID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryInLinks)

			bz, err := cdc.MarshalJSON(types.NewQueryLinksPrams(types.Cid(args[0])))
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var resp types.ResultLinks
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(resp)
		},
	}
}

func GetCmdOutLinks(cdc *codec.Codec) *cobra.Command{
	return &cobra.Command{
		Use:   "out",
		Short: "Query the current out links by given CID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryOutLinks)

			bz, err := cdc.MarshalJSON(types.NewQueryLinksPrams(types.Cid(args[0])))
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var resp types.ResultLinks
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(resp)
		},
	}
}

func GetCmdLinksAmount(cdc *codec.Codec) *cobra.Command{
	return &cobra.Command{
		Use:   "amount",
		Short: "Query the current out links",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLinksAmount)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var resp uint64
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(resp)
		},
	}
}