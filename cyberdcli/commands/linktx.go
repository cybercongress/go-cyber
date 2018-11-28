package commands

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	cbd "github.com/cybercongress/cyberd/app/types"
	. "github.com/cybercongress/cyberd/cyberdcli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagCidFrom = "cid-from"
	flagCidTo   = "cid-to"
)

// LinkTxCmd will create a link tx and sign it with the given key.
func LinkTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link",
		Short: "Create and sign a link tx",
		RunE: func(cmd *cobra.Command, args []string) error {

			txCtx := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			cidFrom := cbd.Cid(viper.GetString(flagCidFrom))
			cidTo := cbd.Cid(viper.GetString(flagCidTo))

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			// ensure that account exists in chain
			_, err = cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := BuildMsg(from, cidFrom, cidTo)

			return utils.CompleteAndBroadcastTxCli(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagCidFrom, "", "Content id to link from")
	cmd.Flags().String(flagCidTo, "", "Content id to link to")

	return cmd
}
