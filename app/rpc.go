package app

import (
	"errors"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	//"github.com/cosmos/cosmos-sdk/x/auth/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cybercongress/go-cyber/merkle"
	ctypes "github.com/cybercongress/go-cyber/types"
	bw "github.com/cybercongress/go-cyber/x/bandwidth"
	"github.com/cybercongress/go-cyber/x/link"
)

type RankedCid struct {
	Cid  link.Cid `json:"cid"`
	Rank uint64  `json:"rank"`
}

func (app *CyberdApp) RpcContext() sdk.Context {
	return app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
}

func (app *CyberdApp) Search(cid string, page, perPage int) ([]RankedCid, int, error) {

	ctx := app.RpcContext()
	cidNumber, exists := app.graphKeeper.GetCidNumber(ctx, link.Cid(cid))
	if !exists || cidNumber > app.rankKeeper.GetLastCidNum() {
		return nil, 0, errors.New("no such cid found")
	}

	rankedCidNumbers, size, err := app.rankKeeper.Search(cidNumber, page, perPage)

	if err != nil {
		return nil, size, err
	}

	result := make([]RankedCid, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, RankedCid{Cid: app.graphKeeper.GetCid(ctx, c.GetNumber()), Rank: c.GetRank()})
	}

	return result, size, nil
}

func (app *CyberdApp) Backlinks(cid string, page, perPage int) ([]RankedCid, int, error) {

	ctx := app.RpcContext()
	cidNumber, exists := app.graphKeeper.GetCidNumber(ctx, link.Cid(cid))
	if !exists || cidNumber > app.rankKeeper.GetLastCidNum() {
		return nil, 0, errors.New("no such cid found")
	}

	rankedCidNumbers, size, err := app.rankKeeper.Backlinks(cidNumber, page, perPage)

	if err != nil {
		return nil, size, err
	}

	result := make([]RankedCid, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, RankedCid{Cid: app.graphKeeper.GetCid(ctx, c.GetNumber()), Rank: c.GetRank()})
	}

	return result, size, nil
}

func (app *CyberdApp) Top(page, perPage int) ([]RankedCid, int, error) {

	ctx := app.RpcContext()

	rankedCidNumbers, size, err := app.rankKeeper.Top(page, perPage)

	if err != nil {
		return nil, size, err
	}

	result := make([]RankedCid, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, RankedCid{Cid: app.graphKeeper.GetCid(ctx, c.GetNumber()), Rank: c.GetRank()})
	}

	return result, size, nil
}

func (app *CyberdApp) Rank(cid string, proof bool) (uint64, []merkle.Proof, error) {

	cidNumber, exists := app.graphKeeper.GetCidNumber(app.RpcContext(), link.Cid(cid))
	if !exists || cidNumber > app.rankKeeper.GetLastCidNum() {
		return 0, nil, errors.New("no such cid found")
	}

	rankValue := app.rankKeeper.GetRankValue(cidNumber)

	if proof {
		proofs := make([]merkle.Proof, 0, int64(math.Log2(float64(app.CidsCount()))))
		proofs = app.rankKeeper.GetMerkleTree().GetIndexProofs(int(cidNumber))
		return rankValue, proofs, nil
	}
	return rankValue, nil, nil
}

func (app *CyberdApp) Account(address sdk.AccAddress) exported.Account {
	return app.accountKeeper.GetAccount(app.RpcContext(), address)
}

func (app *CyberdApp) AccountBandwidth(address sdk.AccAddress) bw.AcсountBandwidth {
	return app.bandwidthMeter.GetCurrentAccountBandwidth(app.RpcContext(), address)
}

func (app *CyberdApp) AccountLinks(address sdk.AccAddress, page, perPage int) ([]link.Link, int, error) {
	ctx := app.RpcContext()

	acc := app.accountKeeper.GetAccount(ctx, address)

	if acc != nil {
		accNumber := ctypes.AccNumber(acc.GetAccountNumber())
		links, total, _ := app.rankKeeper.Accounts(uint64(accNumber), page, perPage)

		result := make([]link.Link, 0, len(links))
		for j, c := range links {
			result = append(result, link.Link{From: app.graphKeeper.GetCid(ctx, j), To: app.graphKeeper.GetCid(ctx, c)})
		}
		return result, total, nil
	} else {
		return nil, 0, nil
	}
}

func (app *CyberdApp) IsLinkExist(from link.Cid, to link.Cid, address sdk.AccAddress) bool {

	ctx := app.RpcContext()
	fromNumber, fromExist := app.graphKeeper.GetCidNumber(ctx, from)
	toNumber, toExists := app.graphKeeper.GetCidNumber(ctx, to)

	if fromExist && toExists {
		if address != nil {
			acc := app.accountKeeper.GetAccount(ctx, address)
			if acc != nil {
				accNumber := ctypes.AccNumber(acc.GetAccountNumber())
				return app.indexKeeper.IsLinkExist(link.NewLink(fromNumber, toNumber, accNumber))
			}
		} else {
			return app.indexKeeper.IsAnyLinkExist(fromNumber, toNumber)
		}
	}
	return false
}

func (app *CyberdApp) CurrentBandwidthPrice() float64 {
	return app.bandwidthMeter.GetCurrentCreditPrice()
}

func (app *CyberdApp) CidsCount() uint64 {
	return app.graphKeeper.GetCidsCount(app.RpcContext())
}

func (app *CyberdApp) LinksCount() uint64 {
	return app.graphKeeper.GetLinksCount(app.RpcContext())
}

func (app *CyberdApp) AccsCount() uint64 {
	return app.accountKeeper.GetNextAccountNumber(app.RpcContext())
}

func (app *CyberdApp) CurrentNetworkLoad() float64 {
	return app.bandwidthMeter.GetCurrentNetworkLoad(app.RpcContext())
}