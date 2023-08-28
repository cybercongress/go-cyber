package types

import (
	"errors"

	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
)

type RankedCidNumber struct {
	number graphtypes.CidNumber
	rank   uint64
}

func (c RankedCidNumber) GetNumber() graphtypes.CidNumber { return c.number }
func (c RankedCidNumber) GetRank() uint64                 { return c.rank }

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
	Load(links graphtypes.Links)
	Search(cidNumber graphtypes.CidNumber, page, perPage uint32) ([]RankedCidNumber, uint32, error)
	Backlinks(cidNumber graphtypes.CidNumber, page, perPage uint32) ([]RankedCidNumber, uint32, error)
	Top(page, perPage uint32) ([]RankedCidNumber, uint32, error)
	GetRankValue(cidNumber graphtypes.CidNumber) uint64
	PutNewLinks(links []graphtypes.CompactLink)
	PutNewRank(rank Rank)
}

type NoopSearchIndex struct{}

func (i NoopSearchIndex) Run() GetError {
	return func() error {
		return nil
	}
}

func (i NoopSearchIndex) Load(links graphtypes.Links) {}
func (i NoopSearchIndex) Search(cidNumber graphtypes.CidNumber, page, perPage uint32) ([]RankedCidNumber, uint32, error) {
	return nil, 0, errors.New("The search API is not enabled on this node")
}

func (i NoopSearchIndex) Backlinks(cidNumber graphtypes.CidNumber, page, perPage uint32) ([]RankedCidNumber, uint32, error) {
	return nil, 0, errors.New("The search API is not enabled on this node")
}

func (i NoopSearchIndex) Top(page, perPage uint32) ([]RankedCidNumber, uint32, error) {
	return nil, 0, errors.New("The search API is not enabled on this node")
}
func (i NoopSearchIndex) PutNewLinks(links []graphtypes.CompactLink) {}
func (i NoopSearchIndex) PutNewRank(rank Rank)                       {}
func (i NoopSearchIndex) GetRankValue(cidNumber graphtypes.CidNumber) uint64 {
	return 0
}
