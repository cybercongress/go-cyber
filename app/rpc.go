package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	. "github.com/cybercongress/cyberd/cosmos/poc/app/storage"
	. "github.com/cybercongress/cyberd/cosmos/poc/app/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (app *CyberdApp) Search(cid string, page, perPage int) ([]RankedCid, int, error) {
	result, totalSize, err := app.memStorage.GetCidRankedLinks(Cid(cid), page, perPage)
	if err != nil {
		return nil, totalSize, err
	}
	return result, totalSize, nil
}

func (app *CyberdApp) Account(address sdk.AccAddress) auth.Account {

	acc := app.accountKeeper.GetAccount(app.NewContext(true, abci.Header{}), address)

	if acc != nil {
		return acc
	}

	// no acc in chain, assume new one, so balance 0 and seq 0
	// acc number -1 is for addresses not in chain
	return &auth.BaseAccount{
		Address:       address,
		Sequence:      0,
		Coins:         make(sdk.Coins, 0),
		AccountNumber: -1,
	}
}
