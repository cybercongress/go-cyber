package cmd

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cybercongress/cyberd/app"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hleb-albau/ethereum-pubkey-collector/crypto"
	"github.com/hleb-albau/ethereum-pubkey-collector/storage"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func LotteryBalancesCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lottery-balances <path-to-pou-file> <path-to-pubkeys-db>",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {

			db, err := storage.OpenDb(args[1])
			if err != nil {
				return err
			}

			accs, balances, err := readPouEthAccounts(args[0])
			if err != nil {
				return err
			}

			resultMap := make(map[string]cyberAccInfo, len(accs))
			for i := range accs {
				ethAddr := common.HexToAddress(accs[i])
				ethRawPubkey := db.GetAddressPublicKey(ethAddr)
				if ethRawPubkey == nil {
					println(accs[i])
					continue
				}
				cosmosAddr := crypto.CosmosAddressFromEthKey(ethRawPubkey)
				accInfo := cyberAccInfo{
					Balance: balances[i],
					Address: crypto.EncodeToHex(cosmosAddr, "cyber"),
				}
				resultMap[accs[i]] = accInfo
			}

			resultJson, err := json.MarshalIndent(resultMap, "", "\t")
			if err != nil {
				return err
			}

			return ioutil.WriteFile("lottery.json", resultJson, 0777)
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")

	return cmd
}

type cyberAccInfo struct {
	Address string `json:"address"`
	Balance int64  `json:"balance"`
}

func readPouEthAccounts(path string) ([]string, []int64, error) {
	accs := make([]string, 0)
	balances := make([]int64, 0)
	pocFile, err := os.Open(path)
	if err != nil {
		return accs, balances, err
	}

	reader := csv.NewReader(bufio.NewReader(pocFile))
	reader.Comma = ' '

	for {

		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return accs, balances, err
		}

		accShare, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return accs, balances, err
		}
		accs = append(accs, line[0])
		balances = append(balances, pouShareToAmt(accShare))
	}

	return accs, balances, nil
}

// output example
/*
{
	"0x9f4062f6153ff4dbf93f6a6f686ed3c906bf0684: : {
		"address" : "cyber1f9yjqmxh6prsmgpcaqj8lmjnxg644n5074zznm",
		"balance" : 1234
	},
	"0x9f4062f6153ff4dbf93f6a6f686ed3c906bf0684: : {
		"address" : "cyber1f9yjqmxh6prsmgpcaqj8lmjnxg644n5074zznm",
		"balance" : 1234
	}
}
*/
