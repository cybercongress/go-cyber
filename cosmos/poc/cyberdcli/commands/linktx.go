package commands

import (
	"github.com/cybercongress/cyberd/cosmos/poc/app/storage"
	. "github.com/cybercongress/cyberd/cosmos/poc/cyberdcli/util"
	"os"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagCidFrom = "cid-from"
	flagCidTo   = "cid-to"
)

// LinkTxCmd will create a link tx and sign it with the given key.
func LinkTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link",
		Short: "Create and sign a link tx",
		RunE: func(cmd *cobra.Command, args []string) error {

			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			cidFrom := storage.Cid(viper.GetString(flagCidFrom))
			cidTo := storage.Cid(viper.GetString(flagCidTo))

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

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagCidFrom, "", "Content id to link from")
	cmd.Flags().String(flagCidTo, "", "Content id to link to")

	return cmd
}
