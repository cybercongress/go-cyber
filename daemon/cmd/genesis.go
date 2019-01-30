package cmd

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/app"
	cbd "github.com/cybercongress/cyberd/types"
	"github.com/spf13/cobra"
)

func GenesisCmds(ctx *server.Context, cdc *codec.Codec) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "genesis",
		Short: "genesis commands",
	}
	rootCmd.AddCommand(AddEulerTokensCmd(ctx, cdc))
	rootCmd.AddCommand(AddMissingEulerTokensCmd(ctx, cdc))
	rootCmd.AddCommand(ChangeGenesisAccsPrefixCmd(ctx, cdc))
	return rootCmd
}

// 0.1.0 Distribution
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Proof Of Code & Proof Of Value Euler addresses.
// Addresses taken from 0x136c1121f21c29415d8cd71f8bb140c7ff187033 Ethereum contract at a 'jan 3, 2019'
// 1 'CBD' token equaling to 1% of genesis supply.
// First address, with 70% allocation is 'Proof Of Use' one.
// The rest is pocv, or 'Proof Of Code' & 'Proof Of Value'

const pouFileCoefficient float64 = 0.7 // means 70%

const tokensToBurnPercentage float64 = 11.8000 // 0x3b6ce0d5fd5f6de16fb6b687207efe86702a6e76
const proofOfUsePercentage float64 = 70.0000   // 0xccf0e5a05bf5c0fb7c2d91737b176b4a2d2fd7f0

var pocvEulerAccs = map[string]float64{
	"cbd1f9yjqmxh6prsmgpcaqj8lmjnxg644n50qjl4vw": 8.288000001,              // 0x9f4062f6153ff4dbf93f6a6f686ed3c906bf0684
	"cbd1hlu0kqwvxmhjjsezr00jdrvs2k537mqhrv02ja": 3.045611111,              // 0x7c4401ae98f12ef6de39ae24cf9fc51f80eba16b
	"cbd1myeyqp96pz3tayjdctflrxpwf45dq3xyj56yk0": 2.1153,                   // 0xf2cb7985a5c3fdd8d7742a73d5dc001bbd32caf8
	"cbd1gannk6qt3s5mnm5smx6xjqqvecu08666hpazlz": 1.5328,                   // 0x002f9caf40a444f20813da783d152bdfaf42852f
	"cbd1sjedcfmqupxcnxudq9w0rxrf87r3c6tvep5fnj": 1.428,                    // 0x00b8fe1a1a2b899418702e32a96e276ff56a4d05
	"cbd1ch4dpd8jxkl7w4wnzdx02utmw4j0xatfks6ulv": 1,                        // 0x8b788b444ca3203bba0fdcae1c482110494e81f1
	"cbd1s3748ghvcwvrws3kxsdc8xnan3qhv77740gnnl": 0.568211111,              // 0x00cff8cf7bff03a9a2a81c01920ffd8cfa7ae9d0
	"cbd14d92r4svhl4qa3g6q48tjekarw2kt67njlaeht": 0.083811111,              // 0x63e65bc441334b27d2178f81f2d701e4e58c158a
	"cbd1up7dk03v4d898vqgmc2y32y7duuylgx8ra7jjj": 0.043511111,              // 0x9d7d6e753f055e40d3767337300e722e934086c1
	"cbd1rqudjcrdwqedffxufmqgsleuguhm7pka6snns3": 0.028,                    // 0x00725d89a2a2fb3b21fd1035b579cbcde4a0991b
	"cbd1hmkqhy8ygl6tnl5g8tc503rwrmmrkjcqtqrsx6": 0.023311111,              // 0x00ca47db1be92c1072e973fd8dc4a082f7d70214
	"cbd1gs92s58t6rkallnml8ufdzrz3038dcylal0nlc": 0.023311111,              // 0x4585c7eaa2cb96d4b59e868929efabeeb8e65b07
	"cbd1h7u5zvduvc3dqrfq9hejm35ktfxh3ha7fra64a": 0.0132800000000000911111, // 0x002b9c5b537a1b6004ed720f32cc808fd6210f26
	"cbd1rl3xnsrkpjfwejqfy7v4kntu64hzxy8dgafh6j": 0.003111111,              // 0x00d3c9033570b8adea9c18780325a45635c55805
	"cbd1xege0g92p6exmzjv58u7vh3s5zkz75v48mlnev": 0.003111111,              // 0x5d01f31f6eda95489ca1e3c6357a9627fa2983de
}

