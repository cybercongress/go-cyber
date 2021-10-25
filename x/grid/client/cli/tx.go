package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/cybercongress/go-cyber/x/grid/types"
)

// GetTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	gridTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	gridTxCmd.AddCommand(
		GetCmdCreateRoute(),
		GetCmdEditRoute(),
		GetCmdDeleteRoute(),
		GetCmdEditRouteName(),
	)

	return gridTxCmd
}

func GetCmdCreateRoute() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-route [destination] [name]",
		Args:  cobra.ExactArgs(2),
		Short: "Create grid route from your address to destination address with provided name",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			src := clientCtx.GetFromAddress()
			dst, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateRoute(src, dst, args[1])

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdEditRoute() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-route [destination] [value]",
		Args:  cobra.ExactArgs(2),
		Short: "Set value of grid route to destination address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			src := clientCtx.GetFromAddress()
			dst, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgEditRoute(src, dst, amount)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdDeleteRoute() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-route [destination]",
		Args:  cobra.ExactArgs(1),
		Short: "Delete your grid route to given destination address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			src := clientCtx.GetFromAddress()
			dst, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteRoute(src, dst)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdEditRouteName() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-route-name [destination] [name]",
		Args:  cobra.ExactArgs(2),
		Short: "Edit name of grid route to given destination address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			src := clientCtx.GetFromAddress()
			dst, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgEditRouteName(src, dst, args[1])

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
