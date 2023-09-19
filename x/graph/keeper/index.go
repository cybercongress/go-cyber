package keeper

import (
	"encoding/binary"
	"io"

	ctypes "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/utils"
	"github.com/cybercongress/go-cyber/x/graph/types"

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

	// Inter-block cache for cyberlinks, reset on every block during Commit
	tkey sdk.StoreKey

	// currentBlockLinks []types.CompactLink
}

func NewIndexKeeper(gk GraphKeeper, tkey sdk.StoreKey) *IndexKeeper {
	return &IndexKeeper{
		GraphKeeper: gk,
		tkey:        tkey,
	}
}

func (i *IndexKeeper) LoadState(rankCtx sdk.Context, freshCtx sdk.Context) {
	inLinks, outLinks, err := i.GraphKeeper.GetAllLinks(rankCtx)
	if err != nil {
		tmos.Exit(err.Error())
	}

	i.currentRankInLinks = inLinks
	i.currentRankOutLinks = outLinks

	newInLinks, newOutLinks, err := i.GraphKeeper.GetAllLinksFiltered(freshCtx, func(l types.CompactLink) bool {
		return !i.currentRankOutLinks.IsLinkExist(types.CidNumber(l.From), types.CidNumber(l.To), ctypes.AccNumber(l.Account))
	})
	if err != nil {
		tmos.Exit(err.Error())
	}

	i.nextRankInLinks = newInLinks
	i.nextRankOutLinks = newOutLinks
}

func (i *IndexKeeper) UpdateRankLinks() {
	i.currentRankInLinks.PutAll(i.nextRankInLinks)
	i.currentRankOutLinks.PutAll(i.nextRankOutLinks)

	i.nextRankInLinks = make(types.Links)
	i.nextRankOutLinks = make(types.Links)
}

func (i *IndexKeeper) MergeContextLinks(ctx sdk.Context) {
	lenLinks := uint64(0)
	iterator := sdk.KVStorePrefixIterator(ctx.TransientStore(i.tkey), types.CyberlinkTStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		link := types.UnmarshalBinaryLink(iterator.Key()[1:])
		i.nextRankOutLinks.Put(types.CidNumber(link.From), types.CidNumber(link.To), ctypes.AccNumber(link.Account))
		i.nextRankInLinks.Put(types.CidNumber(link.To), types.CidNumber(link.From), ctypes.AccNumber(link.Account))
		lenLinks++
	}

	if lenLinks > 0 {
		store := ctx.TransientStore(i.tkey)
		store.Set(types.HasNewLinks, sdk.Uint64ToBigEndian(lenLinks))
	}
}

func (i *IndexKeeper) HasNewLinks(ctx sdk.Context) bool {
	store := ctx.TransientStore(i.tkey)
	hasLinks := store.Get(types.HasNewLinks)
	if hasLinks == nil {
		return false
	}
	return sdk.BigEndianToUint64(hasLinks) > 0
}

// Use transient store because need to commit cyberlinks to cache only when transaction is successful
func (i *IndexKeeper) PutLink(ctx sdk.Context, link types.CompactLink) {
	store := ctx.TransientStore(i.tkey)
	store.Set(types.CyberlinksTStoreKey(link.MarshalBinaryLink()), []byte{1})
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

func (i *IndexKeeper) GetCurrentBlockNewLinks(ctx sdk.Context) []types.CompactLink {
	result := make([]types.CompactLink, 0, 1000)

	iterator := sdk.KVStorePrefixIterator(ctx.TransientStore(i.tkey), types.CyberlinkTStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		link := types.UnmarshalBinaryLink(iterator.Key()[1:])
		if !i.IsAnyLinkExist(types.CidNumber(link.From), types.CidNumber(link.To)) {
			result = append(result, link)
		}
	}

	return result
}

func (i *IndexKeeper) IsAnyLinkExist(from types.CidNumber, to types.CidNumber) bool {
	return i.currentRankOutLinks.IsAnyLinkExist(from, to) || i.nextRankOutLinks.IsAnyLinkExist(from, to)
}

func (i *IndexKeeper) IsLinkExist(link types.CompactLink) bool {
	return i.currentRankOutLinks.IsLinkExist(types.CidNumber(link.From), types.CidNumber(link.To), ctypes.AccNumber(link.Account)) ||
		i.nextRankOutLinks.IsLinkExist(types.CidNumber(link.From), types.CidNumber(link.To), ctypes.AccNumber(link.Account))
}

func (i *IndexKeeper) IsLinkExistInCache(ctx sdk.Context, link types.CompactLink) bool {
	store := ctx.TransientStore(i.tkey)
	return store.Has(link.MarshalBinaryLink())
}

func (i *IndexKeeper) LoadFromReader(ctx sdk.Context, reader io.Reader) (err error) {
	linksCountBytes, err := utils.ReadExactlyNBytes(reader, LinksCountBytesSize)
	if err != nil {
		return
	}
	linksCount := binary.LittleEndian.Uint64(linksCountBytes)

	for j := uint64(0); j < linksCount; j++ {
		linkBytes, err := utils.ReadExactlyNBytes(reader, LinkBytesSize)
		if err != nil {
			return err
		}
		compactLink := types.UnmarshalBinaryLink(linkBytes)
		i.GraphKeeper.SaveLink(ctx, compactLink)
		i.PutLink(ctx, compactLink)
	}
	return
}
