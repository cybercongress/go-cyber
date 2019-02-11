package app

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cbd "github.com/cybercongress/cyberd/types"
	bdwth "github.com/cybercongress/cyberd/x/bandwidth/types"
	cbdlink "github.com/cybercongress/cyberd/x/link/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type RankedCid struct {
	Cid  cbdlink.Cid `json:"cid"`
	Rank float64     `amino:"unsafe" json:"rank"`
}

func (app *CyberdApp) RpcContext() sdk.Context {
	return app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
}

func (app *CyberdApp) Search(cid string, page, perPage int) ([]RankedCid, int, error) {

	ctx := app.RpcContext()
	cidNumber, exists := app.cidNumKeeper.GetCidNumber(ctx, cbdlink.Cid(cid))
	if !exists || cidNumber > app.rankState.GetLastCidNum() {
		return nil, 0, errors.New("no such cid found")
	}

	rankedCidNumbers, size, err := app.rankState.GetCidRankedLinks(cidNumber, page, perPage)

	if err != nil {
		return nil, size, err
	}

	result := make([]RankedCid, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, RankedCid{Cid: app.cidNumKeeper.GetCid(ctx, c.GetNumber()), Rank: c.GetRank()})
	}

	return result, size, nil
}

func (app *CyberdApp) Account(address sdk.AccAddress) auth.Account {
	return app.accountKeeper.GetAccount(app.RpcContext(), address)
}

func (app *CyberdApp) AccountBandwidth(address sdk.AccAddress) bdwth.AcсBandwidth {
	return app.bandwidthMeter.GetCurrentAccBandwidth(app.RpcContext(), address)
}

func (app *CyberdApp) IsLinkExist(from cbdlink.Cid, to cbdlink.Cid, address sdk.AccAddress) bool {

	ctx := app.RpcContext()
	fromNumber, fromExist := app.cidNumKeeper.GetCidNumber(ctx, from)
	toNumber, toExists := app.cidNumKeeper.GetCidNumber(ctx, to)

	if address != nil {
		acc := app.accountKeeper.GetAccount(ctx, address)
		if fromExist && toExists && acc != nil {
			accNumber := cbd.AccNumber(acc.GetAccountNumber())
			return app.linkIndexedKeeper.IsLinkExist(ctx, cbdlink.NewLink(fromNumber, toNumber, accNumber))
		}
		return false
	}

	// link created by someone. used for request with not specified address
	return fromExist && toExists
}

func (app *CyberdApp) CurrentBandwidthPrice() float64 {
	return app.bandwidthMeter.GetCurrentCreditPrice()
}

func (app *CyberdApp) CidsCount() uint64 {
	return app.mainKeeper.GetCidsCount(app.RpcContext())
}

func (app *CyberdApp) LinksCount() uint64 {
	return app.mainKeeper.GetLinksCount(app.RpcContext())
}

func (app *CyberdApp) AccsCount() uint64 {
	return app.accountKeeper.GetNextAccountNumber(app.RpcContext())
}
