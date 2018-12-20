package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/cybercongress/cyberd/x/link/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type LinkIndexedKeeper struct {
	LinkKeeper

	inLinks  map[CidNumber]CidLinks
	outLinks map[CidNumber]CidLinks
}

func NewLinkIndexedKeeper(keeper LinkKeeper) LinkIndexedKeeper {
	return LinkIndexedKeeper{LinkKeeper: keeper}
}

func (i *LinkIndexedKeeper) Load(ctx sdk.Context) {
	inLinks, outLinks, err := i.LinkKeeper.GetAllLinks(ctx)
	if err != nil {
		cmn.Exit(err.Error())
	}

	i.inLinks = inLinks
	i.outLinks = outLinks
}

func (i *LinkIndexedKeeper) Empty() {
	i.inLinks = make(map[CidNumber]CidLinks)
	i.outLinks = make(map[CidNumber]CidLinks)
}

func (i *LinkIndexedKeeper) PutIntoIndex(link Link) {
	Links(i.outLinks).Put(link.From(), link.To(), link.Acc())
	Links(i.inLinks).Put(link.To(), link.From(), link.Acc())
}

func (i *LinkIndexedKeeper) GetOutLinks() map[CidNumber]CidLinks {
	return i.outLinks
}

func (i *LinkIndexedKeeper) GetInLinks() map[CidNumber]CidLinks {
	return i.inLinks
}
