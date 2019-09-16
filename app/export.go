package app

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cybercongress/cyberd/types/coin"
	"github.com/cybercongress/cyberd/util"
	"github.com/cybercongress/cyberd/x/link"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// ExportAppStateAndValidators implements custom application logic that exposes
// various parts of the application's state and set of validators. An error is
// returned if any step getting the state or set of validators fails.
func (app *CyberdApp) ExportAppStateAndValidators(
	serverCtx *server.Context,
) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {

	ctx := app.NewContext(true, abci.Header{})

	err = link.WriteGenesis(ctx, app.cidNumKeeper, app.linkIndexedKeeper, app.Logger())
	if err != nil {
		return nil, nil, err
	}

	// iterate to get the accounts
	var accounts []GenesisAccount
	appendAccount := func(acc auth.Account) (stop bool) {
		account := NewGenesisAccount(acc)
		accounts = append(accounts, account)
		return false
	}
	app.accountKeeper.IterateAccounts(ctx, appendAccount)

	doc, genState, err := LoadGenesisState(serverCtx, app.cdc)
	if err != nil {
		return nil, nil, err
	}

	genState.Accounts = accounts
	genState.GenTxs = []json.RawMessage{}

	genState.Pool.NotBondedTokens = sdk.ZeroInt()
	for _, acc := range accounts {
		genState.Pool.NotBondedTokens = genState.Pool.NotBondedTokens.Add(sdk.NewInt(acc.Coins.AmountOf(coin.CYB).Int64()))
	}

	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	genesisFilePath := util.RootifyPath("export/genesis.json")
	err = util.ExportGenesisFile(genesisFilePath, doc.ChainID, doc.Validators, appState)
	if err != nil {
		return nil, nil, err
	}

	validators = staking.WriteValidators(ctx, app.stakingKeeper)
	return appState, validators, nil
}
