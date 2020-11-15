
package cli

import (
	"bufio"
	//"bytes"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"

	"github.com/cybercongress/go-cyber/x/link/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	linkCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "link transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	linkCmd.AddCommand(flags.PostCommands(
		GetCmdLink(cdc),
	)...)

	return linkCmd
}

func GetCmdLink(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [cid-from] [cid-to]",
		Short: "Create cyberlink",
		Args:  cobra.ExactArgs(2),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create cyberlink.
Example:
$ %s tx link cyberlink QmWZYRj344JSLShtBnrMS4vw5DQ2zsGqrytYKMqcQgEneB QmfZwbahFLTcB3MTMT8TA8si5khhRmzm7zbHToo4WVK3zn
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))

			cidFrom := types.Cid(args[0])
			cidTo   := types.Cid(args[1])

			if _, err := cid.Decode(string(cidFrom)); err != nil {
				return types.ErrInvalidCid
			}

			if _, err := cid.Decode(string(cidTo)); err != nil {
				return types.ErrInvalidCid
			}

			msg := types.NewMsgCyberlink(
				cliCtx.GetFromAddress(),
				[]types.Link{
					{From: cidFrom, To: cidTo},
				},
			)

			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}