package client

import (
	"errors"
	"fmt"
	cli "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/utils"
	cskeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cybercongress/cyberd/app"
	"github.com/cybercongress/cyberd/daemon/rpc"
	bwtps "github.com/cybercongress/cyberd/x/bandwidth/types"
	"github.com/cybercongress/cyberd/x/link"
	cbdlink "github.com/cybercongress/cyberd/x/link/types"
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
	cliCtx      *cli.CLIContext
	txBuilder   *authtxb.TxBuilder
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
	seq, err := cliCtx.GetAccountSequence(addr)

	if err != nil {
		panic(err)
	}

	txBuilder := authtxb.TxBuilder{
		Gas:           10000000000,
		ChainID:       status.NodeInfo.Network,
		AccountNumber: accountNumber,
		TxEncoder:     utils.GetTxEncoder(cdc),
		Sequence:      seq,
	}

	return HttpCyberdClient{
		tdmClient:  tdmHttpClient,
		httpClient: *httpClient,

		chainId: status.NodeInfo.Network,

		passphrase:  passphrase,
		fromAddress: addr,
		cliCtx:      &cliCtx,
		txBuilder:   &txBuilder,
	}
}

func (c HttpCyberdClient) GetChainId() string {
	return c.chainId
}

func (c HttpCyberdClient) IsLinkExist(from cbdlink.Cid, to cbdlink.Cid, addr sdk.AccAddress) (result bool, err error) {
	_, err = c.httpClient.Call("is_link_exist",
		map[string]interface{}{"from": from, "to": to, "address": addr.String()},
		&result,
	)
	return
}

func (c HttpCyberdClient) GetCurrentBandwidthCreditPrice() (float64, error) {
	result := &rpc.ResultBandwidthPrice{}
	_, err := c.httpClient.Call("current_bandwidth_price", map[string]interface{}{}, &result)
	return result.Price, err
}

func (c HttpCyberdClient) GetAccountBandwidth() (result bwtps.AcсBandwidth, err error) {
	_, err = c.httpClient.Call("account_bandwidth",
		map[string]interface{}{"address": c.fromAddress.String()}, &result)
	return
}

func (c HttpCyberdClient) SubmitLinkSync(link Link) error {
	return c.SubmitLinksSync([]Link{link})
}

func (c HttpCyberdClient) SubmitLinksSync(links []Link) error {

	// used to remove duplicated items
	var filter = make(CidsFilter)
	msges := make([]sdk.Msg, 0)

	for _, l := range links {

		if filter.Contains(l.From, l.To) {
			continue
		}

		exists, err := c.IsLinkExist(l.From, l.To, c.fromAddress)
		if err != nil {
			return err
		}
		if !exists {
			msges = append(msges, link.NewMsg(c.fromAddress, l.From, l.To))
		}
		filter.Put(l.From, l.To)
	}
	return c.BroadcastTx(msges)
}

func (c HttpCyberdClient) BroadcastTx(msgs []sdk.Msg) error {

	if len(msgs) == 0 {
		return nil
	}

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
	c.txBuilder.Sequence = c.txBuilder.Sequence + 1
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
