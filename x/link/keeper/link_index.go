package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/cybercongress/cyberd/x/link/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type LinkIndexedKeeper struct {
	LinkKeeper

	// Actual links for current rank calculated state.
	currentRankInLinks  map[CidNumber]CidLinks
	currentRankOurLinks map[CidNumber]CidLinks

	// New links for the next rank calculation.
	// Actually, do we need them in memory?
	// TODO: optimize to not store whole index (store just new links)
	nextRankInLinks  map[CidNumber]CidLinks
	nextRankOutLinks map[CidNumber]CidLinks

	currentBlockLinks []CompactLink
}

func NewLinkIndexedKeeper(keeper LinkKeeper) *LinkIndexedKeeper {
	return &LinkIndexedKeeper{LinkKeeper: keeper}
}

func (i *LinkIndexedKeeper) Load(rankCtx sdk.Context, freshCtx sdk.Context) {
	inLinks, outLinks, err := i.LinkKeeper.GetAllLinks(rankCtx)
	if err != nil {
		cmn.Exit(err.Error())
	}

	i.currentRankInLinks = inLinks
	i.currentRankOurLinks = outLinks

	newInLinks, newOutLinks, err := i.LinkKeeper.GetAllLinks(freshCtx)
	if err != nil {
		cmn.Exit(err.Error())
	}

	i.nextRankInLinks = newInLinks
	i.nextRankOutLinks = newOutLinks
}

func (i *LinkIndexedKeeper) FixLinks() {
	// todo state copied
	i.currentRankInLinks = Links(i.nextRankInLinks).Copy()
	i.currentRankOurLinks = Links(i.nextRankOutLinks).Copy()
}

func (i *LinkIndexedKeeper) EndBlocker() {
	for _, link := range i.currentBlockLinks {
		Links(i.nextRankOutLinks).Put(link.From(), link.To(), link.Acc())
		Links(i.nextRankInLinks).Put(link.To(), link.From(), link.Acc())
	}
	i.currentBlockLinks = make([]CompactLink, 0, 1000) // todo: 1000 hardcoded value
}

func (i *LinkIndexedKeeper) PutIntoIndex(link CompactLink) {
	i.currentBlockLinks = append(i.currentBlockLinks, link)
}

func (i *LinkIndexedKeeper) GetOutLinks() map[CidNumber]CidLinks {
	return i.currentRankOurLinks
}

func (i *LinkIndexedKeeper) GetInLinks() map[CidNumber]CidLinks {
	return i.currentRankInLinks
}

func (i *LinkIndexedKeeper) GetCurrentBlockLinks() []CompactLink {
	return i.currentBlockLinks
}
