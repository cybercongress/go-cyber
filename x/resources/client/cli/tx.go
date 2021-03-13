package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	ctypes "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/resources/types"
)


func NewTxCmd() *cobra.Command {
	resourcesTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	resourcesTxCmd.AddCommand(
		GetCmdConvert(),
	)

	return resourcesTxCmd
}

func GetCmdConvert() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert [amount] [resource] [end-time]",
		Args:  cobra.ExactArgs(3),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx resources convert 10000cyb 100000 --from mykey
`,
				version.Version,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			agent := clientCtx.GetFromAddress()

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return fmt.Errorf("coin %s not a valid coin, please input a valid coin", args[0])
			}

			if amount.Denom != ctypes.CYB {
				return fmt.Errorf("coin %s not a valid coin, please input a valid coin", args[0])
			}

			if args[1] != ctypes.VOLT && args[1] != ctypes.AMPER {
				return fmt.Errorf("resource %s not a valid resource, please input a valid resource", args[1])
			}

			endBlock, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("block period %s not a valid uint, please input a valid block period", args[1])
			}

			msg := types.NewMsgConvert(agent, amount, args[1], int64(endBlock))

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)


	return cmd
}