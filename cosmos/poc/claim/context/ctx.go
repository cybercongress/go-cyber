package context

import (
	cli "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cybercongress/cyberd/cosmos/poc/app"
	"github.com/cybercongress/cyberd/cosmos/poc/claim/common"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"sync"
)

type ClaimContext struct {
	Name          string
	Passphrase    string
	ChainId       string
	ClaimFrom     types.AccAddress
	Codec         *codec.Codec
	CliContext    cli.CLIContext
	ipClaims      map[string]int
	Mtx           *sync.Mutex
	Sequence      *int64
	AccountNumber int64
}

func NewClaimContext() (ClaimContext, error) {
	name := viper.GetString(common.FlagName)
	chainId := viper.GetString(common.FlagChainId)

	cdc := app.MakeCodec()

	cliCtx := newCLIContext(name, chainId, viper.GetString(common.FlagNode)).
		WithCodec(cdc).
		WithAccountDecoder(app.GetAccountDecoder(cdc))

	address, err := types.AccAddressFromBech32(viper.GetString(common.FlagAddress))
	if err != nil {
		return ClaimContext{}, err
	}

	accountNumber, err := cliCtx.GetAccountNumber(address)
	if err != nil {
		return ClaimContext{}, err
	}
	seq, err := cliCtx.GetAccountSequence(address)
	if err != nil {
		return ClaimContext{}, err
	}

	return ClaimContext{
		Name:          name,
		ClaimFrom:     address,
		Passphrase:    viper.GetString(common.FlagPassphrase),
		ChainId:       chainId,
		Codec:         cdc,
		CliContext:    cliCtx,
		ipClaims:      make(map[string]int),
		Mtx:           new(sync.Mutex),
		Sequence:      &seq,
		AccountNumber: accountNumber,
	}, nil
}

func (ctx ClaimContext) IncrementIp(ip string) error {
	cur := ctx.ipClaims[ip]
	if cur > 100 {
		return errors.New("Limit for ip exceeded")
	}
	ctx.ipClaims[ip] = cur + 1
	return nil
}

func (ctx ClaimContext) TxBuilder() (authtxb.TxBuilder, error) {

	return authtxb.TxBuilder{
		ChainID:       ctx.ChainId,
		Gas:           10000000,
		AccountNumber: ctx.AccountNumber,
		Sequence:      *ctx.Sequence,
		Fee:           "",
		Memo:          "",
		Codec:         ctx.Codec,
	}, nil
}

func newCLIContext(accName string, chainId string, nodeEndpoint string) cli.CLIContext {

	node := rpcclient.NewHTTP(nodeEndpoint, "/websocket")
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
