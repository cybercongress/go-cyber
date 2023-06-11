package types

import (
	"errors"
	"sort"
	"time"

	// "time"

	"github.com/cometbft/cometbft/libs/log"

	graphtypes "github.com/cybercongress/go-cyber/v2/x/graph/types"
)

type BaseSearchIndex struct {
	links     []cidLinks
	backlinks []cidLinks
	rank      Rank

	linksChan chan graphtypes.CompactLink
	rankChan  chan Rank
	errChan   chan error

	locked bool

	logger log.Logger
}

func NewBaseSearchIndex(log log.Logger) *BaseSearchIndex {
	return &BaseSearchIndex{
		linksChan: make(chan graphtypes.CompactLink, 1000),
		rankChan:  make(chan Rank, 1),
		errChan:   make(chan error),
		locked:    true,
		logger:    log,
	}
}

func (i *BaseSearchIndex) Run() GetError {
	go i.startListenNewLinks()
	go i.startListenNewRank()

	return i.checkIndexError
}

// LoadState links with zero rank values. No sorting. Index should be unavailable for read
func (i *BaseSearchIndex) Load(links graphtypes.Links) {
	startTime := time.Now()
	i.lock() // lock index for read

	i.links = make([]cidLinks, 0, 1000000)
	i.backlinks = make([]cidLinks, 0, 1000000)

	for from, toCids := range links {
		i.extendIndex(uint64(from))

		for to := range toCids {
			i.putLinkIntoIndex(from, to)

			i.extendReverseIndex(uint64(to))
			i.putBacklinkIntoIndex(from, to)
		}
	}

	i.logger.Info("The node search index is loaded", "time", time.Since(startTime))
}

func (i *BaseSearchIndex) PutNewLinks(links []graphtypes.CompactLink) {
	for _, link := range links {
		i.linksChan <- link
	}
}

func (i *BaseSearchIndex) PutNewRank(rank Rank) {
	i.rankChan <- rank.CopyWithoutTree()
}

func (i *BaseSearchIndex) Search(cidNumber graphtypes.CidNumber, page, perPage uint32) ([]RankedCidNumber, uint32, error) {
	i.logger.Info("Search query", "particle", cidNumber, "page", page, "perPage", perPage)

	if i.locked {
		return nil, 0, errors.New("search index currently unavailable after node restart")
	}

	if uint64(cidNumber) >= uint64(len(i.links)) {
		return []RankedCidNumber{}, 0, nil
	}

	cidLinks := i.links[cidNumber]
	if cidLinks.sortedLinks == nil || len(cidLinks.sortedLinks) == 0 {
		return []RankedCidNumber{}, 0, nil
	}

	totalSize := uint32(len(cidLinks.sortedLinks))
	startIndex := page * perPage
	if startIndex >= totalSize {
		return nil, totalSize, errors.New("page not found")
	}

	endIndex := startIndex + perPage
	if endIndex > totalSize {
		endIndex = startIndex + (totalSize % perPage)
	}

	resultSet := cidLinks.sortedLinks[startIndex:endIndex]

	return resultSet, totalSize, nil
}

func (i *BaseSearchIndex) Backlinks(cidNumber graphtypes.CidNumber, page, perPage uint32) ([]RankedCidNumber, uint32, error) {
	i.logger.Info("Backlinks query", "cid", cidNumber, "page", page, "perPage", perPage)

	if i.locked {
		return nil, 0, errors.New("the search index is currently unavailable after node restart")
	}

	if uint64(cidNumber) >= uint64(len(i.backlinks)) {
		return []RankedCidNumber{}, 0, nil
	}

	cidLinks := i.backlinks[cidNumber]
	if cidLinks.sortedLinks == nil || len(cidLinks.sortedLinks) == 0 {
		return []RankedCidNumber{}, 0, nil
	}

	totalSize := uint32(len(cidLinks.sortedLinks))
	startIndex := page * perPage
	if startIndex >= totalSize {
		return nil, totalSize, errors.New("page not found")
	}

	endIndex := startIndex + perPage
	if endIndex > totalSize {
		endIndex = startIndex + (totalSize % perPage)
	}

	resultSet := cidLinks.sortedLinks[startIndex:endIndex]

	return resultSet, totalSize, nil
}

func (i *BaseSearchIndex) Top(page, perPage uint32) ([]RankedCidNumber, uint32, error) {
	if i.locked {
		return nil, 0, errors.New("the search index is currently unavailable after node restart")
	}

	totalSize := uint32(len(i.rank.TopCIDs))
	startIndex := page * perPage
	if startIndex >= totalSize {
		return nil, totalSize, errors.New("page not found")
	}

	endIndex := startIndex + perPage
	if endIndex > totalSize {
		endIndex = startIndex + (totalSize % perPage)
	}

	resultSet := i.rank.TopCIDs[startIndex:endIndex]

	return resultSet, totalSize, nil
}

// make sure that this link (from-to) is new
func (i *BaseSearchIndex) handleLink(link graphtypes.CompactLink) {
	i.extendIndex(link.From)

	fromIndex := i.links[link.From]
	// in case unlock signal received we could operate on this index otherwise put link in the end of queue and finish
	select {
	case <-fromIndex.unlockSignal:
		i.putLinkIntoIndex(graphtypes.CidNumber(link.From), graphtypes.CidNumber(link.To))
		fromIndex.Unlock()
		break
	default:
		i.linksChan <- link
	}
}

