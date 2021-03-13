package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cybercongress/go-cyber/merkle"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
	"github.com/cybercongress/go-cyber/x/rank/types"
)

type StateKeeper interface {
	Load(sdk.Context, log.Logger)

	BuildSearchIndex(log.Logger) types.SearchIndex
	GetIndexError() error

	Search(cidNumber graphtypes.CidNumber, page, perPage uint32) ([]types.RankedCidNumber, uint32, error)
	Backlinks(cidNumber graphtypes.CidNumber, page, perPage uint32) ([]types.RankedCidNumber, uint32, error)
	Top(page, perPage uint32) ([]types.RankedCidNumber, uint32, error)

	GetRankValue(graphtypes.CidNumber) uint64
	GetNetworkRankHash() []byte
	GetLastCidNum() graphtypes.CidNumber
	GetMerkleTree() *merkle.Tree
	GetLatestBlockNumber(ctx sdk.Context) uint64
}
