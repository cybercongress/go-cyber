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
	// TODO: init???
	newLinks []Link
}

func NewLinkIndexedKeeper(keeper LinkKeeper) LinkIndexedKeeper {
	return LinkIndexedKeeper{LinkKeeper: keeper}
}

// TODO: rewrite to load from previous version to current links and older versions to new links.
func (i *LinkIndexedKeeper) Load(ctx sdk.Context) {
	inLinks, outLinks, err := i.LinkKeeper.GetAllLinks(ctx)
	if err != nil {
		cmn.Exit(err.Error())
	}

	i.inLinks = inLinks
	i.outLinks = outLinks
}

func (i *LinkIndexedKeeper) MergeNewLinks() {
	for _, l := range i.newLinks {
		Links(i.outLinks).Put(l.From(), l.To(), l.Acc())
		Links(i.inLinks).Put(l.To(), l.From(), l.Acc())
	}
	i.newLinks = make([]Link, 0)
}

func (i *LinkIndexedKeeper) PutIntoIndex(link Link) {
	i.newLinks = append(i.newLinks, link) //TODO: optimize addition
}

func (i *LinkIndexedKeeper) GetOutLinks() map[CidNumber]CidLinks {
	return i.outLinks
}

func (i *LinkIndexedKeeper) GetInLinks() map[CidNumber]CidLinks {
	return i.inLinks
}
