package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/cybercongress/go-cyber/x/energy/types"
)

// GetTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	energyTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	energyTxCmd.AddCommand(
		GetCmdCreateEnergyRoute(),
		GetCmdEditEnergyRoute(),
		GetCmdDeleteEnergyRoute(),
		GetCmdEditEnergyRouteAlias(),
	)

	return energyTxCmd
}

func GetCmdCreateEnergyRoute() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-route [destination-addr] [alias]",
		Args:  cobra.ExactArgs(2),
		Short: "Create energy route from your address to destination address with provided alias",
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

			msg := types.NewMsgCreateEnergyRoute(src, dst, args[1])

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

func GetCmdEditEnergyRoute() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-route [destination-addr] [value]",
		Args:  cobra.ExactArgs(2),
		Short: "Route energy value to destination address",
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

			msg := types.NewMsgEditEnergyRoute(src, dst, amount)

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

func GetCmdDeleteEnergyRoute() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-route [destination-addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Delete yours energy route to given destination address",
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

			msg := types.NewMsgDeleteEnergyRoute(src, dst)
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

func GetCmdEditEnergyRouteAlias() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-route-alias [destination-addr] [alias]",
		Args:  cobra.ExactArgs(2),
		Short: "Delete yours energy route to given destination address",
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

			msg := types.NewMsgEditEnergyRouteAlias(src, dst, args[1])

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
