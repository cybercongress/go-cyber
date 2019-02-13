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
	sortedLinks sortedCidNumbers

	unlockSignal chan struct{}
}

func NewCidLinks() cidLinks {
	return cidLinks{
		sortedLinks:  make(sortedCidNumbers, 0),
		unlockSignal: make(chan struct{}, 1),
	}
}

type sortedCidNumbers []RankedCidNumber

// Sort Interface functions
func (links sortedCidNumbers) Len() int { return len(links) }

func (links sortedCidNumbers) Less(i, j int) bool { return links[i].rank < links[j].rank }

func (links sortedCidNumbers) Swap(i, j int) { links[i], links[j] = links[j], links[i] }

// Send unlock signal so others could operate on this index
func (links cidLinks) Unlock() {
	links.unlockSignal <- struct{}{}
}

//todo: put rank values here
type SearchIndex struct {
	linksIndex []cidLinks
	rank       Rank

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
	i.linksIndex = make([]cidLinks, 0, 1000000)
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
		return nil, 0, errors.New("linksIndex currently unavailable after node restart")
	}

	if uint64(cidNumber) >= uint64(len(i.linksIndex)) {
		return []RankedCidNumber{}, 0, nil
	}

	cidRankedLinks := i.linksIndex[cidNumber]
	if len(cidRankedLinks.sortedLinks) == 0 {
		return []RankedCidNumber{}, 0, nil
	}

	totalSize := len(cidRankedLinks.sortedLinks)
	startIndex := page * perPage
	if startIndex >= totalSize {
		return nil, totalSize, errors.New("page not found")
	}

	endIndex := startIndex + perPage
	if endIndex > totalSize {
		endIndex = startIndex + (totalSize % perPage)
	}

	resultSet := cidRankedLinks.sortedLinks[startIndex:endIndex]

	return resultSet, totalSize, nil
}

// make sure that this link (from-to) is new
func (i *SearchIndex) handleLink(link CompactLink) {

	i.extendIndex(uint64(link.From()))

	fromIndex := i.linksIndex[link.From()]
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
	indexLen := uint64(len(i.linksIndex))
	if fromCidNumber >= indexLen {
		for j := indexLen; j <= fromCidNumber; j++ {
			links := NewCidLinks()
			links.Unlock() // allow operations on this index
			i.linksIndex = append(i.linksIndex, links)
		}
	}
}

func (i *SearchIndex) putLinkIntoIndex(from CidNumber, to CidNumber) {
	index := i.linksIndex[uint64(from)].sortedLinks
	// todo: not optimal. replace with some another implementation. may be AVL tree
	rankedTo := RankedCidNumber{to, i.getRankValue(to)}
	pos := sort.Search(len(index), func(i int) bool { return index[i].rank < rankedTo.rank })
	index = append(index, RankedCidNumber{})
	copy(index[pos+1:], index[pos:])
	index[pos] = rankedTo
	i.linksIndex[uint64(from)].sortedLinks = index
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
		case link := <-i.LinksChan: //todo: channel should be buffered
			i.logger.Info("New link to index")
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
			i.logger.Info("New rank to index")
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
	n := len(i.linksIndex) // fix index length to avoid concurrency modification

	// todo: run in parallel
	for j := 0; j < n; j++ {

		<-i.linksIndex[j].unlockSignal // wait till some operations done on this index

		currentSortedLinks := i.linksIndex[j].sortedLinks
		newSortedLinks := make(sortedCidNumbers, 0, len(currentSortedLinks))
		for _, cidNumber := range currentSortedLinks {
			newRankedCid := RankedCidNumber{cidNumber.number, i.getRankValue(cidNumber.number)}
			newSortedLinks = append(newSortedLinks, newRankedCid)
		}
		sort.Stable(sort.Reverse(newSortedLinks))
		i.linksIndex[j].sortedLinks = newSortedLinks
		i.linksIndex[j].Unlock()
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
