package app

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cybercongress/cyberd/x/link"
	"github.com/cybercongress/cyberd/x/mint"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

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
		auth.NewGenesisState(sdk.Coins{}, app.accountKeeper.GetParams(ctx)),
		staking.ExportGenesis(ctx, app.stakingKeeper),
		mint.ExportGenesis(ctx, app.minter),
		distr.ExportGenesis(ctx, app.distrKeeper),
		slashing.ExportGenesis(ctx, app.slashingKeeper),
	)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}
	validators = staking.WriteValidators(ctx, app.stakingKeeper)
	err = link.WriteGenesis(ctx, app.cidNumKeeper, app.linkIndexedKeeper, app.Logger)
	if err != nil {
		return nil, nil, err
	}
	return appState, validators, nil
}
