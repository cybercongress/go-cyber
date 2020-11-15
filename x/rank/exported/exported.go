package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cybercongress/go-cyber/merkle"
	"github.com/cybercongress/go-cyber/x/link"
	"github.com/cybercongress/go-cyber/x/rank/types"
)

type StateKeeper interface {
	Load(sdk.Context, log.Logger)

	BuildSearchIndex(log.Logger) types.SearchIndex
	GetIndexError() error

	Search(cidNumber link.CidNumber, page, perPage int) ([]types.RankedCidNumber, int, error)
	Top(page, perPage int) ([]types.RankedCidNumber, int, error)

	GetRankValue(link.CidNumber) uint64
	GetNetworkRankHash() []byte
	GetLastCidNum() link.CidNumber
	GetMerkleTree() *merkle.Tree
	GetLatestBlockNumber(ctx sdk.Context) uint64

}
