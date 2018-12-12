package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	cli "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	cskeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cybercongress/cyberd/app"
	cbd "github.com/cybercongress/cyberd/app/types"
	"github.com/cybercongress/cyberd/claim/common"
	"github.com/cybercongress/cyberd/x/link"
	"github.com/spf13/viper"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"os"
	"time"
)

func InitAddLink() func([]Link) {

	app.SetPrefix()
	chainId := viper.GetString(common.FlagChainId)
	address := viper.GetString(common.FlagAddress)
	addr, name := accountFromAddress(address)

	cdc := app.MakeCodec()
	cliCtx := newCLIContext(name, chainId).
		WithCodec(cdc).
		WithAccountDecoder(cdc)

	accountNumber, _ := cliCtx.GetAccountNumber(addr)
	seq, _ := cliCtx.GetAccountSequence(addr)

	txCtx := authtxb.TxBuilder{
		ChainID:       chainId,
		Gas:           10000000,
		AccountNumber: accountNumber,
		Sequence:      seq,
		Fee:           "",
		Memo:          "",
		Codec:         cdc,
	}

	return func(links []Link) {

		msges := make([]sdk.Msg, 0, len(links))
		for _, l := range links {
			msges = append(msges, link.NewMsg(addr, cbd.Cid(l.from), cbd.Cid(l.to)))
		}

		sendTx(address, txCtx, cliCtx, msges)
		txCtx.Sequence++
	}
}

func sendTx(address string, txCtx authtxb.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) {

	passphrase := viper.GetString(common.FlagPassphrase)
	txBytes, err := txCtx.BuildAndSign(cliCtx.From, passphrase, msgs)
	if err != nil {
		panic(err)
	}

	result, err := cliCtx.BroadcastTxSync(txBytes)

	if err != nil {
		println("Error during broadcasting tx. Rebrodcasting ...")
		println(err.Error())
		sendTx(address, txCtx, cliCtx, msgs)
	}

	if result.Code != 0 {
		println("Error during broadcasting tx")
		println(string(result.Log))
		time.Sleep(5 * time.Second)
		addr, _ := accountFromAddress(address)
		seq, _ := cliCtx.GetAccountSequence(addr)
		txCtx.Sequence = seq
		sendTx(address, txCtx, cliCtx, msgs)
	}
}

func newCLIContext(accName string, chainId string) cli.CLIContext {

	nodeUrl := viper.GetString(common.FlagNode)
	node := rpcclient.NewHTTP(nodeUrl, "/websocket")
	verifier := &common.NoopVerifier{ChainId: chainId}

	return cli.CLIContext{
		Client:        node,
		NodeURI:       "",
		AccountStore:  "acc",
		From:          accName,
		Height:        0,
		TrustNode:     true,
		UseLedger:     false,
		Async:         false,
		JSON:          false,
		PrintResponse: true,
		Verifier:      verifier,
	}
}

func accountFromAddress(from string) (fromAddr sdk.AccAddress, fromName string) {
	if from == "" {
		return nil, ""
	}

	keybase, err := keys.GetKeyBase()
	if err != nil {
		fmt.Println("no keybase found")
		os.Exit(1)
	}

	var info cskeys.Info
	if addr, err := sdk.AccAddressFromBech32(from); err == nil {
		info, err = keybase.GetByAddress(addr)
		if err != nil {
			fmt.Printf("could not find key %s\n", from)
			os.Exit(1)
		}
	} else {
		info, err = keybase.Get(from)
		if err != nil {
			fmt.Printf("could not find key %s\n", from)
			os.Exit(1)
		}
	}

	fromAddr = info.GetAddress()
	fromName = info.GetName()
	return
}
