package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	. "github.com/cybercongress/cyberd/app/storage"
	. "github.com/cybercongress/cyberd/app/types"
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
	return nil
}

func (app *CyberdApp) IsLinkExist(from Cid, to Cid, address sdk.AccAddress) bool {

	ctx := app.NewContext(true, abci.Header{})
	acc := app.accountKeeper.GetAccount(ctx, address)

	if acc != nil {
		fromNumber, fromExist := app.persistStorages.CidIndex.GetCidIndex(ctx, from)
		toNumber, toExists := app.persistStorages.CidIndex.GetCidIndex(ctx, to)
		if fromExist && toExists {
			accNumber := AccountNumber(acc.GetAccountNumber())
			return app.persistStorages.Links.IsLinkExist(ctx, NewLink(fromNumber, toNumber, accNumber))
		}
	}

	return false
}
