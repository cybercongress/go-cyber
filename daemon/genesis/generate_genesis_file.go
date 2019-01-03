package genesis

import (
	"bufio"
	"encoding/csv"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/app/genesis"
	"github.com/cybercongress/cyberd/x/mint"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/types"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	bitcoinHeightZeroTime = time.Unix(1231006505, 0).UTC() // 2009-01-03 18:15:05 +0000 UTC
	chainId               = "euler"
	genesisSupply         = mint.GenesisSupply
	pocPercentage         = 0.7
)

func GenerateEulerGenesisFile(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-euler-genesis-block",
		Short: "Generate genesis file for euler testnet.",
		RunE: func(_ *cobra.Command, args []string) error {

			// proof of use accs
			pouAccs, err := readPouAccounts()
			if err != nil {
				return err
			}

			appState := genesis.NewDefaultGenesisState()
			appState.Accounts = append(appState.Accounts, getGenesisAccs()...)
			appState.Accounts = append(appState.Accounts, pouAccs...)

			// deduplicate
			addrMap := make(map[string]int64, len(appState.Accounts))
			for _, acc := range appState.Accounts {
				strAddr := acc.Address.String()
				if _, ok := addrMap[strAddr]; ok {
					log.Printf("duplicate account in genesis state: Address %v", acc.Address.String())
				}
				addrMap[strAddr] += acc.Amount
			}

			addrAsArray := make([]genesis.GenesisAccount, 0)
			for k, v := range addrMap {
				addrAsArray = append(addrAsArray, genesis.GenesisAccount{Address: addr(k), Amount: v})
			}
			appState.Accounts = addrAsArray

			appState.StakeData.Pool.LooseTokens = sdk.NewDec(genesisSupply)
			stateAsJson, err := codec.MarshalJSONIndent(cdc, appState)
			if err != nil {
				return err
			}

			genDoc := types.GenesisDoc{
				ChainID:     chainId,
				Validators:  make([]types.GenesisValidator, 0),
				AppState:    stateAsJson,
				GenesisTime: bitcoinHeightZeroTime.AddDate(10, 0, 0),
			}

			err = genDoc.SaveAs(ctx.Config.GenesisFile())
			return err
		},
	}
	return cmd
}

func readPouAccounts() ([]genesis.GenesisAccount, error) {

	accs := make([]genesis.GenesisAccount, 0)
	pocFile, err := os.Open("/home/hlb/.cyberd/proof-of-code")
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(bufio.NewReader(pocFile))
	reader.Comma = ' '

	for {

		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		accAddress, err := sdk.AccAddressFromBech32(line[0])
		if err != nil {
			return nil, err
		}

		accAmtPercent, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return nil, err
		}

		accs = append(accs, genesis.GenesisAccount{
			Address: accAddress,
			Amount:  amt(accAmtPercent),
		})
	}

	return accs, nil
}

// Returns all, except genesis poc accs
func getGenesisAccs() []genesis.GenesisAccount {
	accs := []genesis.GenesisAccount{
		{Address: addr("cbd1f9yjqmxh6prsmgpcaqj8lmjnxg644n50qjl4vw"), Amount: amt(8.288000001)},
		{Address: addr("cbd1hlu0kqwvxmhjjsezr00jdrvs2k537mqhrv02ja"), Amount: amt(3.045611111)},
		{Address: addr("cbd1myeyqp96pz3tayjdctflrxpwf45dq3xyj56yk0"), Amount: amt(2.1153)},
		{Address: addr("cbd1gannk6qt3s5mnm5smx6xjqqvecu08666hpazlz"), Amount: amt(1.5328)},
		{Address: addr("cbd1sjedcfmqupxcnxudq9w0rxrf87r3c6tvep5fnj"), Amount: amt(1.428)},
		{Address: addr("cbd1ch4dpd8jxkl7w4wnzdx02utmw4j0xatfks6ulv"), Amount: amt(1)},
		{Address: addr("cbd1s3748ghvcwvrws3kxsdc8xnan3qhv77740gnnl"), Amount: amt(0.568211111)},
		{Address: addr("cbd14d92r4svhl4qa3g6q48tjekarw2kt67njlaeht"), Amount: amt(0.083811111)},
		{Address: addr("cbd1up7dk03v4d898vqgmc2y32y7duuylgx8ra7jjj"), Amount: amt(0.043511111)},
		{Address: addr("cbd1rqudjcrdwqedffxufmqgsleuguhm7pka6snns3"), Amount: amt(0.028)},
		{Address: addr("cbd1hmkqhy8ygl6tnl5g8tc503rwrmmrkjcqtqrsx6"), Amount: amt(0.023311111)},
		{Address: addr("cbd1gs92s58t6rkallnml8ufdzrz3038dcylal0nlc"), Amount: amt(0.023311111)},
		{Address: addr("cbd1h7u5zvduvc3dqrfq9hejm35ktfxh3ha7fra64a"), Amount: amt(0.013911111)},
		{Address: addr("cbd1rl3xnsrkpjfwejqfy7v4kntu64hzxy8dgafh6j"), Amount: amt(0.003111111)},
		{Address: addr("cbd1xege0g92p6exmzjv58u7vh3s5zkz75v48mlnev"), Amount: amt(0.003111111)},
	}
	return accs
}

func addr(hex string) sdk.AccAddress {
	accAddress, _ := sdk.AccAddressFromBech32(hex)
	return accAddress
}

func amt(pct float64) int64 {
	return int64(pct / 100 * pocPercentage * float64(genesisSupply))
}