func (i *BaseSearchIndex) handleBacklink(link graphtypes.CompactLink) {
	i.extendReverseIndex(link.To)

	toIndex := i.backlinks[link.To]
	// in case unlock signal received we could operate on this index otherwise put link in the end of queue and finish
	select {
	case <-toIndex.unlockSignal:
		i.putBacklinkIntoIndex(graphtypes.CidNumber(link.From), graphtypes.CidNumber(link.To))
		toIndex.Unlock()
		break
	default:
		i.linksChan <- link
	}
}

func (i *BaseSearchIndex) GetRankValue(cid graphtypes.CidNumber) uint64 {
	if i.rank.RankValues == nil || uint64(len(i.rank.RankValues)) <= uint64(cid) {
		return 0
	}
	return i.rank.RankValues[cid]
}

func (i *BaseSearchIndex) extendIndex(fromCidNumber uint64) {
	indexLen := uint64(len(i.links))
	if fromCidNumber >= indexLen {
		for j := indexLen; j <= fromCidNumber; j++ {
			links := NewCidLinks()
			links.Unlock() // allow operations on this index
			i.links = append(i.links, links)
		}
	}
}

func (i *BaseSearchIndex) extendReverseIndex(fromCidNumber uint64) {
	indexLen := uint64(len(i.backlinks))
	if fromCidNumber >= indexLen {
		for j := indexLen; j <= fromCidNumber; j++ {
			backlinks := NewCidLinks()
			backlinks.Unlock() // allow operations on this index
			i.backlinks = append(i.backlinks, backlinks)
		}
	}
}

func (i *BaseSearchIndex) putLinkIntoIndex(from graphtypes.CidNumber, to graphtypes.CidNumber) {
	fromLinks := i.links[uint64(from)].sortedLinks
	rankedTo := RankedCidNumber{to, i.GetRankValue(to)}
	pos := sort.Search(len(fromLinks), func(i int) bool { return fromLinks[i].rank < rankedTo.rank })
	fromLinks = append(fromLinks, RankedCidNumber{})
	copy(fromLinks[pos+1:], fromLinks[pos:])
	fromLinks[pos] = rankedTo
	i.links[uint64(from)].sortedLinks = fromLinks
}

func (i *BaseSearchIndex) putBacklinkIntoIndex(from graphtypes.CidNumber, to graphtypes.CidNumber) {
	toLinks := i.backlinks[uint64(to)].sortedLinks
	rankedFrom := RankedCidNumber{from, i.GetRankValue(from)}
	pos := sort.Search(len(toLinks), func(i int) bool { return toLinks[i].rank < rankedFrom.rank })
	toLinks = append(toLinks, RankedCidNumber{})
	copy(toLinks[pos+1:], toLinks[pos:])
	toLinks[pos] = rankedFrom
	i.backlinks[uint64(to)].sortedLinks = toLinks
}

// for parallel usage
func (i *BaseSearchIndex) startListenNewLinks() {
	defer func() {
		if r := recover(); r != nil {
			i.errChan <- r.(error)
		}
	}()

	// i.logger.Info("The search index is starting to listen to new links")
	for {
		link := <-i.linksChan
		i.handleLink(link)
		i.handleBacklink(link)
	}
}

// for parallel usage
func (i *BaseSearchIndex) startListenNewRank() {
	defer func() {
		if r := recover(); r != nil {
			i.errChan <- r.(error)
		}
	}()

	// i.logger.Info("The search index is starting to listen to new rank")
	for {
		rank := <-i.rankChan // TODO could be problems if recalculation lasts more than rank period
		i.rank = rank
		i.recalculateIndices()
	}
}

func (i *BaseSearchIndex) recalculateIndices() {
	defer i.unlock()
	n := len(i.links) // TODO: fix index length to avoid concurrency modification

	// TODO: run in parallel
	for j := 0; j < n; j++ {

		<-i.links[j].unlockSignal // wait till some operations done on this index

		currentSortedLinks := i.links[j].sortedLinks
		newSortedLinks := make(sortableCidNumbers, 0, len(currentSortedLinks))
		for _, cidNumber := range currentSortedLinks {
			newRankedCid := RankedCidNumber{cidNumber.number, i.GetRankValue(cidNumber.number)}
			newSortedLinks = append(newSortedLinks, newRankedCid)
		}
		sort.Stable(sort.Reverse(newSortedLinks))
		i.links[j].sortedLinks = newSortedLinks
		i.links[j].Unlock()
	}

	// same process for backlinks
	n = len(i.backlinks)

	for j := 0; j < n; j++ {

		<-i.backlinks[j].unlockSignal // wait till some operations done on this index

		currentSortedLinks := i.backlinks[j].sortedLinks
		newSortedLinks := make(sortableCidNumbers, 0, len(currentSortedLinks))
		for _, cidNumber := range currentSortedLinks {
			newRankedCid := RankedCidNumber{cidNumber.number, i.GetRankValue(cidNumber.number)}
			newSortedLinks = append(newSortedLinks, newRankedCid)
		}
		sort.Stable(sort.Reverse(newSortedLinks))
		i.backlinks[j].sortedLinks = newSortedLinks
		i.backlinks[j].Unlock()
	}
}

func (i *BaseSearchIndex) lock() {
	i.locked = true
}

func (i *BaseSearchIndex) unlock() {
	i.locked = false
}

func (i *BaseSearchIndex) checkIndexError() error {
	select {
	case err := <-i.errChan:
		return err
	default:
	}
	return nil
}
