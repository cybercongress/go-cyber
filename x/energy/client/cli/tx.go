package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cybercongress/go-cyber/x/energy/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	energyTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	energyTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreateEnergyRoute(cdc),
		GetCmdEditEnergyRoute(cdc),
		GetCmdDeleteEnergyRoute(cdc),
		GetCmdEditEnergyRouteAlias(cdc),
	)...)

	return energyTxCmd
}

func GetCmdCreateEnergyRoute(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-route [destination-addr] [alias]",
		Args:  cobra.ExactArgs(2),
		Short: "Create energy route from your address to destination address with provided alias",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create energy route from from your address to destination address with provided alias.

Example:
$ %s tx energy create-route cyber1p5ygqrqcw5cn59em9dpgmk7jcnu5twc6mkxlsp cyber19wyeu0fh8lc8ajjq84q9fkwkrwrxxf98su97pc alias --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			src := cliCtx.GetFromAddress()
			dst, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateEnergyRoute(src, dst, args[1])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdEditEnergyRoute(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "edit-route [destination-addr] [value]",
		Args:  cobra.ExactArgs(2),
		Short: "Route energy value to destination address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Route energy from an amount of liquid coins to a destination account from your wallet.

Example:
$ %s tx energy edit-route cyber19wyeu0fh8lc8ajjq84q9fkwkrwrxxf98su97pc 1000stake --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			amount, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			src := cliCtx.GetFromAddress()
			dst, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgEditEnergyRoute(src, dst, amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdDeleteEnergyRoute(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-route [destination-addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Delete yours energy route to given destination address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Delete yours energy route to destination address.

Example:
$ %s tx energy delete-route cyber1p5ygqrqcw5cn59em9dpgmk7jcnu5twc6mkxlsp cyber19wyeu0fh8lc8ajjq84q9fkwkrwrxxf98su97pc --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			src := cliCtx.GetFromAddress()
			dst, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteEnergyRoute(src, dst)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdEditEnergyRouteAlias(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "edit-route-alias [destination-addr] [alias]",
		Args:  cobra.ExactArgs(2),
		Short: "Delete yours energy route to given destination address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Delete yours energy route to destination address.

Example:
$ %s tx energy edit-route-alias cyber19wyeu0fh8lc8ajjq84q9fkwkrwrxxf98su97pc alias --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			src := cliCtx.GetFromAddress()
			dst, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgEditEnergyRouteAlias(src, dst, args[1])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
