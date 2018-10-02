package commands

import (
	client2 "github.com/cybercongress/cyberd/cosmos/poc/client"
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
	flagCid1 = "cid-from"
	flagCid2 = "cid-to"
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

			cid1 := viper.GetString(flagCid1)
			cid2 := viper.GetString(flagCid2)

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			_, err = cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := client2.BuildMsg(from, cid1, cid2)

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagCid1, "", "Content id to link from")
	cmd.Flags().String(flagCid2, "", "Content id to link to")

	return cmd
}
