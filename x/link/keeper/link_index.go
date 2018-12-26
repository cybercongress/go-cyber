package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/cybercongress/cyberd/x/link/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type LinkIndexedKeeper struct {
	LinkKeeper

	// Actual links for current rank state.
	inLinks  map[CidNumber]CidLinks
	outLinks map[CidNumber]CidLinks

	// New links for the next calculation.
	// Actually, do we need them in memory?
	// TODO: optimize to not store whole index (store just new links)
	newInLinks  map[CidNumber]CidLinks
	newOutLinks map[CidNumber]CidLinks
}

func NewLinkIndexedKeeper(keeper LinkKeeper) LinkIndexedKeeper {
	return LinkIndexedKeeper{LinkKeeper: keeper}
}

func (i *LinkIndexedKeeper) Load(rankCtx sdk.Context, freshCtx sdk.Context) {
	inLinks, outLinks, err := i.LinkKeeper.GetAllLinks(rankCtx)
	if err != nil {
		cmn.Exit(err.Error())
	}

	i.inLinks = inLinks
	i.outLinks = outLinks

	newInLinks, newOutLinks, err := i.LinkKeeper.GetAllLinks(freshCtx)
	if err != nil {
		cmn.Exit(err.Error())
	}

	i.newInLinks = newInLinks
	i.newOutLinks = newOutLinks
}

func (i *LinkIndexedKeeper) FixLinks() {
	// todo state copied
	i.inLinks = Links(i.newInLinks).Copy()
	i.outLinks = Links(i.newOutLinks).Copy()
}

func (i *LinkIndexedKeeper) PutIntoIndex(link Link) {
	Links(i.newOutLinks).Put(link.From(), link.To(), link.Acc())
	Links(i.newInLinks).Put(link.To(), link.From(), link.Acc())
}

func (i *LinkIndexedKeeper) GetOutLinks() map[CidNumber]CidLinks {
	return i.outLinks
}

func (i *LinkIndexedKeeper) GetInLinks() map[CidNumber]CidLinks {
	return i.inLinks
}
