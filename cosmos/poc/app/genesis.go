package app

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/cyberd/cosmos/poc/app/coin"
	"github.com/cybercongress/cyberd/cosmos/poc/app/storage"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisState reflects the genesis state of the application.
type GenesisState struct {
	Accounts []*GenesisAccount `json:"accounts"`
}

// GenesisAccount reflects a genesis account the application expects in it's
// genesis state.
type GenesisAccount struct {
	Name    string         `json:"name"`
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// initChainer implements the custom application logic that the BaseApp will
// invoke upon initialization. In this case, it will take the application's
// state provided by 'req' and attempt to deserialize said state. The state
// should contain all the genesis accounts. These accounts will be added to the
// application's account mapper.

func NewGenesisApplier(imms *storage.InMemoryStorage, cdc *codec.Codec, accStorage auth.AccountKeeper) sdk.InitChainer {

	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {

		stateJSON := req.AppStateBytes
		genesisState := new(GenesisState)
		err := cdc.UnmarshalJSON(stateJSON, genesisState)

		if err != nil {
			panic(err)
		}

		for _, gacc := range genesisState.Accounts {
			acc, err := gacc.ToBaseAccount()
			if err != nil {
				panic(err)
			}

			acc.AccountNumber = accStorage.GetNextAccountNumber(ctx)
			accStorage.SetAccount(ctx, acc)
			imms.UpdateStake(acc.Address, acc.Coins.AmountOf(coin.CBD).Int64())
		}

		return abci.ResponseInitChain{}
	}
}

// ToAppAccount converts a GenesisAccount to an AppAccount.
func (ga *GenesisAccount) ToBaseAccount() (acc *auth.BaseAccount, err error) {
	return &auth.BaseAccount{
		Address: ga.Address,
		Coins:   ga.Coins.Sort(),
	}, nil
}
