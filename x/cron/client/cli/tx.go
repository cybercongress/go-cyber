package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/cybercongress/go-cyber/x/cron/types"
)

// GetTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	cronTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cronTxCmd.AddCommand(
		GetCmdAddJob(),
		GetCmdRemoveJob(),
		GetCmdChangeJobCID(),
		GetCmdChangeJobLabel(),
		GetCmdChangeJobCallData(),
		GetCmdChangeJobGasPrice(),
		GetCmdChangeJobPeriod(),
		GetCmdChangeJobBlock(),
	)

	return cronTxCmd
}

func GetCmdAddJob() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-job [contract] [trigger-period] [trigger-block] [load-calldata] [load-gasprice] [label] [cid]",
		Args:  cobra.ExactArgs(7),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron add-job cyber___ 10 0 "{"release":{}}" 1 "label" Qm___ --from mykey
`,
				version.Version,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			period, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("block period %s not a valid uint, please input a valid block period", args[1])
			}

			block, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("block %s not a valid uint, please input a valid block ", args[2])
			}

			gasprice, err := sdk.ParseCoinNormalized(args[4])
			if err != nil {
				return err
			}

			msg := types.NewMsgCronAddJob(
				creator, contract,
				types.NewTrigger(period, block),
				types.NewLoad(args[3], gasprice),
				args[5], args[6],
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)


	return cmd
}

func GetCmdRemoveJob() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-job [contract] [label]",
		Args:  cobra.ExactArgs(2),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron remove-job cyber___ "label" --from mykey
`,
				version.Version,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgCronRemoveJob(
				creator, contract, args[1],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)


	return cmd
}

func GetCmdChangeJobCID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-job-cid [contract] [label] [cid]",
		Args:  cobra.ExactArgs(3),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron change-job-cid cyber___ "label" Qm___ --from mykey
`,
				version.Version,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgCronChangeJobCID(
				creator, contract, args[1], args[2],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)


	return cmd
}

func GetCmdChangeJobLabel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-job-label [contract] [label] [newLabel]",
		Args:  cobra.ExactArgs(3),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron change-job-label cyber___ "label" "newLabel" --from mykey
`,
				version.Version,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgCronChangeJobLabel(
				creator, contract, args[1], args[2],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)


	return cmd
}

func GetCmdChangeJobCallData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-job-calldata [contract] [label] [calldata]",
		Args:  cobra.ExactArgs(3),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron change-job-calldata cyber___ "label" "{"release":{}}" --from mykey
`,
				version.Version,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgCronChangeCallData(
				creator, contract, args[1], args[2],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)


	return cmd
}

func GetCmdChangeJobGasPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-job-gasprice [contract] [label] [gasprice]",
		Args:  cobra.ExactArgs(3),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron change-job-gasprice cyber___ "label" 10 --from mykey
`,
				version.Version,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			gasprice, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgCronChangeGasPrice(
				creator, contract, args[1], gasprice,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)


	return cmd
}

func GetCmdChangeJobPeriod() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-job-period [contract] [label] [period]",
		Args:  cobra.ExactArgs(3),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron change-job-period cyber___ "label" 10 --from mykey
`,
				version.Version,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			period, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("block period %s not a valid uint, please input a valid block period", args[2])
			}

			msg := types.NewMsgCronChangeJobPeriod(
				creator, contract, args[1], period,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)


	return cmd
}

func GetCmdChangeJobBlock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-job-block [contract] [label] [block]",
		Args:  cobra.ExactArgs(3),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron change-job-block cyber___ "label" 10 --from mykey
`,
				version.Version,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			block, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("block %s not a valid uint, please input a valid block ", args[2])
			}

			msg := types.NewMsgCronChangeJobBlock(
				creator, contract, args[1], block,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)


	return cmd
}