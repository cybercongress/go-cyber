package rank

import (
	"errors"
	. "github.com/cybercongress/cyberd/x/link/types"
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

type GetError func() error

type SearchIndex interface {
	Run() GetError
	Load(links Links)
	Search(cidNumber CidNumber, page, perPage int) ([]RankedCidNumber, int, error)
	GetRankValue(cidNumber CidNumber) float64
	PutNewLinks(links []CompactLink)
	PutNewRank(rank Rank)
}

type NoopSearchIndex struct{}

func (i NoopSearchIndex) Run() GetError {
	return func() error {
		return nil
	}
}

func (i NoopSearchIndex) Load(links Links) {}
func (i NoopSearchIndex) Search(cidNumber CidNumber, page, perPage int) ([]RankedCidNumber, int, error) {
	return nil, 0, errors.New("search is not allowed on this node")
}
func (i NoopSearchIndex) PutNewLinks(links []CompactLink) {}
func (i NoopSearchIndex) PutNewRank(rank Rank)            {}
func (i NoopSearchIndex) GetRankValue(cidNumber CidNumber) float64 {
	return 0
}