var pocvAccs = map[string]float64{
	"cyber1f9yjqmxh6prsmgpcaqj8lmjnxg644n5074zznm": 8.288000001,              // 0x9f4062f6153ff4dbf93f6a6f686ed3c906bf0684
	"cyber1hlu0kqwvxmhjjsezr00jdrvs2k537mqhatjadg": 3.045611111,              // 0x7c4401ae98f12ef6de39ae24cf9fc51f80eba16b
	"cyber1myeyqp96pz3tayjdctflrxpwf45dq3xyvn8nf6": 2.1153,                   // 0xf2cb7985a5c3fdd8d7742a73d5dc001bbd32caf8
	"cyber1gannk6qt3s5mnm5smx6xjqqvecu08666fxq4qh": 1.5328,                   // 0x002f9caf40a444f20813da783d152bdfaf42852f
	"cyber1sjedcfmqupxcnxudq9w0rxrf87r3c6tv8xf7v8": 1.428,                    // 0x00b8fe1a1a2b899418702e32a96e276ff56a4d05
	"cyber1ch4dpd8jxkl7w4wnzdx02utmw4j0xatfgh8tqe": 1,                        // 0x8b788b444ca3203bba0fdcae1c482110494e81f1
	"cyber1s3748ghvcwvrws3kxsdc8xnan3qhv777tg4yv2": 0.568211111,              // 0x00cff8cf7bff03a9a2a81c01920ffd8cfa7ae9d0
	"cyber14d92r4svhl4qa3g6q48tjekarw2kt67nvcqwg7": 0.083811111,              // 0x63e65bc441334b27d2178f81f2d701e4e58c158a
	"cyber1up7dk03v4d898vqgmc2y32y7duuylgx8a6r9d8": 0.043511111,              // 0x9d7d6e753f055e40d3767337300e722e934086c1
	"cyber1rqudjcrdwqedffxufmqgsleuguhm7pkayhwy0y": 0.028,                    // 0x00725d89a2a2fb3b21fd1035b579cbcde4a0991b
	"cyber1hmkqhy8ygl6tnl5g8tc503rwrmmrkjcq4878e0": 0.023311111,              // 0x00ca47db1be92c1072e973fd8dc4a082f7d70214
	"cyber1gs92s58t6rkallnml8ufdzrz3038dcylrcjyqd": 0.023311111,              // 0x4585c7eaa2cb96d4b59e868929efabeeb8e65b07
	"cyber1h7u5zvduvc3dqrfq9hejm35ktfxh3ha7hyqd2g": 0.0132800000000000911111, // 0x002b9c5b537a1b6004ed720f32cc808fd6210f26
	"cyber1rl3xnsrkpjfwejqfy7v4kntu64hzxy8dk65q98": 0.003111111,              // 0x00d3c9033570b8adea9c18780325a45635c55805
	"cyber1xege0g92p6exmzjv58u7vh3s5zkz75v4euzyxe": 0.003111111,              // 0x5d01f31f6eda95489ca1e3c6357a9627fa2983de
}

// 0.1.0 Distribution
func AddEulerTokensCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "add-euler-tokens <pou-file-path>",
		RunE: func(cmd *cobra.Command, args []string) error {
			//todo refactor this, see #199
			return GenerateEulerGenesisFileCmd(ctx, cdc).Execute()
		},
	}
	return cmd
}

