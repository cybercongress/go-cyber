package rank

import (
	"errors"
	"github.com/cybercongress/cyberd/x/link/keeper"
	. "github.com/cybercongress/cyberd/x/link/types"
	"github.com/tendermint/tendermint/libs/log"
	"sort"
	"time"
)

type RankedCidNumber struct {
	number CidNumber
	rank   float64
}

func (c RankedCidNumber) GetNumber() CidNumber { return c.number }
func (c RankedCidNumber) GetRank() float64     { return c.rank }

//
// Local type for sorting
type cidLinks struct {
	sortedLinks sortableCidNumbers

	unlockSignal chan struct{}
}

func NewCidLinks() cidLinks {
	return cidLinks{
		sortedLinks:  make(sortableCidNumbers, 0),
		unlockSignal: make(chan struct{}, 1),
	}
}

type sortableCidNumbers []RankedCidNumber

// Sort Interface functions
func (links sortableCidNumbers) Len() int { return len(links) }

func (links sortableCidNumbers) Less(i, j int) bool { return links[i].rank < links[j].rank }

func (links sortableCidNumbers) Swap(i, j int) { links[i], links[j] = links[j], links[i] }

// Send unlock signal so others could operate on this index
func (links cidLinks) Unlock() {
	links.unlockSignal <- struct{}{}
}

//todo: put rank values here
type SearchIndex struct {
	links []cidLinks
	rank  Rank

	LinksChan chan CompactLink
	RankChan  chan Rank
	errChan   chan error

	locked bool

	logger log.Logger
}

func NewSearchIndex(log log.Logger) *SearchIndex {
	return &SearchIndex{
		LinksChan: make(chan CompactLink, 1000),
		RankChan:  make(chan Rank, 1),
		errChan:   make(chan error),
		locked:    true,
		logger:    log,
	}
}

func (i *SearchIndex) Lock() {
	i.locked = true
}

func (i *SearchIndex) Unlock() {
	i.locked = false
}

// Load links with zero rank values. No sorting. Index should be unavailable for read
func (i *SearchIndex) Load(linkIndexedKeeper *keeper.LinkIndexedKeeper) {
	i.Lock() // lock index for read
	i.links = make([]cidLinks, 0, 1000000)
	startTime := time.Now()
	for from, toCids := range linkIndexedKeeper.GetNextOutLinks() {
		i.extendIndex(uint64(from))

		for to := range toCids {
			i.putLinkIntoIndex(from, to)
		}
	}
	i.logger.Info("Search index loaded!", "time", time.Since(startTime))
}

func (i *SearchIndex) Search(cidNumber CidNumber, page, perPage int) ([]RankedCidNumber, int, error) {

	i.logger.Info("Search query", "cid", cidNumber, "page", page, "perPage", perPage)

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

	totalSize := len(cidLinks.sortedLinks)
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

// make sure that this link (from-to) is new
func (i *SearchIndex) handleLink(link CompactLink) {

	i.extendIndex(uint64(link.From()))

	fromIndex := i.links[link.From()]
	// in case unlock signal received we could operate on this index otherwise put link in the end of queue and finish
	select {
	case _ = <-fromIndex.unlockSignal:
		i.putLinkIntoIndex(link.From(), link.To())
		fromIndex.Unlock()
		break
	default:
		i.LinksChan <- link
	}
}

func (i *SearchIndex) getRankValue(cid CidNumber) float64 {
	if i.rank.Values == nil || uint64(len(i.rank.Values)) <= uint64(cid) {
		return 0
	}
	return i.rank.Values[cid]
}

func (i *SearchIndex) extendIndex(fromCidNumber uint64) {
	indexLen := uint64(len(i.links))
	if fromCidNumber >= indexLen {
		for j := indexLen; j <= fromCidNumber; j++ {
			links := NewCidLinks()
			links.Unlock() // allow operations on this index
			i.links = append(i.links, links)
		}
	}
}

func (i *SearchIndex) putLinkIntoIndex(from CidNumber, to CidNumber) {
	fromLinks := i.links[uint64(from)].sortedLinks
	// todo: not optimal. replace with some another implementation. may be AVL tree
	rankedTo := RankedCidNumber{to, i.getRankValue(to)}
	pos := sort.Search(len(fromLinks), func(i int) bool { return fromLinks[i].rank < rankedTo.rank })
	fromLinks = append(fromLinks, RankedCidNumber{})
	copy(fromLinks[pos+1:], fromLinks[pos:])
	fromLinks[pos] = rankedTo
	i.links[uint64(from)].sortedLinks = fromLinks
}

// for parallel usage
func (i *SearchIndex) startListenNewLinks() {
	defer func() {
		if r := recover(); r != nil {
			i.errChan <- r.(error)
		}
	}()

	i.logger.Info("Search index starting listen new links")
	for {
		select {
		case link := <-i.LinksChan:
			i.handleLink(link)
			break
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// for parallel usage
func (i *SearchIndex) startListenNewRank() {
	defer func() {
		if r := recover(); r != nil {
			i.errChan <- r.(error)
		}
	}()

	i.logger.Info("Search index starting listen new rank")
	for {
		select {
		case rank := <-i.RankChan: //todo: could be problems if recalculation lasts more than rank period
			i.rank = rank
			i.recalculateIndices()
			break
		default:
			time.Sleep(100 * time.Millisecond)
		}

	}
}

func (i *SearchIndex) recalculateIndices() {
	defer i.Unlock()
	n := len(i.links) // fix index length to avoid concurrency modification

	// todo: run in parallel
	for j := 0; j < n; j++ {

		<-i.links[j].unlockSignal // wait till some operations done on this index

		currentSortedLinks := i.links[j].sortedLinks
		newSortedLinks := make(sortableCidNumbers, 0, len(currentSortedLinks))
		for _, cidNumber := range currentSortedLinks {
			newRankedCid := RankedCidNumber{cidNumber.number, i.getRankValue(cidNumber.number)}
			newSortedLinks = append(newSortedLinks, newRankedCid)
		}
		sort.Stable(sort.Reverse(newSortedLinks))
		i.links[j].sortedLinks = newSortedLinks
		i.links[j].Unlock()
	}
}

func (i *SearchIndex) checkIndexError() error {
	select {
	case err := <-i.errChan:
		return err
	default:
	}
	return nil
}
