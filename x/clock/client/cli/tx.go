package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/cybercongress/go-cyber/v6/x/clock/types"
)

// NewTxCmd returns a root CLI command handler for certain modules/Clock
// transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Clock subcommands.",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewRegisterClockContract(),
		NewUnregisterClockContract(),
		NewUnjailClockContract(),
	)
	return txCmd
}

// NewRegisterClockContract returns a CLI command handler for registering a
// contract for the clock module.
func NewRegisterClockContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [contract_bech32]",
		Short: "Register a clock contract.",
		Long:  "Register a clock contract. Sender must be admin of the contract.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			senderAddress := cliCtx.GetFromAddress()
			contractAddress := args[0]

			msg := &types.MsgRegisterClockContract{
				SenderAddress:   senderAddress.String(),
				ContractAddress: contractAddress,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewUnregisterClockContract returns a CLI command handler for unregistering a
// contract for the clock module.
func NewUnregisterClockContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unregister [contract_bech32]",
		Short: "Unregister a clock contract.",
		Long:  "Unregister a clock contract. Sender must be admin of the contract.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			senderAddress := cliCtx.GetFromAddress()
			contractAddress := args[0]

			msg := &types.MsgUnregisterClockContract{
				SenderAddress:   senderAddress.String(),
				ContractAddress: contractAddress,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewUnjailClockContract returns a CLI command handler for unjailing a
// contract for the clock module.
func NewUnjailClockContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unjail [contract_bech32]",
		Short: "Unjail a clock contract.",
		Long:  "Unjail a clock contract. Sender must be admin of the contract.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			senderAddress := cliCtx.GetFromAddress()
			contractAddress := args[0]

			msg := &types.MsgUnjailClockContract{
				SenderAddress:   senderAddress.String(),
				ContractAddress: contractAddress,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
