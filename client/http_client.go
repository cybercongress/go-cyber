package client

import (
	"errors"
	"fmt"
	cli "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	cskeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cybercongress/cyberd/app"
	cbd "github.com/cybercongress/cyberd/app/types"
	"github.com/cybercongress/cyberd/x/link"
	tdmClient "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/lib/client"
	"os"
)

type HttpCyberdClient struct {
	// tdm client
	tdmClient tdmClient.Client
	// transport client
	httpClient rpcclient.JSONRPCClient

	// general fields
	nodeUrl string
	chainId string

	// fields used by local keys store to sing transactions
	passphrase  string
	fromAddress sdk.AccAddress
	cliCtx      cli.CLIContext
	txBuilder   authtxb.TxBuilder
}

func NewHttpCyberdClient(nodeUrl string, passphrase string, singAddr string) CyberdClient {

	tdmHttpClient := tdmClient.NewHTTP(nodeUrl, "/websocket")
	httpClient := rpcclient.NewJSONRPCClient(nodeUrl)
	status, err := tdmHttpClient.Status()
	if err != nil {
		panic(err)
	}

	cdc := app.MakeCodec()
	app.SetPrefix()
	addr, cliAddrName := accountFromAddress(singAddr)
	verifier := &NoopVerifier{ChainId: status.NodeInfo.Network}
	cliCtx := cli.CLIContext{
		Client:        tdmHttpClient,
		NodeURI:       nodeUrl,
		AccountStore:  "acc",
		From:          cliAddrName,
		TrustNode:     true,
		Async:         false,
		PrintResponse: true,
		Verifier:      verifier,
	}.WithCodec(cdc).WithAccountDecoder(cdc)

	accountNumber, _ := cliCtx.GetAccountNumber(addr)
	txBuilder := authtxb.TxBuilder{
		Gas:           1000000,
		ChainID:       status.NodeInfo.Network,
		AccountNumber: accountNumber,
		Codec:         cdc,
	}

	return HttpCyberdClient{
		tdmClient:  tdmHttpClient,
		httpClient: *httpClient,

		chainId: status.NodeInfo.Network,

		passphrase:  passphrase,
		fromAddress: addr,
		cliCtx:      cliCtx,
		txBuilder:   txBuilder,
	}
}

func (c HttpCyberdClient) GetChainId() string {
	return c.chainId
}

/*func (c HttpCyberdClient) GetCurrentBandwidthCreditPrice() (float64, error) {

}

func (c HttpCyberdClient) GetAccount(address sdk.AccAddress) (auth.Account, error) {

}

func (c HttpCyberdClient) GetAccountBandwidth(address sdk.AccAddress) (bdwth.AcсBandwidth, error) {

}*/

func (c HttpCyberdClient) SubmitLinkSync(link Link) error {
	return c.SubmitLinksSync([]Link{link})
}

func (c HttpCyberdClient) SubmitLinksSync(links []Link) error {
	msges := make([]sdk.Msg, 0, len(links))
	for _, l := range links {
		msges = append(msges, link.NewMsg(c.fromAddress, cbd.Cid(l.From), cbd.Cid(l.To)))
	}
	return c.BroadcastTx(msges)
}

func (c HttpCyberdClient) BroadcastTx(msgs []sdk.Msg) error {

	seq, err := c.cliCtx.GetAccountSequence(c.fromAddress)
	if err != nil {
		return err
	}
	c.txBuilder.Sequence = seq

	txBytes, err := c.txBuilder.BuildAndSign(c.cliCtx.From, c.passphrase, msgs)
	if err != nil {
		panic(err)
	}

	result, err := c.cliCtx.BroadcastTxSync(txBytes)
	if err != nil {
		println("Error during broadcasting tx. Rebrodcasting ...")
		println(err.Error())

		_ = c.BroadcastTx(msgs)
	}

	if result.Code != 0 {
		return errors.New(string(result.Log))
	}
	return nil
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
