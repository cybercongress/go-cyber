package app

import (
	"errors"
	"github.com/cybercongress/cyberd/x/rank"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/cyberd/merkle"
	cbd "github.com/cybercongress/cyberd/types"
	bdwth "github.com/cybercongress/cyberd/x/bandwidth/types"
	cbdlink "github.com/cybercongress/cyberd/x/link/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type RankedCid struct {
	Cid      cbdlink.Cid      `json:"cid"`
	Rank     float64          `amino:"unsafe" json:"rank"`
	Accounts []sdk.AccAddress `json:"accounts"`
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

	rankedCidNumbers, size, err := app.rankState.Search(cidNumber, page, perPage)
	if err != nil {
		return nil, size, err
	}

	accNumbers := make(map[rank.RankedCidNumber][]cbd.AccNumber)
	app.linkIndexedKeeper.Iterate(ctx, func(link cbdlink.CompactLink) {
		for _, c := range rankedCidNumbers {
			if link.From() == cidNumber && link.To() == c.GetNumber() {
				accNumbers[c] = append(accNumbers[c], link.Acc())
				break
			}
		}
	})

	result := make([]RankedCid, 0, len(rankedCidNumbers))
	for c, an := range accNumbers {
		rc := RankedCid{
			Cid:      app.cidNumKeeper.GetCid(ctx, c.GetNumber()),
			Rank:     c.GetRank(),
			Accounts: make([]sdk.AccAddress, 0, len(an)),
		}

		app.accountKeeper.IterateAccounts(ctx, func(account auth.Account) (stop bool) {
			for _, an := range an {
				if cbd.AccNumber(account.GetAccountNumber()) == an {
					rc.Accounts = append(rc.Accounts, account.GetAddress())
					return true
				}
			}
			return false
		})
		result = append(result, rc)
	}

	return result, size, nil
}

func (app *CyberdApp) Rank(cid string, proof bool) (float64, []merkle.Proof, error) {

	cidNumber, exists := app.cidNumKeeper.GetCidNumber(app.RpcContext(), cbdlink.Cid(cid))
	if !exists || cidNumber > app.rankState.GetLastCidNum() {
		return 0, nil, errors.New("no such cid found")
	}

	rankValue := app.rankState.GetRankValue(cidNumber)

	if proof {
		proofs := make([]merkle.Proof, 0, int64(math.Log2(float64(app.CidsCount()))))
		proofs = app.rankState.GetMerkleTree().GetIndexProofs(int(cidNumber))
		return rankValue, proofs, nil
	}
	return rankValue, nil, nil
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

	if fromExist && toExists {
		if address != nil {
			acc := app.accountKeeper.GetAccount(ctx, address)
			if acc != nil {
				accNumber := cbd.AccNumber(acc.GetAccountNumber())
				return app.linkIndexedKeeper.IsLinkExist(cbdlink.NewLink(fromNumber, toNumber, accNumber))
			}
		} else {
			return app.linkIndexedKeeper.IsAnyLinkExist(fromNumber, toNumber)
		}
	}
	return false
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
