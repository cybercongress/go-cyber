package app

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// todo what is purpose of this functions?
// ExportAppStateAndValidators implements custom application logic that exposes
// various parts of the application's state and set of validators. An error is
// returned if any step getting the state or set of validators fails.
func (app *CyberdApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})

	// iterate to get the accounts
	var accounts []GenesisAccount
	appendAccount := func(acc auth.Account) (stop bool) {
		account := NewGenesisAccountI(acc)
		accounts = append(accounts, account)
		return false
	}
	app.accountKeeper.IterateAccounts(ctx, appendAccount)
	genState := NewGenesisState(
		accounts,
		stake.ExportGenesis(ctx, app.stakeKeeper),
		slashing.ExportGenesis(ctx, app.slashingKeeper),
	)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}
	validators = stake.WriteValidators(ctx, app.stakeKeeper)
	return appState, validators, nil
}
