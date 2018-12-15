package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	. "github.com/cybercongress/cyberd/app/storage"
	cbd "github.com/cybercongress/cyberd/app/types"
	"github.com/cybercongress/cyberd/x/bandwidth"
	bdwth "github.com/cybercongress/cyberd/x/bandwidth/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (app *CyberdApp) RpcContext() sdk.Context {
	return app.NewContext(true, abci.Header{})
}

func (app *CyberdApp) Search(cid string, page, perPage int) ([]RankedCid, int, error) {
	return app.memStorage.GetCidRankedLinks(cbd.Cid(cid), page, perPage)
}

func (app *CyberdApp) Account(address sdk.AccAddress) auth.Account {
	return app.accountKeeper.GetAccount(app.RpcContext(), address)
}

func (app *CyberdApp) AccountBandwidth(address sdk.AccAddress) (bdwth.AccountBandwidth, error) {
	accBdwth, err := app.accBandwidthKeeper.GetAccountBandwidth(app.RpcContext(), address)
	if err != nil {
		accBdwth.Recover(app.latestBlockHeight+1, bandwidth.RecoveryPeriod)
	}
	return accBdwth, err
}

func (app *CyberdApp) IsLinkExist(from cbd.Cid, to cbd.Cid, address sdk.AccAddress) bool {

	ctx := app.RpcContext()
	acc := app.accountKeeper.GetAccount(ctx, address)

	if acc != nil {
		fromNumber, fromExist := app.persistStorages.CidIndex.GetCidIndex(ctx, from)
		toNumber, toExists := app.persistStorages.CidIndex.GetCidIndex(ctx, to)
		if fromExist && toExists {
			accNumber := cbd.AccountNumber(acc.GetAccountNumber())
			return app.persistStorages.Links.IsLinkExist(ctx, cbd.NewLink(fromNumber, toNumber, accNumber))
		}
	}

	return false
}
