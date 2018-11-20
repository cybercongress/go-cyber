package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cbd "github.com/cybercongress/cyberd/cosmos/poc/app/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Main search entry
func (app *CyberdApp) Search(cid cbd.Cid, page, perPage int) ([]cbd.CidWithRank, int, error) {

	ctx := app.BaseApp.NewContext(false, abci.Header{})

	if cidNumber, ok := app.persistStorages.CidIndex.GetCidIndex(ctx, cid); ok {

		cidNumbersWithRanks, totalSize, err := app.rank.GetSortedByRankLinkedCids(cidNumber, page, perPage)
		if err != nil || totalSize == 0 {
			return nil, 0, err
		}

		result := make([]cbd.CidWithRank, 0, totalSize)
		for _, i := range cidNumbersWithRanks {
			cid := app.persistStorages.CidIndex.GetCid(ctx, i.CidNumber())
			result = append(result, cbd.CidWithRank{Cid: cid, Rank: i.Rank()})
		}
		return result, totalSize, nil
	}

	return nil, 0, cbd.ErrCidNotFound()
}

// Used to query addresses
func (app *CyberdApp) Account(address sdk.AccAddress) auth.Account {

	acc := app.accStorage.GetAccount(app.NewContext(true, abci.Header{}), address)

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
