package keeper

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/util"
	. "github.com/cybercongress/cyberd/x/link/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"io"
)

type LinkIndexedKeeper struct {
	LinkKeeper

	// Actual links for current rank calculated state.
	currentRankInLinks  Links
	currentRankOutLinks Links

	// New links for the next rank calculation.
	nextRankInLinks  Links
	nextRankOutLinks Links

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
	i.currentRankOutLinks = outLinks

	newInLinks, newOutLinks, err := i.LinkKeeper.GetAllLinksFiltered(freshCtx, func(l CompactLink) bool {
		return !i.currentRankOutLinks.IsLinkExist(l.From(), l.To(), l.Acc())
	})

	if err != nil {
		cmn.Exit(err.Error())
	}

	i.nextRankInLinks = newInLinks
	i.nextRankOutLinks = newOutLinks
}

func (i *LinkIndexedKeeper) FixLinks() {
	i.currentRankInLinks.PutAll(i.nextRankInLinks)
	i.currentRankOutLinks.PutAll(i.nextRankOutLinks)
	i.nextRankInLinks = make(Links)
	i.nextRankOutLinks = make(Links)
}

// return true if this block has new links
func (i *LinkIndexedKeeper) EndBlocker() bool {
	hasNewLinks := len(i.currentBlockLinks) > 0
	for _, link := range i.currentBlockLinks {
		i.nextRankOutLinks.Put(link.From(), link.To(), link.Acc())
		i.nextRankInLinks.Put(link.To(), link.From(), link.Acc())
	}
	i.currentBlockLinks = make([]CompactLink, 0, 1000) // todo: 1000 hardcoded value
	return hasNewLinks
}

func (i *LinkIndexedKeeper) PutIntoIndex(link CompactLink) {
	i.currentBlockLinks = append(i.currentBlockLinks, link)
}

func (i *LinkIndexedKeeper) PutLink(ctx sdk.Context, link CompactLink) {
	if !ctx.IsCheckTx() {
		i.currentBlockLinks = append(i.currentBlockLinks, link)
	}
	i.LinkKeeper.PutLink(ctx, link)
}

func (i *LinkIndexedKeeper) GetOutLinks() Links {
	return i.currentRankOutLinks
}

func (i *LinkIndexedKeeper) GetInLinks() Links {
	return i.currentRankInLinks
}

func (i *LinkIndexedKeeper) GetNextOutLinks() Links {
	return i.nextRankOutLinks
}

func (i *LinkIndexedKeeper) GetCurrentBlockLinks() []CompactLink {
	return i.currentBlockLinks
}

func (i *LinkIndexedKeeper) GetCurrentBlockNewLinks() []CompactLink {
	result := make([]CompactLink, 0, len(i.currentBlockLinks))
	for _, link := range i.currentBlockLinks {
		if !i.IsAnyLinkExist(link.From(), link.To()) {
			result = append(result, link)
		}
	}
	return result
}

func (i *LinkIndexedKeeper) IsAnyLinkExist(from CidNumber, to CidNumber) bool {
	return i.currentRankOutLinks.IsAnyLinkExist(from, to) || i.nextRankOutLinks.IsAnyLinkExist(from, to)
}

func (i *LinkIndexedKeeper) IsLinkExist(link CompactLink) bool {
	return i.currentRankOutLinks.IsLinkExist(link.From(), link.To(), link.Acc()) ||
		i.nextRankOutLinks.IsLinkExist(link.From(), link.To(), link.Acc())
}

//todo: remove duplicated method (BaseLinksKeeper)
func (i *LinkIndexedKeeper) LoadFromReader(ctx sdk.Context, reader io.Reader) (err error) {
	linksCountBytes, err := util.ReadExactlyNBytes(reader, LinksCountBytesSize)
	if err != nil {
		return
	}
	linksCount := binary.LittleEndian.Uint64(linksCountBytes)

	for j := uint64(0); j < linksCount; j++ {
		linkBytes, err := util.ReadExactlyNBytes(reader, LinkBytesSize)
		if err != nil {
			return err
		}
		i.PutLink(ctx, UnmarshalBinaryLink(linkBytes))
	}
	return
}
