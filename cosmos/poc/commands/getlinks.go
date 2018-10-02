package commands

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cybercongress/cyberd/cosmos/poc/app"
	"github.com/spf13/cobra"
)

type ContentIdLinks struct {
	ContentID  string         `json:"cid"`
	LinkedCIDS map[string]int `json:"linkedCids"`
}

// GetAccountCmd returns a query account that will display the state of the
// account at a given address.
func GetLinksCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "links [cid]",
		Short: "Query cid links",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// find the key to look up the account
			cid := args[0]

			key, err := json.Marshal(cid)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(key, storeName)
			if err != nil || len(res) == 0 {
				return err
			}

			links := new(app.ContentIdLinks)
			err = json.Unmarshal(res, &links)
			if err != nil {
				panic(err)
			}

			output, err := wire.MarshalJSONIndent(cdc, links)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
}
