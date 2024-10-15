package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"

	"github.com/cybercongress/go-cyber/v4/x/graph/types"
)

func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Graph transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdLink(),
	)

	return txCmd
}

func GetCmdLink() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cyberlink [cid-from] [cid-to]",
		Short: "Create cyberlink",
		Args:  cobra.ExactArgs(2),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create cyberlink.
Example:
$ %s tx link cyberlink QmWZYRj344JSLShtBnrMS4vw5DQ2zsGqrytYKMqcQgEneB QmfZwbahFLTcB3MTMT8TA8si5khhRmzm7zbHToo4WVK3zn
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cidFrom := types.Cid(args[0])
			cidTo := types.Cid(args[1])

			if _, err := cid.Decode(string(cidFrom)); err != nil {
				return types.ErrInvalidParticle
			}

			if _, err := cid.Decode(string(cidTo)); err != nil {
				return types.ErrInvalidParticle
			}

			msg := types.NewMsgCyberlink(
				clientCtx.GetFromAddress(),
				[]types.Link{
					{From: string(cidFrom), To: string(cidTo)},
				},
			)

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
