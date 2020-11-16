package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cybercongress/go-cyber/x/energy/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group power queries under a subcommand
	energyQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	energyQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryParams(queryRoute, cdc),
		GetCmdQuerySourceRoutes(queryRoute, cdc),
		GetCmdQueryDestinationRoutes(queryRoute, cdc),
		GetCmdQuerySourceRoutedEnergy(queryRoute, cdc),
		GetCmdQueryDestinationRoutedEnergy(queryRoute, cdc),
		GetCmdQueryRoute(queryRoute, cdc),
		GetCmdQueryRoutes(queryRoute, cdc))...)

	return energyQueryCmd
}

func GetCmdQueryParams(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current energy module parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as staking parameters.

Example:
$ %s query energy params
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryParams)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(bz, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}

func GetCmdQuerySourceRoutes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "routes-from [source-addr]",
		Short: "Query all energy routes made from source account",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query routes from an individual source account.

Example:
$ %s query energy routes-from cyber1p5ygqrqcw5cn59em9dpgmk7jcnu5twc6mkxlsp
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			src, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.NewQuerySourceParams(src))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QuerySourceRoutes)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var resp types.Routes
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(resp)
		},
	}
}

func GetCmdQueryDestinationRoutes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "routes-to [destination-addr]",
		Short: "Query all energy routes made to destination account",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query routes to an individual destination account.

Example:
$ %s query energy routes-to cyber19wyeu0fh8lc8ajjq84q9fkwkrwrxxf98su97pc
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			dst, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.NewQueryDestinationParams(dst))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDestinationRoutes)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var resp types.Routes
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(resp)
		},
	}
}

func GetCmdQuerySourceRoutedEnergy(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "routed-from [src-addr]",
		Short: "Query energy routed from source account",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query energy routed from an individual source account.

Example:
$ %s query energy routed-from cyber1p5ygqrqcw5cn59em9dpgmk7jcnu5twc6mkxlsp
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			src, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.NewQuerySourceParams(src))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QuerySourceRoutedEnergy)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var resp sdk.Int
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(resp)
		},
	}
}

func GetCmdQueryDestinationRoutedEnergy(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "routed-to [destination-addr]",
		Short: "Query energy routed to destination account",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query energy routed on an individual destination account.

Example:
$ %s query energy routed-to cyber19wyeu0fh8lc8ajjq84q9fkwkrwrxxf98su97pc
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			dst, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.NewQueryDestinationParams(dst))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDestinationRoutedEnergy)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var resp sdk.Int
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(resp)
		},
	}
}

func GetCmdQueryRoute(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "route  [source-addr] [destination-addr]",
		Short: "Query energy route for given source and destination accounts",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query energy route for given source and destination accounts.

Example:
$ %s query energy route cyber1p5ygqrqcw5cn59em9dpgmk7jcnu5twc6mkxlsp cyber19wyeu0fh8lc8ajjq84q9fkwkrwrxxf98su97pc
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			src, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			dst, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.NewQueryRouteParams(src, dst))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRoute)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var resp types.Route
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(resp)
		},
	}
}

func GetCmdQueryRoutes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "routes",
		Short: "Query all energy routes",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all energy routes.

Example:
$ %s query energy routes
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRoutes)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var resp types.Routes
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(resp)
		},
	}
}