// 0.1.1 Distribution fix for 0.1.0
// Our current distribution is 1% for proof of use, and 0.7% for others.
// In 0.1.1 launch we should:
//   - add rest 0.3% to others part. See issue #156.
//   - send 11.8000% to dead address
func AddMissingEulerTokensCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use: "add-missing-euler-tokens",
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Println("")
			fmt.Println("Fixing euler 0.1.0 network genesis distribution")

			sdk.GetConfig().SetBech32PrefixForAccount("cbd", "cbdpub")
			sdk.GetConfig().SetBech32PrefixForValidator("cbdvaloper", "cbdvaloperpub")
			sdk.GetConfig().SetBech32PrefixForConsensusNode("cbdvalcons", "cbdvalconspub")
			doc, state, err := loadGenesisState(ctx, cdc)
			if err != nil {
				return err
			}

			addrWithMissingTokens := make(map[string]int64, len(state.Accounts))
			//pocv accs
			totalPocvAddition := int64(0)
			for pocvAddr, percentage := range pocvEulerAccs {
				tokensToAdd := percentageToAmt(percentage) - pouPercentageToAmt(percentage)
				totalPocvAddition += tokensToAdd
				addrWithMissingTokens[pocvAddr] += tokensToAdd
			}
			fmt.Printf("%v added to pocv\n", totalPocvAddition)

			//pou accs
			totalPouAddition := int64(0)
			pouEulerAccs, shares, err := readPouAccounts()
			if err != nil {
				return err
			}
			for i := range pouEulerAccs {
				tokensToAdd := pouShareToAmt(shares[i]) - pouPercentageToAmt(shares[i])
				totalPouAddition += tokensToAdd
				addrWithMissingTokens[pouEulerAccs[i].String()] += tokensToAdd
			}
			fmt.Printf("%v added to pou\n", totalPouAddition)

			// burn addr acc
			addrWithMissingTokens[cbd.GetBurnAddress().String()] = percentageToAmt(tokensToBurnPercentage)

			// add tokens
			for accIndex, acc := range state.Accounts {
				tokensToAdd, ok := addrWithMissingTokens[acc.Address.String()]
				if ok {
					acc.Amount += tokensToAdd
					delete(addrWithMissingTokens, acc.Address.String())
					state.Accounts[accIndex] = acc
				}
			}

			if len(addrWithMissingTokens) != 0 {
				fmt.Println("")
				return fmt.Errorf("%v does not included in genesis.json file. Reverting", addrWithMissingTokens)
			}
			fmt.Println("Fixed!")
			return saveGenesisState(ctx, cdc, doc, state)

		},
	}
	return cmd
}

func ChangeGenesisAccsPrefixCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use: "change-genesis-accs-prefix",
		RunE: func(cmd *cobra.Command, args []string) error {

			// set old prefix for parsing file
			sdk.GetConfig().SetBech32PrefixForAccount("cbd", "cbdpub")
			sdk.GetConfig().SetBech32PrefixForValidator("cbdvaloper", "cbdvaloperpub")
			sdk.GetConfig().SetBech32PrefixForConsensusNode("cbdvalcons", "cbdvalconspub")
			doc, state, err := loadGenesisState(ctx, cdc)
			if err != nil {
				return err
			}

			fmt.Println("")
			fmt.Println("Changing accs prefix from 'cbd' to 'cyber' (0.1.0 -> 0.1.1 change)")

			// set new prefix to persis file
			app.SetPrefix()
			return saveGenesisState(ctx, cdc, doc, state)
		},
	}
	return cmd
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// calculate proof of use(pou) tokens amount.
// pou file contains pou percentage(-> sum of all addresses is 100%)
// So to get overall network percentage multiply it on pouFileCoefficient (0.7)
func pouPercentageToAmt(pouPercentage float64) int64 {
	return int64(pouPercentage / 100 * pouFileCoefficient * float64(genesisSupply))
}

func pouShareToAmt(share float64) int64 {
	return int64(share * pouFileCoefficient * float64(genesisSupply))
}

func percentageToAmt(percentage float64) int64 {
	return int64(percentage / 100 * float64(genesisSupply))
}
