package keeper

import (
	"encoding/binary"

	"github.com/cybercongress/go-cyber/util"
	"github.com/cybercongress/go-cyber/x/link/types"

	"io"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmos "github.com/tendermint/tendermint/libs/os"
)

type IndexKeeper struct {
	GraphKeeper

	// Actual links for current rank calculated state.
	currentRankInLinks  types.Links
	currentRankOutLinks types.Links

	// New links for the next rank calculation.
	nextRankInLinks  types.Links
	nextRankOutLinks types.Links

	currentBlockLinks []types.CompactLink
}

func NewIndexKeeper(gk GraphKeeper) *IndexKeeper {
	return &IndexKeeper{
		GraphKeeper: gk,
	}
}

func (i *IndexKeeper) Load(rankCtx sdk.Context, freshCtx sdk.Context) {
	inLinks, outLinks, err := i.GraphKeeper.GetAllLinks(rankCtx)
	if err != nil {
		tmos.Exit(err.Error())
	}

	i.currentRankInLinks = inLinks
	i.currentRankOutLinks = outLinks

	newInLinks, newOutLinks, err := i.GraphKeeper.GetAllLinksFiltered(freshCtx, func(l types.CompactLink) bool {
		return !i.currentRankOutLinks.IsLinkExist(l.From(), l.To(), l.Account())
	})

	if err != nil {
		tmos.Exit(err.Error())
	}

	i.nextRankInLinks = newInLinks
	i.nextRankOutLinks = newOutLinks
}

func (i *IndexKeeper) FixLinks() {
	i.currentRankInLinks.PutAll(i.nextRankInLinks)
	i.currentRankOutLinks.PutAll(i.nextRankOutLinks)

	i.nextRankInLinks = make(types.Links)
	i.nextRankOutLinks = make(types.Links)
}

// return true if this block has new links
func (i *IndexKeeper) EndBlocker() bool {
	hasNewLinks := len(i.currentBlockLinks) > 0

	for _, link := range i.currentBlockLinks {
		i.nextRankOutLinks.Put(link.From(), link.To(), link.Account())
		i.nextRankInLinks.Put(link.To(), link.From(), link.Account())
	}

	i.currentBlockLinks = make([]types.CompactLink, 0, 1000)
	return hasNewLinks
}

func (i *IndexKeeper) PutIntoIndex(link types.CompactLink) {
	i.currentBlockLinks = append(i.currentBlockLinks, link)
}

func (i *IndexKeeper) PutLink(ctx sdk.Context, link types.CompactLink) {
	if !ctx.IsCheckTx() {
		i.currentBlockLinks = append(i.currentBlockLinks, link)
	}
}

func (i *IndexKeeper) GetOutLinks() types.Links {
	return i.currentRankOutLinks
}

func (i *IndexKeeper) GetInLinks() types.Links {
	return i.currentRankInLinks
}

func (i *IndexKeeper) GetNextOutLinks() types.Links {
	return i.nextRankOutLinks
}

func (i *IndexKeeper) GetCurrentBlockLinks() []types.CompactLink {
	return i.currentBlockLinks
}

func (i *IndexKeeper) GetCurrentBlockNewLinks() []types.CompactLink {
	result := make([]types.CompactLink, 0, len(i.currentBlockLinks))
	for _, link := range i.currentBlockLinks {
		if !i.IsAnyLinkExist(link.From(), link.To()) {
			result = append(result, link)
		}
	}
	return result
}

func (i *IndexKeeper) IsAnyLinkExist(from types.CidNumber, to types.CidNumber) bool {
	return i.currentRankOutLinks.IsAnyLinkExist(from, to) || i.nextRankOutLinks.IsAnyLinkExist(from, to)
}

func (i *IndexKeeper) IsLinkExist(link types.CompactLink) bool {
	return i.currentRankOutLinks.IsLinkExist(link.From(), link.To(), link.Account()) ||
		i.nextRankOutLinks.IsLinkExist(link.From(), link.To(), link.Account())
}

func (i *IndexKeeper) LoadFromReader(ctx sdk.Context, reader io.Reader) (err error) {
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
		//fmt.Println("link: ", types.UnmarshalBinaryLink(linkBytes))
		link := types.UnmarshalBinaryLink(linkBytes)
		i.GraphKeeper.PutLink(ctx, link)
		i.PutLink(ctx, link)
	}
	//fmt.Println("Links count: ",i.GraphKeeper.GetLinksCount(ctx))
	//ctx.CacheContext()
	return
}


