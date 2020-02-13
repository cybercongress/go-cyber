package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cybercongress/cyberd/merkle"
	"github.com/cybercongress/cyberd/x/link"
	"github.com/cybercongress/cyberd/x/rank/internal/types"
)

type StateKeeper interface {
	SetParams(sdk.Context, types.Params)
	GetParams(sdk.Context) types.Params

	Load(sdk.Context, log.Logger)
	BuildSearchIndex(log.Logger) types.SearchIndex

	EndBlocker(sdk.Context, log.Logger)

	Search(cidNumber link.CidNumber, page, perPage int) ([]types.RankedCidNumber, int, error)
	Top(page, perPage int) ([]types.RankedCidNumber, int, error)

	GetRankValue(link.CidNumber) float64
	GetNetworkRankHash() []byte

	GetLastCidNum() link.CidNumber
	GetMerkleTree() *merkle.Tree
	GetIndexError() error
}
